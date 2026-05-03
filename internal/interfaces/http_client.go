package interfaces

import (
	"context"

	"reqium/internal/models"
)

type HTTPClient interface {
	Do(ctx context.Context, req models.Request) (models.Response, error)
}
