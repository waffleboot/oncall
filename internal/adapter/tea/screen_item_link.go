package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateItemLink(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItemLinks
		}
	}

	return m, nil
}

func (m *TeaModel) viewItemLink() string {
	return fmt.Sprintf("%d\n", m.selectedLink.ID)
}

func (m *TeaModel) resetItemLink(link model.ItemLink) {
	m.selectedLink = link
}
