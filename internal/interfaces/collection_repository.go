package interfaces

import (
	"context"

	"reqium/internal/models"
)

type CollectionRepository interface {
	Save(ctx context.Context, collection models.Collection) error
	List(ctx context.Context) ([]models.Collection, error)
	Get(ctx context.Context, name string) (models.Collection, error)
	Delete(ctx context.Context, name string) error
}
