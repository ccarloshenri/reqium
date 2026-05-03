package domain

import "errors"

var (
	ErrInvalidURL     = errors.New("invalid url")
	ErrMissingURL     = errors.New("url is required")
	ErrInvalidMethod  = errors.New("invalid http method")
	ErrInvalidTimeout = errors.New("timeout must be greater than zero")
	ErrInvalidJSON    = errors.New("invalid json body")
)
