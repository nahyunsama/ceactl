package llmanalysis

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChat_SendsExpectedPayload(t *testing.T) {
	var gotPath, gotContentType string
	var got chatRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotContentType = r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":{"role":"assistant","content":"hello back"}}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "gemma4:e2b")
	client.HTTP = server.Client()

	gotReply, err := client.Chat(context.Background(), "sys prompt", "user prompt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotReply != "hello back" {
		t.Errorf("got reply %q, want %q", gotReply, "hello back")
	}
	if gotPath != "/api/chat" {
		t.Errorf("got path %q, want /api/chat", gotPath)
	}
	if gotContentType != "application/json" {
		t.Errorf("got Content-Type %q, want application/json", gotContentType)
	}
	if got.Model != "gemma4:e2b" {
		t.Errorf("got model %q, want gemma4:e2b", got.Model)
	}
	if got.Stream {
		t.Errorf("got stream=true, want false")
	}
	if got.Options.NumCtx != defaultNumCtx {
		t.Errorf("got num_ctx %d, want %d", got.Options.NumCtx, defaultNumCtx)
	}
	if len(got.Messages) != 2 || got.Messages[0].Role != "system" || got.Messages[0].Content != "sys prompt" {
		t.Errorf("got messages[0] = %+v, want system/sys prompt", got.Messages[0])
	}
	if len(got.Messages) != 2 || got.Messages[1].Role != "user" || got.Messages[1].Content != "user prompt" {
		t.Errorf("got messages[1] = %+v, want user/user prompt", got.Messages[1])
	}
}

func TestChat_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"boom"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "gemma4:e2b")
	client.HTTP = server.Client()

	_, err := client.Chat(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("expected an error for non-200 status code, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error %q does not mention status code 500", err.Error())
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("error %q does not include response body", err.Error())
	}
}

func TestChat_UnreachableServer(t *testing.T) {
	client := NewClient("http://127.0.0.1:0", "gemma4:e2b")

	_, err := client.Chat(context.Background(), "sys", "user")
	if err == nil {
		t.Fatal("expected an error when the server is unreachable, got nil")
	}
}
