package app

import (
	"context"
	"time"

	"reqium/internal/interfaces"
	"reqium/internal/models"
)

type EnvironmentService struct {
	repo interfaces.EnvironmentRepository
}

func NewEnvironmentService(repo interfaces.EnvironmentRepository) *EnvironmentService {
	return &EnvironmentService{repo: repo}
}

func (s *EnvironmentService) Create(ctx context.Context, name string) error {
	return s.repo.Save(ctx, models.Environment{
		Name:      name,
		Variables: map[string]string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *EnvironmentService) Set(ctx context.Context, name string, key string, value string) error {
	env, err := s.repo.Get(ctx, name)
	if err != nil {
		env = models.Environment{Name: name, Variables: map[string]string{}}
	}
	if env.Variables == nil {
		env.Variables = map[string]string{}
	}
	env.Variables[key] = value
	return s.repo.Save(ctx, env)
}

func (s *EnvironmentService) List(ctx context.Context) ([]models.Environment, error) {
	return s.repo.List(ctx)
}

func (s *EnvironmentService) Use(ctx context.Context, name string) error {
	return s.repo.SetActive(ctx, name)
}

func (s *EnvironmentService) Active(ctx context.Context) (models.Environment, error) {
	return s.repo.Active(ctx)
}

func (s *EnvironmentService) Variables(ctx context.Context, name string) (map[string]string, error) {
	if name == "" {
		env, err := s.repo.Active(ctx)
		if err != nil {
			return map[string]string{}, nil
		}
		return env.Variables, nil
	}
	env, err := s.repo.Get(ctx, name)
	if err != nil {
		return nil, err
	}
	return env.Variables, nil
}
