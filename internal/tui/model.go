package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"reqium/internal/enums"
	"reqium/internal/models"
)

type mode int

const (
	modeDashboard mode = iota
	modeRequestForm
	modeEnvForm
)

type tab int

const (
	tabHistory tab = iota
	tabCollections
	tabEnvironments
)

type model struct {
	ctx          context.Context
	services     Services
	width        int
	height       int
	mode         mode
	activeTab    tab
	history      []models.HistoryEntry
	collections  []models.Collection
	environments []models.Environment
	requestForm  requestForm
	envForm      envForm
	response     string
	status       string
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
		ctx:          ctx,
		services:     services,
		width:        112,
		height:       36,
		mode:         modeDashboard,
		activeTab:    tabHistory,
		history:      history,
		collections:  collections,
		environments: environments,
		requestForm:  newRequestForm(),
		envForm:      newEnvForm(),
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

type requestForm struct {
	methodIndex int
	focus       int
	url         textinput.Model
	env         textinput.Model
	headers     textarea.Model
	body        textarea.Model
}

func newRequestForm() requestForm {
	url := textinput.New()
	url.Placeholder = "https://api.example.com/users or {{base_url}}/users"
	url.Prompt = ""
	url.CharLimit = 500
	url.Width = 72
	url.Focus()

	env := textinput.New()
	env.Placeholder = "dev, staging, prod or blank for active"
	env.Prompt = ""
	env.CharLimit = 80
	env.Width = 40

	headers := textarea.New()
	headers.Placeholder = "Content-Type: application/json\nAuthorization: Bearer {{token}}"
	headers.CharLimit = 2000
	headers.SetWidth(72)
	headers.SetHeight(4)

	body := textarea.New()
	body.Placeholder = `{"name":"John"}`
	body.CharLimit = 10000
	body.SetWidth(72)
	body.SetHeight(6)

	return requestForm{
		methodIndex: 0,
		focus:       1,
		url:         url,
		env:         env,
		headers:     headers,
		body:        body,
	}
}

func (f requestForm) method() string {
	methods := []enums.HTTPMethod{
		enums.MethodGET,
		enums.MethodPOST,
		enums.MethodPUT,
		enums.MethodPATCH,
		enums.MethodDELETE,
	}
	return methods[f.methodIndex].String()
}

type envForm struct {
	focus textinputFocus
	env   textinput.Model
	key   textinput.Model
	value textinput.Model
}

type textinputFocus int

func newEnvForm() envForm {
	env := textinput.New()
	env.Placeholder = "dev"
	env.Prompt = ""
	env.CharLimit = 80
	env.Width = 36
	env.Focus()

	key := textinput.New()
	key.Placeholder = "base_url"
	key.Prompt = ""
	key.CharLimit = 120
	key.Width = 36

	value := textinput.New()
	value.Placeholder = "https://api.example.com"
	value.Prompt = ""
	value.CharLimit = 500
	value.Width = 72

	return envForm{env: env, key: key, value: value}
}

func (m model) reload() (model, error) {
	history, err := m.services.History.List(m.ctx, 20)
	if err != nil {
		return m, err
	}
	collections, err := m.services.Collections.List(m.ctx)
	if err != nil {
		return m, err
	}
	environments, err := m.services.Environments.List(m.ctx)
	if err != nil {
		return m, err
	}
	m.history = history
	m.collections = collections
	m.environments = environments
	return m, nil
}

func (m model) activeEnvironmentName() string {
	for _, env := range m.environments {
		if env.Active {
			return env.Name
		}
	}
	return "none"
}

type requestSentMsg struct {
	output string
	err    error
}

type envSavedMsg struct {
	err error
}

type refreshedMsg struct {
	model model
	err   error
}
