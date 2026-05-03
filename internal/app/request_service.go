package app

import (
	"context"
	"fmt"

	"reqium/internal/domain"
	"reqium/internal/interfaces"
)

type RequestService struct {
	client    interfaces.HTTPClient
	formatter interfaces.Formatter
}

func NewRequestService(client interfaces.HTTPClient, formatter interfaces.Formatter) *RequestService {
	return &RequestService{client: client, formatter: formatter}
}

func (s *RequestService) Send(ctx context.Context, req domain.Request) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()

	response, err := s.client.Do(timeoutCtx, req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}

	output, err := s.formatter.Format(response)
	if err != nil {
		return "", fmt.Errorf("format response: %w", err)
	}

	return output, nil
}
