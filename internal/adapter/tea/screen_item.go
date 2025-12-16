package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateItem(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuEditItem.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItems
			return m, m.getItems
		case "enter", " ":
			switch g, _ := m.menuEditItem.GetGroup(); g {
			case "exit":
				m.currentScreen = screenItems
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
				m.currentScreen = screenTitle
				m.resetItemTitle()
			case "item_type":
				m.currentScreen = screenItemType
				m.menuItemType.JumpToGroup(string(m.selectedItem.Type))
			case "item_nodes":
				m.currentScreen = screenNodes
			case "item_notes":
				m.currentScreen = screenNotes
			case "item_links":
				m.currentScreen = screenLinks
				m.resetItemLinks("new")
			case "item_vms":
				m.currentScreen = screenVMs
				m.resetVMs("new")
			}
		case "t":
			m.menuEditItem.JumpToGroup("item_title")
			m.currentScreen = screenTitle
			m.resetItemTitle()
		case "l":
			m.menuEditItem.JumpToGroup("item_links")
			m.currentScreen = screenLinks
			m.resetItemLinks("new")
		case "v":
			m.menuEditItem.JumpToGroup("vms")
			m.currentScreen = screenVMs
			m.resetVMs("new")
		case "s":
			return m.toggleSleep(m.selectedItem)
		}
	case itemClosedMsg:
		m.currentScreen = screenItems
		return m, m.getItems
	case itemDeletedMsg:
		m.currentScreen = screenItems
		return m, m.getItems
	case model.Item:
		m.selectedItem = msg
		m.resetEditItem("")
	}
	return m, nil
}

func (m *TeaModel) viewItem() string {
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
	s.WriteString(m.menuEditItem.GenerateMenu())

	return s.String()
}

func (m *TeaModel) resetEditItem(toGroup string) {
	m.menuEditItem.ResetMenu()

	m.menuEditItem.AddGroup("exit")

	if !m.selectedItem.IsClosed() {
		m.menuEditItem.AddGroup("item_type")
	}

	m.menuEditItem.AddGroup("item_title")
	m.menuEditItem.AddGroup("item_nodes")
	m.menuEditItem.AddGroup("item_vms")
	m.menuEditItem.AddGroup("item_links")
	m.menuEditItem.AddGroup("item_notes")
	m.menuEditItem.AddDelimiter()

	if m.selectedItem.IsActive() {
		m.menuEditItem.AddGroup("sleep")
	}

	if m.selectedItem.IsSleep() {
		m.menuEditItem.AddGroup("awake")
	}

	if !m.selectedItem.IsClosed() {
		m.menuEditItem.AddGroup("close")
	}

	m.menuEditItem.AddDelimiter()
	m.menuEditItem.AddGroup("delete")

	if toGroup != "" {
		m.menuEditItem.JumpToGroup(toGroup)
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
