package tea

import tea "github.com/charmbracelet/bubbletea"

func (m *TeaModel) updateItemVMs(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenEditItem
		}
	}

	return m, nil
}

func (m *TeaModel) viewItemVMs() string {
	return "vms\n"
}
