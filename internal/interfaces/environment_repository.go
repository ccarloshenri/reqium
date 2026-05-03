package interfaces

import (
	"context"

	"reqium/internal/models"
)

type EnvironmentRepository interface {
	Save(ctx context.Context, env models.Environment) error
	List(ctx context.Context) ([]models.Environment, error)
	Get(ctx context.Context, name string) (models.Environment, error)
	Delete(ctx context.Context, name string) error
	SetActive(ctx context.Context, name string) error
	Active(ctx context.Context) (models.Environment, error)
}
