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
			if m.cursor < m.menu().maxCursor() {
				m.cursor++
			}
		case "enter", " ":
			switch g, p := m.menu().getGroup(m.cursor); {
			case g == "exit" && p == 0:
				return m.prev, m.prev.Init()
			case g == "delete" && p == 0:
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

	s.WriteString(fmt.Sprintf("  #%d\n\n", m.item.ID))

	s.WriteString(m.addMenu(0, "Выйти"))
	s.WriteString(m.addMenu(1, "Удалить"))
	return s.String()
}

func (m *editModel) addMenu(i int, text string) string {
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

func (m *editModel) menu() menu {
	var n menu
	n.addGroup("exit", 1)
	n.addGroup("delete", 1)
	return n
}
