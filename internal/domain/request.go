package domain

import (
	"encoding/json"
	"net/url"
	"slices"
	"strings"
	"time"
)

var allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

func (r Request) Validate() error {
	if strings.TrimSpace(r.URL) == "" {
		return ErrMissingURL
	}

	parsed, err := url.ParseRequestURI(r.URL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return ErrInvalidURL
	}

	if !slices.Contains(allowedMethods, strings.ToUpper(r.Method)) {
		return ErrInvalidMethod
	}

	if r.Timeout <= 0 {
		return ErrInvalidTimeout
	}

	if len(r.Body) > 0 && !json.Valid(r.Body) {
		return ErrInvalidJSON
	}

	return nil
}
