package app

import (
	"context"
	"time"

	"reqium/internal/interfaces"
	"reqium/internal/models"
)

type HistoryService struct {
	repo           interfaces.HistoryRepository
	requestService *RequestService
}

func NewHistoryService(repo interfaces.HistoryRepository, requestService *RequestService) *HistoryService {
	return &HistoryService{repo: repo, requestService: requestService}
}

func (s *HistoryService) List(ctx context.Context, limit int) ([]models.HistoryEntry, error) {
	return s.repo.List(ctx, limit)
}

func (s *HistoryService) Get(ctx context.Context, id string) (models.HistoryEntry, error) {
	return s.repo.Get(ctx, id)
}

func (s *HistoryService) Replay(ctx context.Context, id string, timeout time.Duration) (models.Response, error) {
	entry, err := s.repo.Get(ctx, id)
	if err != nil {
		return models.Response{}, err
	}
	return s.requestService.Send(ctx, requestFromHistory(entry, timeout))
}
