package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type editModel struct {
	controller *controller
	service    port.Service
	item       model.Item
	prev       tea.Model
	cursor     int
}

func NewEditModel(controller *controller, service port.Service, item model.Item, prev tea.Model) *editModel {
	return &editModel{controller: controller, service: service, prev: prev, item: item}
}

func (m *editModel) Init() tea.Cmd {
	m.cursor = 0
	return nil
}

func (m *editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m.prev, nil
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < 1 {
				m.cursor++
			}
		case "enter", " ":
			if m.cursor == 0 {
				return m.prev, m.prev.Init()
			}
			if m.cursor == 1 {
				if err := m.service.DeleteItem(m.item); err != nil {
					return m.controller.errorModel(err.Error(), m.prev), nil
				}
				return m.prev, m.prev.Init()
			}
		}
	}
	return m, nil
}

func (m *editModel) View() string {
	var s strings.Builder

	s.WriteString("  ")
	s.WriteString(fmt.Sprintf("#%d", m.item.ID))
	s.WriteString("\n")

	s.WriteString(m.menu(0, "Выйти"))
	s.WriteString(m.menu(1, "Удалить"))
	return s.String()
}

func (m *editModel) menu(i int, text string) string {
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
