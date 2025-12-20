package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateItem(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuItem.Update(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItems
			return m, m.getItems
		case "enter", " ":
			switch g, _ := m.menuItem.GetGroup(); g {
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
					if _, err := m.itemService.CloseItem(m.selectedItem); err != nil {
						return fmt.Errorf("close item: %w", err)
					}
					return itemClosedMsg{}
				}
			case "delete":
				return m, func() tea.Msg {
					if _, err := m.itemService.DeleteItem(m.selectedItem); err != nil {
						return fmt.Errorf("delete item: %w", err)
					}
					return itemDeletedMsg{}
				}
			case string(screenTitle):
				m.currentScreen = screenTitle
				m.resetTitle()
			case string(screenItemType):
				m.currentScreen = screenItemType
				m.menuItemType.JumpToGroup(string(m.selectedItem.Type))
			case string(screenNodes):
				m.currentScreen = screenNodes
				m.resetNodes("new")
			case string(screenNotes):
				m.currentScreen = screenNotes
				m.resetNotes("new")
			case string(screenLinks):
				m.currentScreen = screenLinks
				m.resetLinks("new")
			case string(screenVMs):
				m.currentScreen = screenVMs
				m.resetVMs("new")
			case string(screenConsoleLogs):
				m.currentScreen = screenConsoleLogs
				m.resetConsoleLogs("new")
			}
		case "t":
			m.menuItem.JumpToGroup(string(screenTitle))
			m.currentScreen = screenTitle
			m.resetTitle()
		case "l":
			m.menuItem.JumpToGroup(string(screenLinks))
			m.currentScreen = screenLinks
			m.resetLinks("new")
		case "v":
			m.menuItem.JumpToGroup(string(screenVMs))
			m.currentScreen = screenVMs
			m.resetVMs("new")
		case "h":
			m.menuItem.JumpToGroup(string(screenNodes))
			m.currentScreen = screenNodes
			m.resetNodes("new")
		case "n":
			m.menuItem.JumpToGroup(string(screenNotes))
			m.currentScreen = screenNotes
			m.resetNotes("new")
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
	s.WriteString(fmt.Sprintf("  #%d - %s - %s%s\n\n", m.selectedItem.ID, m.selectedItem.Type, m.selectedItem.MenuItem(), state))
	s.WriteString(m.menuItem.View())

	return s.String()
}

func (m *TeaModel) resetEditItem(toGroup string) {
	m.menuItem.ResetMenu()

	m.menuItem.AddGroup("exit")

	if !m.selectedItem.IsClosed() {
		m.menuItem.AddGroup(string(screenItemType))
	}

	m.menuItem.AddGroup(string(screenTitle))
	m.menuItem.AddGroup(string(screenConsoleLogs))
	m.menuItem.AddGroup(string(screenNodes))
	m.menuItem.AddGroup(string(screenVMs))
	m.menuItem.AddGroup(string(screenLinks))
	m.menuItem.AddGroup(string(screenNotes))
	m.menuItem.AddDelimiter()

	var needDelimiter bool

	if m.selectedItem.InProgress() {
		m.menuItem.AddGroup("sleep")
		needDelimiter = true
	}

	if m.selectedItem.IsSleep() {
		m.menuItem.AddGroup("awake")
		needDelimiter = true
	}

	if !m.selectedItem.IsClosed() {
		m.menuItem.AddGroup("close")
		needDelimiter = true
	}

	if needDelimiter {
		m.menuItem.AddDelimiter()
	}
	m.menuItem.AddGroup("delete")

	if toGroup != "" {
		m.menuItem.JumpToGroup(toGroup)
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
