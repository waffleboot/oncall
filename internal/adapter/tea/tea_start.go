package tea

import (
	tea "github.com/charmbracelet/bubbletea"
)

type startModel struct {
	controller *controller
}

func NewStartModel(controller *controller) *startModel {
	return &startModel{controller: controller}
}

func (s *startModel) Init() tea.Cmd {
	return nil
}

func (s *startModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s *startModel) View() string {
	return ""
}
