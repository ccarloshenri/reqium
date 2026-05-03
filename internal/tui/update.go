package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"reqium/internal/models"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.mode == modeRequestForm {
			return m.updateRequestForm(msg)
		}
		if m.mode == modeEnvForm {
			return m.updateEnvForm(msg)
		}

		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "n":
			m.requestForm = newRequestForm()
			m.mode = modeRequestForm
			m.status = ""
			m.err = nil
			return m, nil
		case "v":
			m.envForm = newEnvForm()
			m.mode = modeEnvForm
			m.status = ""
			m.err = nil
			return m, nil
		case "r":
			return m, m.refreshCmd()
		case "h", "left":
			if m.activeTab == tabHistory {
				m.activeTab = tabEnvironments
			} else {
				m.activeTab--
			}
		case "l", "right", "tab":
			if m.activeTab == tabEnvironments {
				m.activeTab = tabHistory
			} else {
				m.activeTab++
			}
		case "1":
			m.activeTab = tabHistory
		case "2":
			m.activeTab = tabCollections
		case "3":
			m.activeTab = tabEnvironments
		}
	case requestSentMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = ""
			return m, nil
		}
		m.response = msg.output
		m.status = "Request completed and saved to history."
		m.err = nil
		m.mode = modeDashboard
		m.activeTab = tabHistory
		return m, m.refreshCmd()
	case envSavedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.status = ""
			return m, nil
		}
		m.status = "Environment variable saved."
		m.err = nil
		m.mode = modeDashboard
		m.activeTab = tabEnvironments
		return m, m.refreshCmd()
	case refreshedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		return msg.model, nil
	}
	return m, nil
}

func (m model) updateRequestForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeDashboard
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "ctrl+s":
		return m, m.sendRequestCmd()
	case "ctrl+space":
		m.requestForm = m.applyVariableCompletion(m.requestForm)
		return m, nil
	case "tab", "ctrl+n":
		m.requestForm = focusRequestField(m.requestForm, false)
		return m, nil
	case "shift+tab", "ctrl+p":
		m.requestForm = focusRequestField(m.requestForm, true)
		return m, nil
	case "ctrl+right":
		m.requestForm = cycleMethod(m.requestForm, false)
		return m, nil
	case "ctrl+left":
		m.requestForm = cycleMethod(m.requestForm, true)
		return m, nil
	case "left", "right":
		if m.requestForm.focus == 0 {
			m.requestForm = cycleMethod(m.requestForm, msg.String() == "left")
			return m, nil
		}
	}

	var cmd tea.Cmd
	switch m.requestForm.focus {
	case 1:
		m.requestForm.url, cmd = m.requestForm.url.Update(msg)
	case 2:
		m.requestForm.env, cmd = m.requestForm.env.Update(msg)
	case 3:
		m.requestForm.headers, cmd = m.requestForm.headers.Update(msg)
	case 4:
		m.requestForm.body, cmd = m.requestForm.body.Update(msg)
	}
	return m, cmd
}

func (m model) updateEnvForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modeDashboard
		return m, nil
	case "ctrl+c":
		return m, tea.Quit
	case "ctrl+s", "enter":
		return m, m.saveEnvCmd()
	case "tab", "ctrl+n":
		m.envForm = focusEnvField(m.envForm, false)
		return m, nil
	case "shift+tab", "ctrl+p":
		m.envForm = focusEnvField(m.envForm, true)
		return m, nil
	}

	var cmd tea.Cmd
	switch m.envForm.focus {
	case 0:
		m.envForm.env, cmd = m.envForm.env.Update(msg)
	case 1:
		m.envForm.key, cmd = m.envForm.key.Update(msg)
	case 2:
		m.envForm.value, cmd = m.envForm.value.Update(msg)
	}
	return m, cmd
}

func (m model) applyVariableCompletion(form requestForm) requestForm {
	suggestions := m.variableSuggestions()
	if len(suggestions) == 0 {
		return form
	}
	selected := suggestions[0]

	switch form.focus {
	case 1:
		form.url.SetValue(completeVariable(form.url.Value(), selected))
		form.url.CursorEnd()
	case 3:
		form.headers.SetValue(completeVariable(form.headers.Value(), selected))
		form.headers.CursorEnd()
	case 4:
		form.body.SetValue(completeVariable(form.body.Value(), selected))
		form.body.CursorEnd()
	}
	return form
}

func completeVariable(value string, variable string) string {
	start := strings.LastIndex(value, "{{")
	if start == -1 {
		return value
	}
	replacement := "{{" + variable + "}}"
	end := strings.Index(value[start:], "}}")
	if end == -1 {
		return value[:start] + replacement
	}
	end = start + end + len("}}")
	return value[:start] + replacement + value[end:]
}

func focusRequestField(form requestForm, backwards bool) requestForm {
	clearRequestFocus(&form)
	if backwards {
		form.focus--
		if form.focus < 0 {
			form.focus = 4
		}
	} else {
		form.focus++
		if form.focus > 4 {
			form.focus = 0
		}
	}
	applyRequestFocus(&form)
	return form
}

func clearRequestFocus(form *requestForm) {
	form.url.Blur()
	form.env.Blur()
	form.headers.Blur()
	form.body.Blur()
}

func applyRequestFocus(form *requestForm) {
	switch form.focus {
	case 1:
		form.url.Focus()
	case 2:
		form.env.Focus()
	case 3:
		form.headers.Focus()
	case 4:
		form.body.Focus()
	}
}

func cycleMethod(form requestForm, backwards bool) requestForm {
	if backwards {
		form.methodIndex--
		if form.methodIndex < 0 {
			form.methodIndex = 4
		}
		return form
	}
	form.methodIndex++
	if form.methodIndex > 4 {
		form.methodIndex = 0
	}
	return form
}

func focusEnvField(form envForm, backwards bool) envForm {
	form.env.Blur()
	form.key.Blur()
	form.value.Blur()
	if backwards {
		form.focus--
		if form.focus < 0 {
			form.focus = 2
		}
	} else {
		form.focus++
		if form.focus > 2 {
			form.focus = 0
		}
	}
	switch form.focus {
	case 0:
		form.env.Focus()
	case 1:
		form.key.Focus()
	case 2:
		form.value.Focus()
	}
	return form
}

func (m model) sendRequestCmd() tea.Cmd {
	form := m.requestForm
	return func() tea.Msg {
		headers, err := parseHeaderBlock(form.headers.Value())
		if err != nil {
			return requestSentMsg{err: err}
		}

		req := models.Request{
			Method:  form.method(),
			URL:     strings.TrimSpace(form.url.Value()),
			Headers: headers,
			Body:    []byte(strings.TrimSpace(form.body.Value())),
			Timeout: 30 * time.Second,
		}
		if strings.TrimSpace(form.body.Value()) == "" {
			req.Body = nil
		}

		variables, err := environmentVariables(m)
		if err != nil {
			return requestSentMsg{err: err}
		}
		req, err = m.services.Resolver.ResolveRequest(req, variables)
		if err != nil {
			return requestSentMsg{err: err}
		}

		response, err := m.services.Requests.Send(m.ctx, req)
		if err != nil {
			return requestSentMsg{err: err}
		}
		output, err := m.services.Formatter.Format(response)
		return requestSentMsg{output: output, err: err}
	}
}

func environmentVariables(m model) (map[string]string, error) {
	envName := strings.TrimSpace(m.requestForm.env.Value())
	if envName != "" {
		env, err := m.services.Environments.Get(m.ctx, envName)
		if err != nil {
			return nil, err
		}
		return env.Variables, nil
	}
	env, err := m.services.Environments.Active(m.ctx)
	if err != nil {
		return map[string]string{}, nil
	}
	return env.Variables, nil
}

func (m model) saveEnvCmd() tea.Cmd {
	form := m.envForm
	return func() tea.Msg {
		env := strings.TrimSpace(form.env.Value())
		key := strings.TrimSpace(form.key.Value())
		value := strings.TrimSpace(form.value.Value())
		if env == "" || key == "" {
			return envSavedMsg{err: fmt.Errorf("environment and key are required")}
		}
		if err := m.services.EnvService.Set(m.ctx, env, key, value); err != nil {
			return envSavedMsg{err: err}
		}
		_ = m.services.EnvService.Use(m.ctx, env)
		return envSavedMsg{}
	}
}

func (m model) refreshCmd() tea.Cmd {
	return func() tea.Msg {
		next, err := m.reload()
		return refreshedMsg{model: next, err: err}
	}
}

func parseHeaderBlock(input string) (map[string]string, error) {
	headers := map[string]string{}
	normalized := strings.ReplaceAll(input, ";", "\n")
	for _, line := range strings.Split(normalized, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok || strings.TrimSpace(key) == "" {
			return nil, fmt.Errorf("invalid header %q: expected Key: Value", line)
		}
		headers[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	return headers, nil
}
