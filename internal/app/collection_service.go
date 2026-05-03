package app

import (
	"context"
	"time"

	"reqium/internal/interfaces"
	"reqium/internal/models"
)

type CollectionService struct {
	repo interfaces.CollectionRepository
}

func NewCollectionService(repo interfaces.CollectionRepository) *CollectionService {
	return &CollectionService{repo: repo}
}

func (s *CollectionService) Create(ctx context.Context, name string) error {
	return s.repo.Save(ctx, models.Collection{
		Name:      name,
		Requests:  []models.SavedRequest{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *CollectionService) AddRequest(ctx context.Context, collectionName string, requestName string, req models.Request) error {
	collection, err := s.repo.Get(ctx, collectionName)
	if err != nil {
		collection = models.Collection{Name: collectionName}
	}
	collection.Requests = append(collection.Requests, models.SavedRequest{
		ID:      newID(),
		Name:    requestName,
		Method:  req.Method,
		URL:     req.URL,
		Headers: req.Headers,
		Body:    req.Body,
	})
	return s.repo.Save(ctx, collection)
}

func (s *CollectionService) List(ctx context.Context) ([]models.Collection, error) {
	return s.repo.List(ctx)
}

func (s *CollectionService) Get(ctx context.Context, name string) (models.Collection, error) {
	return s.repo.Get(ctx, name)
}
