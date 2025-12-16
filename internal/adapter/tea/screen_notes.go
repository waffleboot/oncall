package tea

import tea "github.com/charmbracelet/bubbletea"

func (m *TeaModel) updateNotes(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItem
		}
	}

	return m, nil
}

func (m *TeaModel) viewNotes() string {
	return "notes\n"
}
