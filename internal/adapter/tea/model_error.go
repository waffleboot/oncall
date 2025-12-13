package tea

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ErrorModel struct {
	message string
	prev    tea.Model
}

func NewErrorModel(message string, prev tea.Model) *ErrorModel {
	return &ErrorModel{message: message, prev: prev}
}

func (m *ErrorModel) Init() tea.Cmd {
	return nil
}

func (m *ErrorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *ErrorModel) View() string {
	return m.message
}
