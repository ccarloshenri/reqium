package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"reqium/internal/app"
	"reqium/internal/interfaces"
)

type Services struct {
	History      interfaces.HistoryRepository
	Environments interfaces.EnvironmentRepository
	Collections  interfaces.CollectionRepository
	Requests     *app.RequestService
	EnvService   *app.EnvironmentService
	Resolver     interfaces.VariableResolver
	Formatter    interfaces.Formatter
}

func Run(ctx context.Context, services Services) error {
	model, err := newModel(ctx, services)
	if err != nil {
		return err
	}
	_, err = tea.NewProgram(model).Run()
	return err
}
