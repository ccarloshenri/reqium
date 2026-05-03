package errors

import stderrors "errors"

var (
	ErrInvalidURL       = stderrors.New("invalid url")
	ErrMissingURL       = stderrors.New("url is required")
	ErrInvalidMethod    = stderrors.New("invalid http method")
	ErrInvalidTimeout   = stderrors.New("timeout must be greater than zero")
	ErrInvalidJSON      = stderrors.New("invalid json body")
	ErrNotFound         = stderrors.New("not found")
	ErrVariableNotFound = stderrors.New("variable not found")
)
