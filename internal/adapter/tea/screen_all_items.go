package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateAllItems(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.allItemsMenu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, tea.Quit
		case "enter", " ":
			switch g, p := m.allItemsMenu.GetGroup(); g {
			case "exit":
				return m, tea.Quit
			case "new":
				return m, func() tea.Msg {
					item := m.itemService.CreateItem()
					if err := m.itemService.UpdateItem(item); err != nil {
						return fmt.Errorf("create item: %w", err)
					}

					return itemCreatedMsg{item: item}
				}
			case "close_journal":
				return m, func() tea.Msg {
					if err := m.journalService.CloseJournal(); err != nil {
						return fmt.Errorf("close journal: %w", err)
					}
					return m.getItems()
				}
			case "print_journal":
			case "items":
				m.selectedItem = p
				m.currentScreen = screenEditItem
				m.editItemMenu.JumpToGroup("exit")
				m.resetEditItemMenu()
			}
		}
	}
	return m, nil
}

func (m *TeaModel) viewAllItems() string {
	return m.allItemsMenu.GenerateMenu()
}

func (m *TeaModel) resetAllItemsMenu() {
	m.allItemsMenu.ResetMenu()
	m.allItemsMenu.AddGroup("exit")
	m.allItemsMenu.AddGroup("new")
	m.allItemsMenu.AddGroup("close_journal")
	m.allItemsMenu.AddGroup("print_journal")
	m.allItemsMenu.AddGroupWithItems("items", len(m.items))
	if len(m.items) > 0 {
		m.allItemsMenu.AdjustCursor()
	} else {
		m.allItemsMenu.JumpToGroup("new")
	}
}

func (m *TeaModel) itemLabel(item model.Item) string {
	switch {
	case item.IsSleep():
		return fmt.Sprintf("? #%d - %s", item.ID, item.Type)
	case item.IsClosed():
		return fmt.Sprintf("x #%d - %s", item.ID, item.Type)
	default:
		return fmt.Sprintf("  #%d - %s", item.ID, item.Type)
	}
}
