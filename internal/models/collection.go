package models

import "time"

type Collection struct {
	Name      string         `json:"name"`
	Requests  []SavedRequest `json:"requests"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type SavedRequest struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body,omitempty"`
}

func (r SavedRequest) ToRequest(timeout time.Duration) Request {
	return Request{
		Method:  r.Method,
		URL:     r.URL,
		Headers: r.Headers,
		Body:    r.Body,
		Timeout: timeout,
	}
}
