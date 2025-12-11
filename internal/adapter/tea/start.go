package tea

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/port"
)

type startModel struct {
	controller *controller
	service    port.OnCallService
	cursor     int
}

func NewStartModel(controller *controller, service port.OnCallService) *startModel {
	return &startModel{controller: controller, service: service}
}

func (m *startModel) Init() tea.Cmd {
	return nil
}

func (m *startModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.service.Items()) {
				m.cursor++
			}
		case "enter", " ":
			if m.cursor == 0 {
				m.service.AddItem()
			} else {
				i := m.cursor - 1
				items := m.service.Items()
				if i < len(items) {
					return m.controller.editModel(m, items[i]), nil
				}

			}
		}
	}
	return m, nil
}

func (m *startModel) View() string {
	var s strings.Builder

	if m.cursor == 0 {
		s.WriteString("> Новая запись\n")
	} else {
		s.WriteString("  Новая запись\n")
	}

	for i, item := range m.service.Items() {
		if m.cursor == i+1 {
			s.WriteString("> ")
		} else {
			s.WriteString("  ")
		}
		s.WriteString(item)
		s.WriteString("\n")
	}

	return s.String()
}
