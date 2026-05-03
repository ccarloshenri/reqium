package domain

import (
	"errors"
	"testing"
	"time"
)

func TestRequestValidate(t *testing.T) {
	tests := []struct {
		name string
		req  Request
		err  error
	}{
		{
			name: "valid request",
			req: Request{
				Method:  "GET",
				URL:     "https://api.example.com/users",
				Headers: map[string]string{},
				Timeout: time.Second,
			},
		},
		{
			name: "missing url",
			req: Request{
				Method:  "GET",
				Headers: map[string]string{},
				Timeout: time.Second,
			},
			err: ErrMissingURL,
		},
		{
			name: "invalid url",
			req: Request{
				Method:  "GET",
				URL:     "not-a-url",
				Headers: map[string]string{},
				Timeout: time.Second,
			},
			err: ErrInvalidURL,
		},
		{
			name: "invalid method",
			req: Request{
				Method:  "TRACE",
				URL:     "https://api.example.com/users",
				Headers: map[string]string{},
				Timeout: time.Second,
			},
			err: ErrInvalidMethod,
		},
		{
			name: "invalid timeout",
			req: Request{
				Method:  "GET",
				URL:     "https://api.example.com/users",
				Headers: map[string]string{},
			},
			err: ErrInvalidTimeout,
		},
		{
			name: "invalid json body",
			req: Request{
				Method:  "POST",
				URL:     "https://api.example.com/users",
				Headers: map[string]string{},
				Body:    []byte(`{"name":`),
				Timeout: time.Second,
			},
			err: ErrInvalidJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.err) {
				t.Fatalf("expected %v, got %v", tt.err, err)
			}
		})
	}
}
