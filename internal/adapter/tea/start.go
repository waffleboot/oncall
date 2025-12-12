package tea

import (
	"fmt"

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
	return m.menu().generateMenu(m.cursor)
}

func (m *startModel) menu() menu {
	items := m.service.GetItems()

	var n menu

	n.labelGen = func(group string, pos int) string {
		switch {
		case group == "new" && pos == 0:
			return fmt.Sprintf("Новое обращение %d %d", m.cursor, n.maxCursor())
		case group == "items":
			return fmt.Sprintf("#%d", items[pos].ID)
		}
		return "xwx"
	}

	n.addGroup("new", 1)
	n.addGroup("items", len(items))
	return n
}
