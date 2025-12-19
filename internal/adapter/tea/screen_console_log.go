package tea

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateConsoleLog(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m.exitScreen()
		case "q":
			return m.exitScreen()
		}
	case string:
		if msg == "exit" {
			m.currentScreen = screenConsoleLogs
			return m, m.getItem
		}

	}

	return m, nil
}

func (m *TeaModel) viewConsoleLog() string {
	return ""
}

func (m *TeaModel) resetConsoleLog() {
}
