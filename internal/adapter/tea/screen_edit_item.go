package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateEditItem(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editItemMenu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenAllItems
			return m, m.getItems
		case "enter", " ":
			switch g, _ := m.editItemMenu.GetGroup(); g {
			case "exit":
				m.currentScreen = screenAllItems
				return m, m.getItems
			case "sleep":
				return m, func() tea.Msg {
					if item, err := m.itemService.SleepItem(m.selectedItem); err != nil {
						return fmt.Errorf("sleep item: %w", err)
					} else {
						return item
					}
				}
			case "awake":
				return m, func() tea.Msg {
					if item, err := m.itemService.AwakeItem(m.selectedItem); err != nil {
						return fmt.Errorf("awake item: %w", err)
					} else {
						return item
					}
				}
			case "close":
				return m, func() tea.Msg {
					if err := m.itemService.CloseItem(m.selectedItem); err != nil {
						return fmt.Errorf("close item: %w", err)
					}
					return itemClosedMsg{}
				}
			case "delete":
				return m, func() tea.Msg {
					if err := m.itemService.DeleteItem(m.selectedItem); err != nil {
						return fmt.Errorf("delete item: %w", err)
					}
					return itemDeletedMsg{}
				}
			case "item_title":
				m.currentScreen = screenItemTitle
				m.resetItemTitle()
			case "item_type":
				m.currentScreen = screenItemType
				m.editItemTypeMenu.JumpToGroup(string(m.selectedItem.Type))
			case "item_nodes":
				m.currentScreen = screenItemNodes
			case "item_notes":
				m.currentScreen = screenItemNotes
			case "item_links":
				m.currentScreen = screenItemLinks
				m.resetItemLinks("new")
			case "item_vms":
				m.currentScreen = screenItemVMs
			}
		case "s":
			return m.toggleSleep(m.selectedItem)
		}
	case itemClosedMsg:
		m.currentScreen = screenAllItems
		return m, m.getItems
	case itemDeletedMsg:
		m.currentScreen = screenAllItems
		return m, m.getItems
	case model.Item:
		m.selectedItem = msg
		m.resetEditItem("")
	}
	return m, nil
}

func (m *TeaModel) viewEditItem() string {
	var state string

	switch {
	case m.selectedItem.IsSleep():
		state = " - в ожидании"
	case m.selectedItem.IsClosed():
		switch m.selectedItem.Type {
		case model.ItemTypeAsk:
			state = " - закрыто"
		default:
			state = " - закрыт"
		}
	}

	var s strings.Builder
	s.WriteString(fmt.Sprintf("  #%d - %s - %s%s\n\n", m.selectedItem.Num, m.selectedItem.Type, m.selectedItem.TitleForView(), state))
	s.WriteString(m.editItemMenu.GenerateMenu())

	return s.String()
}

func (m *TeaModel) resetEditItem(toGroup string) {
	m.editItemMenu.ResetMenu()

	m.editItemMenu.AddGroup("exit")

	if !m.selectedItem.IsClosed() {
		m.editItemMenu.AddGroup("item_type")
	}

	m.editItemMenu.AddGroup("item_title")
	m.editItemMenu.AddGroup("item_nodes")
	m.editItemMenu.AddGroup("item_vms")
	m.editItemMenu.AddGroup("item_links")
	m.editItemMenu.AddGroup("item_notes")
	m.editItemMenu.AddDelimiter()

	if m.selectedItem.IsActive() {
		m.editItemMenu.AddGroup("sleep")
	}

	if m.selectedItem.IsSleep() {
		m.editItemMenu.AddGroup("awake")
	}

	if !m.selectedItem.IsClosed() {
		m.editItemMenu.AddGroup("close")
	}

	m.editItemMenu.AddDelimiter()
	m.editItemMenu.AddGroup("delete")

	if toGroup != "" {
		m.editItemMenu.JumpToGroup(toGroup)
	}
}

func (m *TeaModel) toggleSleep(item model.Item) (tea.Model, tea.Cmd) {
	if item.IsSleep() {
		return m, func() tea.Msg {
			if item, err := m.itemService.AwakeItem(item); err != nil {
				return fmt.Errorf("awake: %w", err)
			} else {
				return item
			}
		}
	} else {
		return m, func() tea.Msg {
			if item, err := m.itemService.SleepItem(item); err != nil {
				return fmt.Errorf("sleep: %w", err)
			} else {
				return item
			}
		}
	}
}
