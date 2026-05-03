package app

import (
	"context"
	"fmt"
	"time"

	"reqium/internal/interfaces"
	"reqium/internal/models"
)

type RequestService struct {
	client      interfaces.HTTPClient
	formatter   interfaces.Formatter
	historyRepo interfaces.HistoryRepository
}

func NewRequestService(client interfaces.HTTPClient, formatter interfaces.Formatter) *RequestService {
	return &RequestService{client: client, formatter: formatter}
}

func NewRequestServiceWithHistory(client interfaces.HTTPClient, formatter interfaces.Formatter, historyRepo interfaces.HistoryRepository) *RequestService {
	return &RequestService{client: client, formatter: formatter, historyRepo: historyRepo}
}

func (s *RequestService) Send(ctx context.Context, req models.Request) (models.Response, error) {
	if err := req.Validate(); err != nil {
		return models.Response{}, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()

	response, err := s.client.Do(timeoutCtx, req)
	if err != nil {
		_ = s.saveHistory(ctx, req, models.Response{}, err)
		return models.Response{}, fmt.Errorf("request failed: %w", err)
	}

	if err := s.saveHistory(ctx, req, response, nil); err != nil {
		return models.Response{}, fmt.Errorf("save history: %w", err)
	}

	return response, nil
}

func (s *RequestService) SendFormatted(ctx context.Context, req models.Request) (string, error) {
	if s.formatter == nil {
		return "", fmt.Errorf("formatter is required")
	}

	response, err := s.Send(ctx, req)
	if err != nil {
		return "", err
	}
	output, err := s.formatter.Format(response)
	if err != nil {
		return "", fmt.Errorf("format response: %w", err)
	}

	return output, nil
}

func (s *RequestService) saveHistory(ctx context.Context, req models.Request, response models.Response, requestErr error) error {
	if s.historyRepo == nil {
		return nil
	}

	entry := models.HistoryEntry{
		ID:         newID(),
		Method:     req.Method,
		URL:        req.URL,
		Headers:    req.Headers,
		Body:       req.Body,
		StatusCode: response.StatusCode,
		Response:   response.Body,
		Duration:   response.Duration,
		ExecutedAt: time.Now(),
	}
	if requestErr != nil {
		entry.Error = requestErr.Error()
	}
	return s.historyRepo.Save(ctx, entry)
}

func newID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func requestFromHistory(entry models.HistoryEntry, timeout time.Duration) models.Request {
	return models.Request{
		Method:  entry.Method,
		URL:     entry.URL,
		Headers: entry.Headers,
		Body:    entry.Body,
		Timeout: timeout,
	}
}
