package httpclienttest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"reqium/internal/domain"
	httpinfra "reqium/internal/infrastructure/http"
)

func TestNetHTTPClientDoSendsRequestAndReturnsResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("X-Test") != "true" {
			t.Fatalf("expected X-Test header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"created":true}`))
	}))
	defer server.Close()

	response, err := httpinfra.NewNetHTTPClient().Do(context.Background(), domain.Request{
		Method:  "POST",
		URL:     server.URL,
		Headers: map[string]string{"X-Test": "true"},
		Body:    []byte(`{"name":"John"}`),
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.StatusCode)
	}
	if string(response.Body) != `{"created":true}` {
		t.Fatalf("unexpected body: %s", response.Body)
	}
	if response.Duration <= 0 {
		t.Fatal("expected duration to be recorded")
	}
	if response.Headers["Content-Type"][0] != "application/json" {
		t.Fatalf("unexpected content type: %v", response.Headers["Content-Type"])
	}
}

func TestNetHTTPClientDoReturnsContextTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	_, err := httpinfra.NewNetHTTPClient().Do(ctx, domain.Request{
		Method:  "GET",
		URL:     server.URL,
		Headers: map[string]string{},
		Timeout: time.Second,
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}
}
