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

type testPayload struct {
	XMLName xml.Name `xml:"testRequest"`
	Foo     string   `xml:"foo,attr"`
}

func TestPostXML_Success(t *testing.T) {
	var gotMethod, gotContentType string
	var gotBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotContentType = r.Header.Get("content-type")
		gotBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<aaaLogin outCookie="abc123"/>`))
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTP: server.Client()}

	got, err := client.PostXML(context.Background(), testPayload{Foo: "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(got) != `<aaaLogin outCookie="abc123"/>` {
		t.Errorf("got body %q, unexpected", got)
	}
	if gotMethod != http.MethodPost {
		t.Errorf("got method %q, want POST", gotMethod)
	}
	if gotContentType != "application/x-www-form-urlencoded" {
		t.Errorf("got content-type %q, want application/x-www-form-urlencoded", gotContentType)
	}
	if !strings.Contains(string(gotBody), `foo="bar"`) {
		t.Errorf("server received body %q, want it to contain foo=\"bar\"", gotBody)
	}
}

func TestPostXML_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`<error>boom</error>`))
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTP: server.Client()}

	_, err := client.PostXML(context.Background(), testPayload{Foo: "bar"})
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
