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
	builder    port.ItemBuilder
	cursor     int
}

func NewStartModel(controller *controller, service port.Service, builder port.ItemBuilder) *startModel {
	return &startModel{controller: controller, service: service, builder: builder}
}

func (m *startModel) Init() tea.Cmd {
	m.cursor = 0
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
			if m.cursor < m.menu().maxCursor() {
				m.cursor++
			}
		case "enter", " ":
			switch g, p := m.menu().getGroup(m.cursor); {
			case g == "new" && p == 0:
				newItem := m.builder.CreateItem()

				if err := m.service.AddItem(newItem); err != nil {
					return m.controller.errorModel(err.Error(), m), nil
				}

				next := m.controller.editModel(newItem, m)

				return next, next.Init()
			case g == "items":
				next := m.controller.editModel(m.service.GetItems()[p], m)
				return next, next.Init()
			}
		}
	}
	return m, nil
}

func (m *startModel) View() string {
	var s strings.Builder

	s.WriteString(m.addMenu(0, "Новое обращение"))

	for i, item := range m.service.GetItems() {
		s.WriteString(m.addMenu(i+1, fmt.Sprintf("#%d", item.ID)))
	}

	return s.String()
}

func (m *startModel) addMenu(i int, text string) string {
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

func (m *startModel) menu() menu {
	var n menu
	n.addGroup("new", 1)
	n.addGroup("items", len(m.service.GetItems()))
	return n
}
