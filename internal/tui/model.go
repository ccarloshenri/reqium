package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"

	"reqium/internal/models"
)

type tab int

const (
	tabHistory tab = iota
	tabCollections
	tabEnvironments
)

type model struct {
	activeTab    tab
	history      []models.HistoryEntry
	collections  []models.Collection
	environments []models.Environment
	err          error
}

func newModel(ctx context.Context, services Services) (model, error) {
	history, err := services.History.List(ctx, 20)
	if err != nil {
		return model{}, err
	}
	collections, err := services.Collections.List(ctx)
	if err != nil {
		return model{}, err
	}
	environments, err := services.Environments.List(ctx)
	if err != nil {
		return model{}, err
	}
	return model{
		activeTab:    tabHistory,
		history:      history,
		collections:  collections,
		environments: environments,
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}
