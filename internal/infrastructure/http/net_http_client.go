package http

import (
	"bytes"
	"context"
	"io"
	nethttp "net/http"
	"time"

	"reqium/internal/domain"
)

type NetHTTPClient struct {
	client *nethttp.Client
}

func NewNetHTTPClient() *NetHTTPClient {
	return &NetHTTPClient{client: &nethttp.Client{}}
}

func (c *NetHTTPClient) Do(ctx context.Context, req domain.Request) (domain.Response, error) {
	httpReq, err := nethttp.NewRequestWithContext(ctx, req.Method, req.URL, bytes.NewReader(req.Body))
	if err != nil {
		return domain.Response{}, err
	}

	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	start := time.Now()
	httpResp, err := c.client.Do(httpReq)
	duration := time.Since(start)
	if err != nil {
		return domain.Response{}, err
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return domain.Response{}, err
	}

	return domain.Response{
		StatusCode: httpResp.StatusCode,
		Headers:    map[string][]string(httpResp.Header),
		Body:       body,
		Duration:   duration,
	}, nil
}
