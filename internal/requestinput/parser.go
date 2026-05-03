package requestinput

import (
	"fmt"
	"strings"
	"time"

	"reqium/internal/enums"
	"reqium/internal/interfaces"
	"reqium/internal/models"
)

type BodyOptions struct {
	Body     string
	BodyFile string
}

func ParseHeaders(values []string) (map[string]string, error) {
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

func LoadBody(opts BodyOptions, reader interfaces.FileReader) ([]byte, error) {
	if opts.Body != "" && opts.BodyFile != "" {
		return nil, fmt.Errorf("--body and --body-file cannot be used together")
	}
	if opts.BodyFile != "" {
		return reader.Read(opts.BodyFile)
	}
	if opts.Body != "" {
		return []byte(opts.Body), nil
	}
	return nil, nil
}

func BuildRequest(method string, url string, headers []string, bodyOpts BodyOptions, timeoutSec int, reader interfaces.FileReader) (models.Request, error) {
	parsedHeaders, err := ParseHeaders(headers)
	if err != nil {
		return models.Request{}, err
	}

	body, err := LoadBody(bodyOpts, reader)
	if err != nil {
		return models.Request{}, err
	}

	if len(body) > 0 && !enums.MethodAllowsBody(method) {
		return models.Request{}, fmt.Errorf("body is only accepted for POST, PUT, and PATCH")
	}

	return models.Request{
		Method:  method,
		URL:     url,
		Headers: parsedHeaders,
		Body:    body,
		Timeout: time.Duration(timeoutSec) * time.Second,
	}, nil
}
