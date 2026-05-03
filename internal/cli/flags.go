package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"reqium/internal/enums"
	reqerrors "reqium/internal/errors"
	"reqium/internal/interfaces"
	"reqium/internal/models"
)

type requestOptions struct {
	headers    []string
	body       string
	bodyFile   string
	timeoutSec int
	pretty     bool
	env        string
}

func parseHeaders(values []string) (map[string]string, error) {
	headers := make(map[string]string, len(values))
	for _, header := range values {
		key, value, ok := strings.Cut(header, ":")
		if !ok || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("invalid header %q: expected Key: Value", header)
		}
		headers[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return headers, nil
}

func loadBody(opts requestOptions, reader interfaces.FileReader) ([]byte, error) {
	if opts.body != "" && opts.bodyFile != "" {
		return nil, fmt.Errorf("--body and --body-file cannot be used together")
	}
	if opts.bodyFile != "" {
		return reader.Read(opts.bodyFile)
	}
	if opts.body != "" {
		return []byte(opts.body), nil
	}
	return nil, nil
}

func buildRequest(method string, url string, opts requestOptions, reader interfaces.FileReader) (models.Request, error) {
	headers, err := parseHeaders(opts.headers)
	if err != nil {
		return models.Request{}, err
	}

	body, err := loadBody(opts, reader)
	if err != nil {
		return models.Request{}, err
	}

	if len(body) > 0 && !methodAllowsBody(method) {
		return models.Request{}, fmt.Errorf("body is only accepted for POST, PUT, and PATCH")
	}

	return models.Request{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
		Timeout: time.Duration(opts.timeoutSec) * time.Second,
	}, nil
}

func parseTimeout(value string) (int, error) {
	seconds, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("timeout must be a whole number of seconds")
	}
	if seconds <= 0 {
		return 0, reqerrors.ErrInvalidTimeout
	}
	return seconds, nil
}

func methodAllowsBody(method string) bool {
	return enums.MethodAllowsBody(method)
}
