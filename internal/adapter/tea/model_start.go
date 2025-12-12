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

type StartModel struct {
	controller  *Controller
	itemService port.ItemService
	items       []model.Item
	menu        *Menu
}

func NewStartModel(
	controller *Controller,
	itemService port.ItemService,
) (*StartModel, error) {
	m := &StartModel{
		controller:  controller,
		itemService: itemService,
	}

	m.menu = NewMenu(func(group string, pos int) string {
		switch {
		case group == startNew && pos == 0:
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

	items, err := itemService.GetItems()
	if err != nil {
		return nil, fmt.Errorf("get items: %w", err)
	}

	m.resetMenu(items)

	return m, nil
}

func (m *StartModel) Init() tea.Cmd {
	return nil
}

func (m *StartModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				newItem := m.itemService.CreateItem()

				if err := m.itemService.AddItem(newItem); err != nil {
					return m.controller.errorModel(err.Error(), m.resetItems), nil
				}

				next := m.controller.editModel(newItem, func() (tea.Model, tea.Cmd) {
					return m.resetItemsAndJump(newItem.ID)
				})

				return next, next.Init()
			case startItems:
				next := m.controller.editModel(m.items[p], m.resetItems)
				return next, next.Init()
			case startClose:
				next := m.controller.closeJournalModel(m.resetItems)
				return next, next.Init()
			case startExit:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *StartModel) View() string {
	return m.menu.GenerateMenu()
}

func (m *StartModel) resetItems() (tea.Model, tea.Cmd) {
	items, err := m.itemService.GetItems()
	if err != nil {
		return m.controller.errorModel(err.Error(), m.resetItems), nil
	}

	m.resetMenu(items)

	if g, _ := m.menu.GetGroup(); g == "" {
		m.menu.JumpToGroup(startNew)
		m.menu.MoveCursorUp()
	}

	return m, nil
}

func (m *StartModel) resetItemsAndJump(itemID int) (tea.Model, tea.Cmd) {
	items, err := m.itemService.GetItems()
	if err != nil {
		return m.controller.errorModel(err.Error(), m.resetItems), nil
	}

	m.resetMenu(items)

	m.menu.JumpToItem(startItems, func(pos int) (found bool) {
		return m.items[pos].ID == itemID
	})

	return m, nil
}

func (m *StartModel) resetMenu(items []model.Item) {
	m.items = items
	m.menu.ResetMenu()
	m.menu.AddGroup(startExit)
	m.menu.AddGroup(startNew)
	m.menu.AddGroup(startClose)
	m.menu.AddGroup(startPrint)
	m.menu.AddGroupWithItems(startItems, len(items))
}

func (m *StartModel) itemLabel(item model.Item) string {
	switch {
	case item.IsSleep():
		return fmt.Sprintf("? #%d - %s", item.ID, item.Type)
	case item.IsClosed():
		return fmt.Sprintf("x #%d - %s", item.ID, item.Type)
	default:
		return fmt.Sprintf("  #%d - %s", item.ID, item.Type)
	}
}
