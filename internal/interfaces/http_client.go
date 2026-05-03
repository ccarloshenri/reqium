package interfaces

import (
	"context"

	"reqium/internal/domain"
)

type HTTPClient interface {
	Do(ctx context.Context, req domain.Request) (domain.Response, error)
}
