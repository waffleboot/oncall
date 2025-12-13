package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
	"go.uber.org/zap"
)

const (
	startNew   = "new"
	startItems = "items"
	startClose = "close"
	startPrint = "print"
	startExit  = "exit"
)

type (
	ModelStart struct {
		controller  *Controller
		itemService port.ItemService
		items       []model.Item
		menu        *Menu
		log         *zap.Logger
	}
)

func NewModelStart(controller *Controller, itemService port.ItemService, log *zap.Logger) *ModelStart {
	m := &ModelStart{controller: controller, itemService: itemService, log: log}
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
	m.resetMenu()
	m.menu.JumpToGroup(startNew)
	return m
}

func (m *ModelStart) Init() tea.Cmd {
	return m.getItems
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
					item := m.itemService.CreateItem()
					if err := m.itemService.AddItem(item); err != nil {
						return fmt.Errorf("add item: %w", err)
					}

					return item
				}
			case startItems:
				next := m.controller.modelEdit(m.items[p], m)
				return next, next.Init()
			case startClose:
				return m.controller.modelCloseJournal(m), nil
			case startExit:
				return m, tea.Quit
			}
		}
	case []model.Item:
		m.items = msg
		m.resetMenu()
		return m, nil
	case model.Item:
		m.items = append(m.items, msg)
		m.resetMenu()
		m.menu.JumpToItem(startItems, func(pos int) (found bool) {
			return m.items[pos].ID == msg.ID
		})
		return m, nil
	case error:
		return m.controller.modelError(msg.Error(), m), nil
	}
	return m, m.getItems
}

func (m *ModelStart) View() string {
	return m.menu.GenerateMenu()
}

func (m *ModelStart) resetMenu() {
	m.menu.ResetMenu()
	m.menu.AddGroup(startExit)
	m.menu.AddGroup(startNew)
	m.menu.AddGroup(startClose)
	m.menu.AddGroup(startPrint)
	m.menu.AddGroupWithItems(startItems, len(m.items))
	m.menu.AdjustCursor()
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

func (m *ModelStart) getItems() tea.Msg {
	items, err := m.itemService.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}
	return items
}
