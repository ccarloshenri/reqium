package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"reqium/internal/domain"
)

type mockHTTPClient struct {
	request domain.Request
	called  bool
	err     error
}

func (c *mockHTTPClient) Do(ctx context.Context, req domain.Request) (domain.Response, error) {
	c.called = true
	c.request = req
	if c.err != nil {
		return domain.Response{}, c.err
	}
	return domain.Response{StatusCode: 200, Body: []byte(`{"ok":true}`)}, nil
}

type mockFormatter struct {
	called bool
	err    error
}

func (f *mockFormatter) Format(response domain.Response) (string, error) {
	f.called = true
	if f.err != nil {
		return "", f.err
	}
	return "formatted", nil
}

func TestRequestServiceSendOrchestratesClientAndFormatter(t *testing.T) {
	client := &mockHTTPClient{}
	formatter := &mockFormatter{}
	service := NewRequestService(client, formatter)

	output, err := service.Send(context.Background(), domain.Request{
		Method:  "GET",
		URL:     "https://api.example.com/users",
		Headers: map[string]string{},
		Timeout: time.Second,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "formatted" {
		t.Fatalf("expected formatted output, got %q", output)
	}
	if !client.called {
		t.Fatal("expected client to be called")
	}
	if !formatter.called {
		t.Fatal("expected formatter to be called")
	}
}

func TestRequestServiceSendReturnsValidationErrorBeforeClient(t *testing.T) {
	client := &mockHTTPClient{}
	service := NewRequestService(client, &mockFormatter{})

	_, err := service.Send(context.Background(), domain.Request{
		Method:  "GET",
		Headers: map[string]string{},
		Timeout: time.Second,
	})
	if !errors.Is(err, domain.ErrMissingURL) {
		t.Fatalf("expected missing url, got %v", err)
	}
	if client.called {
		t.Fatal("client should not be called for invalid request")
	}
}
