package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

const (
	startNew   = "new"
	startItems = "items"
	startClose = "close"
	startPrint = "print"
	startExit  = "exit"
)

type ModelStart struct {
	controller  *Controller
	itemService port.ItemService
	items       []model.Item
	menu        *Menu
}

func NewStartModel(controller *Controller, itemService port.ItemService) *ModelStart {
	m := &ModelStart{controller: controller, itemService: itemService}
	m.menu = NewMenu(func(group string, pos int) string {
		switch {
		case group == startNew:
			return "Новое обращение"
		case group == startItems:
			return m.itemLabel(m.items[pos])
		case group == startClose:
			return "Закрыть журнал"
		case group == startPrint:
			return "Распечатать журнал"
		case group == startExit:
			return "Exit"
		}
		return ""
	})
	return m
}

func (m *ModelStart) Init() tea.Cmd {
	return func() tea.Msg {
		items, err := m.itemService.GetItems()
		if err != nil {
			return fmt.Errorf("get items: %w", err)
		}
		return items
	}
}

func (m *ModelStart) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter", " ":
			switch g, p := m.menu.GetGroup(); g {
			case startNew:
				return m, func() tea.Msg {
					newItem := m.itemService.CreateItem()
					if err := m.itemService.AddItem(newItem); err != nil {
						return fmt.Errorf("add item: %w", err)
					}
					return newItem
				}
			case startItems:
				next := m.controller.modelEdit(m.items[p].ID)
				return next, next.Init()
			case startClose:
				return m.controller.modelCloseJournal(), nil
			case startExit:
				return m, tea.Quit
			}

		}
	case []model.Item:
		m.items = msg
		m.resetMenu()
		return m, nil
	case model.Item:
		next := m.controller.modelEdit(msg.ID)
		return next, next.Init()
	case error:
		return m.controller.modelError(msg.Error(), m), nil
	}
	return m, nil
}

func (m *ModelStart) View() string {
	return m.menu.GenerateMenu()
}

// func (m *ModelStart) resetItems() (tea.Model, tea.Cmd) {
// 	if g, _ := m.menu.GetGroup(); g == "" {
// 		m.menu.JumpToGroup(startNew)
// 		m.menu.MoveCursorUp()
// 	}

// 	return m, nil
// }

// func (m *ModelStart) resetItemsAndJump(itemID int) (tea.Model, tea.Cmd) {

// 	m.menu.JumpToItem(startItems, func(pos int) (found bool) {
// 		return m.items[pos].ID == itemID
// 	})

// 	return m, nil
// }

func (m *ModelStart) resetMenu() {
	m.menu.ResetMenu()

	m.menu.AddGroup(startExit)
	m.menu.AddGroup(startNew)
	m.menu.AddGroup(startClose)
	m.menu.AddGroup(startPrint)
	m.menu.AddGroupWithItems(startItems, len(m.items))

	if g, _ := m.menu.GetGroup(); g == "" {
		m.menu.JumpToGroup(startNew)
		m.menu.MoveCursorUp()
	}
}

func (m *ModelStart) itemLabel(item model.Item) string {
	switch {
	case item.IsSleep():
		return fmt.Sprintf("? #%d - %s", item.ID, item.Type)
	case item.IsClosed():
		return fmt.Sprintf("x #%d - %s", item.ID, item.Type)
	default:
		return fmt.Sprintf("  #%d - %s", item.ID, item.Type)
	}
}
