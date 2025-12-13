package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type ItemTypeModel struct {
	controller  *Controller
	itemService port.ItemService
	item        model.Item
	menu        *Menu
	prev        Prev
}

func NewItemTypeModel(
	controller *Controller,
	itemService port.ItemService,
	item model.Item,
) *ItemTypeModel {
	m := &ItemTypeModel{
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

func (m *ItemTypeModel) Init() tea.Cmd {
	return nil
}

func (m *ItemTypeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m.prev()
		case "enter", " ":
			g, _ := m.menu.GetGroup()
			return m, func() tea.Msg {
				if err := m.itemService.SetItemType(m.item, model.ItemType(g)); err != nil {
					return fmt.Errorf("set item type: %w", err)
				}
				return m.item
			}
		}
	case error:
		return m.controller.errorModel(msg.Error(), m), nil
	case model.Item:
		next := m.controller.editModel(msg)
		return next, next.Init()
	}
	return m, nil
}

func (m *ItemTypeModel) View() string {
	var s strings.Builder

	s.WriteString("  Тип обращения:\n\n")
	s.WriteString(m.menu.GenerateMenu())

	return s.String()
}
