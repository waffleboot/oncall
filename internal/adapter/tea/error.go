package tea

import (
	tea "github.com/charmbracelet/bubbletea"
)

type errorModel struct {
	message string
	prev    tea.Model
}

func NewErrorModel(message string, prev tea.Model) *errorModel {
	return &errorModel{message: message, prev: prev}
}

func (m *errorModel) Init() tea.Cmd {
	return nil
}

func (m *errorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.prev, m.prev.Init()
		}
	}
	return m, nil
}

func (m *errorModel) View() string {
	return m.message
}
