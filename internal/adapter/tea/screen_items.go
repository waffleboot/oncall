package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateAllItems(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuAllItems.ProcessMsg(msg) {
		return m, nil
	}

	newItem := func() tea.Msg {
		if item, err := m.itemService.CreateItem(); err != nil {
			return fmt.Errorf("create item: %w", err)
		} else {
			return itemCreatedMsg{item: item}
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			return m, tea.Quit
		case "enter", " ":
			switch g, p := m.menuAllItems.GetGroup(); g {
			case "exit":
				return m, tea.Quit
			case "new":
				return m, newItem
			case "close_journal":
				return m, func() tea.Msg {
					if err := m.journalService.CloseJournal(); err != nil {
						return fmt.Errorf("close journal: %w", err)
					}
					return m.getItems()
				}
			case "print_journal":
				m.printJournal = true
				return m, tea.Quit
			case "items":
				m.selectedItem = m.items[p]
				m.resetEditItem("exit")
				m.currentScreen = screenItem
			}
		case "n":
			return m, newItem
		case "s":
			if g, p := m.menuAllItems.GetGroup(); g == "items" {
				return m.toggleSleep(m.items[p])
			}
		}
	case []model.Item:
		m.resetAllItems(msg)
	case model.Item:
		return m, m.getItems
	case itemCreatedMsg:
		m.selectedItem = msg.item
		m.resetEditItem("exit")
		m.currentScreen = screenItemType
		m.menuItemType.JumpToGroup(string(m.selectedItem.Type))
		return m, nil
	}
	return m, nil
}

func (m *TeaModel) viewAllItems() string {
	return m.menuAllItems.GenerateMenu()
}

func (m *TeaModel) resetAllItems(items []model.Item) {
	m.items = items
	m.menuAllItems.ResetMenu()
	m.menuAllItems.AddGroup("exit")
	m.menuAllItems.AddGroup("new")
	m.menuAllItems.AddGroup("close_journal")
	m.menuAllItems.AddGroup("print_journal")
	m.menuAllItems.AddGroupWithItems("items", len(m.items))
	if len(m.items) == 0 {
		m.menuAllItems.JumpToGroup("new")
	} else {
		m.menuAllItems.AdjustCursor()
	}
}
