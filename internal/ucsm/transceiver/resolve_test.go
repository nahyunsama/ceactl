package transceiver

import (
	"context"
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResolveClass_SendsExpectedPayload(t *testing.T) {
	var got resolveClassRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if err := xml.Unmarshal(body, &got); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<configResolveClass><outConfigs></outConfigs></configResolveClass>`))
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTP: server.Client(), Cookie: "abc123"}

	data, err := client.ResolveClass(context.Background(), "computeItem")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Cookie != "abc123" {
		t.Errorf("server received cookie %q, want abc123", got.Cookie)
	}
	if got.ClassID != "computeItem" {
		t.Errorf("server received classId %q, want computeItem", got.ClassID)
	}
	if got.InHierarchical != "false" {
		t.Errorf("server received inHierarchical %q, want false", got.InHierarchical)
	}
	if string(data) != `<configResolveClass><outConfigs></outConfigs></configResolveClass>` {
		t.Errorf("got response %q, unexpected", data)
	}
}

func TestResolveClass_PropagatesPostXMLError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<error>boom</error>`))
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTP: server.Client()}

	_, err := client.ResolveClass(context.Background(), "computeItem")
	if err == nil {
		t.Fatal("expected an error when the server returns a non-200 status, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error %q does not mention status code 500", err.Error())
	}
}
