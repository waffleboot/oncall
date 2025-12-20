package tea

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateItems(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuAllItems.Update(msg) {
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
		case "esc", "q", "e":
			return m.exitScreen()
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if len(m.items) > 0 {
				i, _ := strconv.Atoi(msg.String())
				if i > len(m.items) {
					break
				}
				m.menuAllItems.JumpToItem("items", func(pos int) (found bool) {
					return pos == i-1
				})
				m.currentScreen = screenItem
				m.selectedItem = m.items[i-1]
				m.resetEditItem("exit")
			}
		case "enter", " ":
			switch g, p := m.menuAllItems.GetGroup(); g {
			case "exit":
				return m.exitScreen()
			case "new":
				return m, newItem
			case "close_journal":
				return m, func() tea.Msg {
					if err := m.journalService.CloseJournal(); err != nil {
						return fmt.Errorf("close journal: %w", err)
					}
					return "exit"
				}
			case "print_journal":
				m.printJournal = true
				return m.exitScreen()
			case "next":
				m.currentScreen = screenUsers
				m.resetUsers(m.userService.GetUser())
				return m, nil
			case "items":
				m.currentScreen = screenItem
				m.selectedItem = m.items[p]
				m.resetEditItem("exit")
			}
		case "n":
			return m, newItem
		case "p":
			m.printJournal = true
			return m, tea.Quit
		case "s":
			if g, p := m.menuAllItems.GetGroup(); g == "items" {
				return m.toggleSleep(m.items[p])
			}
		}
	case string:
		if msg == "exit" {
			return m, tea.Quit
		}
	case []model.Item:
		m.resetItems(msg)
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

func (m *TeaModel) viewItems() string {
	return m.menuAllItems.View()
}

func (m *TeaModel) resetItems(items []model.Item) {
	m.items = items
	m.menuAllItems.ResetMenu()
	m.menuAllItems.AddGroup("exit")
	m.menuAllItems.AddGroup("new")
	m.menuAllItems.AddGroup("close_journal")
	m.menuAllItems.AddGroup("print_journal")
	m.menuAllItems.AddGroup("next")
	m.menuAllItems.AddGroupWithItems("items", len(m.items))
	if len(m.items) == 0 {
		m.menuAllItems.JumpToGroup("new")
	} else {
		m.menuAllItems.AdjustCursor()
	}
}
