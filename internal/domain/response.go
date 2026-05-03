package domain

import "time"

type Response struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
	Duration   time.Duration
}
