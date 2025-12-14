package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateEditItem(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.screenPop()
			return m, nil
		}
	}
	return m, nil
}

func (m *TeaModel) viewEditItem() string {
	return fmt.Sprintf("#%d\n", m.items[m.selectedItem].ID)
}
