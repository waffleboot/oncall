package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateConsoleLogs(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuConsoleLogs.Update(msg) {
		return m, nil
	}

	newConsoleLog := func() tea.Msg {
		consoleLog := m.selectedItem.CreateConsoleLog()
		if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
			return fmt.Errorf("update item: %w", err)
		}
		return consoleLogCreatedMsg{consoleLog: consoleLog}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItem
			return m, m.getItem
		case "d":
			if g, p := m.menuConsoleLogs.GetGroup(); g == "console_logs" {
				return m, func() tea.Msg {
					m.selectedItem.DeleteConsoleLog(m.consoleLogs[p])
					if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return m.getItem()
				}
			}
		case "n":
			return m, newConsoleLog
		case "enter", " ":
			switch g, p := m.menuConsoleLogs.GetGroup(); g {
			case "exit":
				m.currentScreen = screenItem
				return m, m.getItem
			case "new":
				return m, newConsoleLog
			case "console_logs":
				m.selectedConsoleLog = m.consoleLogs[p]
				m.currentScreen = screenConsoleLog
				m.resetConsoleLog()
			}
		}
	case consoleLogCreatedMsg:
		m.selectedConsoleLog = msg.consoleLog
		m.currentScreen = screenConsoleLog
		m.resetConsoleLog()
	case model.Item:
		m.selectedItem = msg
		m.resetConsoleLogs("")
	}

	return m, nil
}

func (m *TeaModel) viewConsoleLogs() string {
	return m.menuConsoleLogs.View()
}

func (m *TeaModel) resetConsoleLogs(toGroup string) {
	m.consoleLogs = m.selectedItem.ActiveConsoleLogs()
	m.menuConsoleLogs.ResetMenu()
	m.menuConsoleLogs.AddGroup("exit")
	m.menuConsoleLogs.AddGroup("new")
	m.menuConsoleLogs.AddGroupWithItems("console_logs", len(m.consoleLogs))
	if toGroup != "" {
		m.menuConsoleLogs.JumpToGroup(toGroup)
	} else {
		m.menuConsoleLogs.AdjustCursor()
	}
}
