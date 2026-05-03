package models

import "time"

type HistoryEntry struct {
	ID          string            `json:"id"`
	Name        string            `json:"name,omitempty"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Body        []byte            `json:"body,omitempty"`
	StatusCode  int               `json:"status_code"`
	Response    []byte            `json:"response,omitempty"`
	Duration    time.Duration     `json:"duration"`
	ExecutedAt  time.Time         `json:"executed_at"`
	Error       string            `json:"error,omitempty"`
	Environment string            `json:"environment,omitempty"`
}
