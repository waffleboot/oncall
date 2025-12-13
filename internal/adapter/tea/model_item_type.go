package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type ModelItemType struct {
	controller  *Controller
	itemService port.ItemService
	item        model.Item
	menu        *Menu
}

func NewModelItemType(
	controller *Controller,
	itemService port.ItemService,
	item model.Item,
) *ModelItemType {
	m := &ModelItemType{
		controller:  controller,
		itemService: itemService,
		item:        item,
	}

	m.menu = NewMenu(func(group string, pos int) string {
		switch model.ItemType(group) {
		case model.ItemTypeInc:
			return "Инцидент"
		case model.ItemTypeAdhoc:
			return "Adhoc"
		case model.ItemTypeAsk:
			return "Обращение"
		case model.ItemTypeAlert:
			return "Alert"
		}
		return ""
	})

	m.menu.AddGroup(string(model.ItemTypeInc))
	m.menu.AddGroup(string(model.ItemTypeAdhoc))
	m.menu.AddGroup(string(model.ItemTypeAsk))
	m.menu.AddGroup(string(model.ItemTypeAlert))
	m.menu.JumpToGroup(string(item.Type))

	return m
}

func (m *ModelItemType) Init() tea.Cmd {
	return nil
}

func (m *ModelItemType) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			next := m.controller.modelEdit(m.item.ID)
			return next, next.Init()
		case "enter", " ":
			g, _ := m.menu.GetGroup()
			return m, func() tea.Msg {
				if err := m.itemService.SetItemType(m.item, model.ItemType(g)); err != nil {
					return fmt.Errorf("set item type: %w", err)
				}
				return "done"
			}
		}
	case error:
		return m.controller.modelError(msg.Error(), m), nil
	case string:
		if msg == "done" {
			next := m.controller.modelEdit(m.item.ID)
			return next, next.Init()
		}
	}
	return m, nil
}

func (m *ModelItemType) View() string {
	var s strings.Builder

	s.WriteString("  Тип обращения:\n\n")
	s.WriteString(m.menu.GenerateMenu())

	return s.String()
}
