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

func TestLogin_Success(t *testing.T) {
	var got loginRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if err := xml.Unmarshal(body, &got); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<aaaLogin outCookie="abc123"/>`))
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTP: server.Client()}

	if err := client.Login(context.Background(), "admin", "secret"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client.Cookie != "abc123" {
		t.Errorf("got Cookie %q, want abc123", client.Cookie)
	}
	if got.InName != "admin" || got.InPassword != "secret" {
		t.Errorf("server received (inName=%q, inPassword=%q), want (admin, secret)", got.InName, got.InPassword)
	}
}

func TestLogin_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<aaaLogin errorCode="551" errorDescr="Authentication failed"/>`))
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTP: server.Client()}

	err := client.Login(context.Background(), "admin", "wrong-password")
	if err == nil {
		t.Fatal("expected an error for failed login, got nil")
	}
	if !strings.Contains(err.Error(), "Authentication failed") {
		t.Errorf("error %q does not contain the server's errorDescr", err.Error())
	}
	if client.Cookie != "" {
		t.Errorf("got Cookie %q, want empty on failed login", client.Cookie)
	}
}

func TestLogout_NoCookieSkipsRequest(t *testing.T) {
	called := false

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTP: server.Client()}

	if err := client.Logout(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("Logout should not contact the server when Cookie is empty")
	}
}

func TestLogout_Success(t *testing.T) {
	var got logoutRequest

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if err := xml.Unmarshal(body, &got); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<aaaLogout/>`))
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTP: server.Client(), Cookie: "abc123"}

	if err := client.Logout(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client.Cookie != "" {
		t.Errorf("got Cookie %q, want empty after successful logout", client.Cookie)
	}
	if got.InCookie != "abc123" {
		t.Errorf("server received inCookie %q, want abc123", got.InCookie)
	}
}

func TestLogout_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<aaaLogout errorCode="552" errorDescr="Invalid cookie"/>`))
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTP: server.Client(), Cookie: "abc123"}

	err := client.Logout(context.Background())
	if err == nil {
		t.Fatal("expected an error for failed logout, got nil")
	}
	if !strings.Contains(err.Error(), "Invalid cookie") {
		t.Errorf("error %q does not contain the server's errorDescr", err.Error())
	}
	if client.Cookie != "abc123" {
		t.Errorf("got Cookie %q, want unchanged (abc123) after failed logout", client.Cookie)
	}
}
