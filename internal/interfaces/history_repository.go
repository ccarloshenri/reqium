package interfaces

import (
	"context"

	"reqium/internal/models"
)

type HistoryRepository interface {
	Save(ctx context.Context, entry models.HistoryEntry) error
	List(ctx context.Context, limit int) ([]models.HistoryEntry, error)
	Get(ctx context.Context, id string) (models.HistoryEntry, error)
}
