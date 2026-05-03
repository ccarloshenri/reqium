package workflowtest

import (
	"errors"
	"testing"
	"time"

	reqerrors "reqium/internal/errors"
	"reqium/internal/models"
)

func TestRequestValidate(t *testing.T) {
	tests := []struct {
		name string
		req  models.Request
		err  error
	}{
		{
			name: "valid request",
			req: models.Request{
				Method:  "GET",
				URL:     "https://api.example.com/users",
				Headers: map[string]string{},
				Timeout: time.Second,
			},
		},
		{
			name: "missing url",
			req: models.Request{
				Method:  "GET",
				Headers: map[string]string{},
				Timeout: time.Second,
			},
			err: reqerrors.ErrMissingURL,
		},
		{
			name: "invalid url",
			req: models.Request{
				Method:  "GET",
				URL:     "not-a-url",
				Headers: map[string]string{},
				Timeout: time.Second,
			},
			err: reqerrors.ErrInvalidURL,
		},
		{
			name: "invalid method",
			req: models.Request{
				Method:  "TRACE",
				URL:     "https://api.example.com/users",
				Headers: map[string]string{},
				Timeout: time.Second,
			},
			err: reqerrors.ErrInvalidMethod,
		},
		{
			name: "invalid timeout",
			req: models.Request{
				Method:  "GET",
				URL:     "https://api.example.com/users",
				Headers: map[string]string{},
			},
			err: reqerrors.ErrInvalidTimeout,
		},
		{
			name: "invalid json body",
			req: models.Request{
				Method:  "POST",
				URL:     "https://api.example.com/users",
				Headers: map[string]string{},
				Body:    []byte(`{"name":`),
				Timeout: time.Second,
			},
			err: reqerrors.ErrInvalidJSON,
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
