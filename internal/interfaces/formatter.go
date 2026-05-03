package interfaces

import "reqium/internal/domain"

type Formatter interface {
	Format(response domain.Response) (string, error)
}
