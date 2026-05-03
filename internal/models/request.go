package models

import (
	"encoding/json"
	"net/url"
	"strings"
	"time"

	"reqium/internal/enums"
	reqerrors "reqium/internal/errors"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

func (r Request) Validate() error {
	if strings.TrimSpace(r.URL) == "" {
		return reqerrors.ErrMissingURL
	}

	parsed, err := url.ParseRequestURI(r.URL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return reqerrors.ErrInvalidURL
	}

	if !enums.ValidHTTPMethod(strings.ToUpper(r.Method)) {
		return reqerrors.ErrInvalidMethod
	}

	if r.Timeout <= 0 {
		return reqerrors.ErrInvalidTimeout
	}

	if len(r.Body) > 0 && !json.Valid(r.Body) {
		return reqerrors.ErrInvalidJSON
	}

	return nil
}
