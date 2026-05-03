package interfaces

import "reqium/internal/models"

type VariableResolver interface {
	Resolve(input string, variables map[string]string) (string, error)
	ResolveRequest(req models.Request, variables map[string]string) (models.Request, error)
}
