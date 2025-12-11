package tea

import (
	tea "github.com/charmbracelet/bubbletea"
)

type editModel struct {
	item string
	prev tea.Model
}

func NewEditModel(item string, prev tea.Model) *editModel {
	return &editModel{prev: prev, item: item}
}

func (m *editModel) Init() tea.Cmd {
	return nil
}

func (m *editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.prev, nil
		}
	}
	return m, nil
}

func (m *editModel) View() string {
	return m.item
}
