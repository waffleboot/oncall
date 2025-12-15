package tea

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateItemLinks(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenEditItem
		case "enter", " ":
			switch g, _ := m.editItemLinksMenu.GetGroup(); g {
			case "new":

			}
		}
	}

	return m, nil
}

func (m *TeaModel) viewItemLinks() string {
	return m.editItemLinksMenu.GenerateMenu()
}

func (m *TeaModel) resetItemLinksMenu() {
	links := m.selectedItem.LiveLinks()
	m.editItemLinksMenu.ResetMenu()
	m.editItemLinksMenu.AddGroup("new")
	m.editItemLinksMenu.AddGroupWithItems("links", len(links))
}
