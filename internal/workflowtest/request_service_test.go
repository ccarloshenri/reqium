package workflowtest

import (
	"context"
	"errors"
	"testing"
	"time"

	"reqium/internal/app"
	reqerrors "reqium/internal/errors"
	"reqium/internal/models"
)

type mockHTTPClient struct {
	request models.Request
	called  bool
	err     error
}

func (c *mockHTTPClient) Do(ctx context.Context, req models.Request) (models.Response, error) {
	c.called = true
	c.request = req
	if c.err != nil {
		return models.Response{}, c.err
	}
	return models.Response{StatusCode: 200, Body: []byte(`{"ok":true}`)}, nil
}

type mockFormatter struct {
	called bool
	err    error
}

func (f *mockFormatter) Format(response models.Response) (string, error) {
	f.called = true
	if f.err != nil {
		return "", f.err
	}
	return "formatted", nil
}

func TestRequestServiceSendFormattedOrchestratesClientAndFormatter(t *testing.T) {
	client := &mockHTTPClient{}
	formatter := &mockFormatter{}
	service := app.NewRequestService(client, formatter)

	output, err := service.SendFormatted(context.Background(), models.Request{
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
	service := app.NewRequestService(client, &mockFormatter{})

	_, err := service.Send(context.Background(), models.Request{
		Method:  "GET",
		Headers: map[string]string{},
		Timeout: time.Second,
	})
	if !errors.Is(err, reqerrors.ErrMissingURL) {
		t.Fatalf("expected missing url, got %v", err)
	}
	if client.called {
		t.Fatal("client should not be called for invalid request")
	}
}
