package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/port"
)

type startModel struct {
	controller *controller
	service    port.Service
	cursor     int
}

func NewStartModel(controller *controller, service port.Service) *startModel {
	return &startModel{controller: controller, service: service}
}

func (m *startModel) Init() tea.Cmd {
	return nil
}

func (m *startModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
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
				items := m.service.Items()

				newItem := fmt.Sprintf("Item %d", len(items)+1)

				if err := m.service.AddItem(newItem); err != nil {
					return m.controller.errorModel(err.Error(), m), nil
				}

				items = m.service.Items()

				for i := range items {
					if items[i] == newItem {
						m.cursor = i + 1
						break
					}
				}

			} else {
				i := m.cursor - 1
				items := m.service.Items()
				if i < len(items) {
					return m.controller.editModel(items[i], m), nil
				}

			}
		}
	}
	return m, nil
}

func (m *startModel) View() string {
	var s strings.Builder

	items := m.service.Items()

	if len(items) == 0 {
		m.cursor = 0
	} else if m.cursor > len(items) {
		m.cursor = len(items)
	}

	s.WriteString(m.menu(0, "Новая запись"))

	for i, item := range items {
		s.WriteString(m.menu(i+1, item))
	}

	return s.String()
}

func (m *startModel) menu(i int, text string) string {
	var s strings.Builder
	if m.cursor == i {
		s.WriteString("> ")
	} else {
		s.WriteString("  ")
	}
	s.WriteString(text)
	s.WriteString("\n")
	return s.String()
}
