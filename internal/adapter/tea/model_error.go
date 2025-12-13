package tea

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ModelError struct {
	message string
	prev    tea.Model
}

func NewModelError(message string, prev tea.Model) *ModelError {
	return &ModelError{message: message, prev: prev}
}

func (m *ModelError) Init() tea.Cmd {
	return nil
}

func (m *ModelError) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.prev, nil
		case "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *ModelError) View() string {
	return m.message
}
