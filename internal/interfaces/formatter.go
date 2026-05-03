package interfaces

import "reqium/internal/models"

type Formatter interface {
	Format(response models.Response) (string, error)
}
