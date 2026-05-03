package cli

import (
	"fmt"
	"strconv"

	"reqium/internal/enums"
	reqerrors "reqium/internal/errors"
	"reqium/internal/interfaces"
	"reqium/internal/models"
	"reqium/internal/requestinput"
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
	return requestinput.ParseHeaders(values)
}

func loadBody(opts requestOptions, reader interfaces.FileReader) ([]byte, error) {
	return requestinput.LoadBody(requestinput.BodyOptions{Body: opts.body, BodyFile: opts.bodyFile}, reader)
}

func buildRequest(method string, url string, opts requestOptions, reader interfaces.FileReader) (models.Request, error) {
	return requestinput.BuildRequest(
		method,
		url,
		opts.headers,
		requestinput.BodyOptions{Body: opts.body, BodyFile: opts.bodyFile},
		opts.timeoutSec,
		reader,
	)
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
