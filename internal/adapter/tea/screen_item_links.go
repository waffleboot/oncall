package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateItemLinks(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editItemLinksMenu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenEditItem
		case "enter", " ":
			switch g, _ := m.editItemLinksMenu.GetGroup(); g {
			case "new":
				return m, func() tea.Msg {
					link := m.selectedItem.CreateItemLink()
					if err := m.itemService.UpdateItemLink(m.selectedItem, link); err != nil {
						return fmt.Errorf("update item link: %w", err)
					}
					return itemLinkCreatedMsg{link: link}
				}
			}
		}
	case itemLinkCreatedMsg:
		m.selectedLink = msg.link
		m.currentScreen = screenItemLink
	}

	return m, nil
}

func (m *TeaModel) viewItemLinks() string {
	return m.editItemLinksMenu.GenerateMenu()
}

func (m *TeaModel) resetItemLinksMenu() {
	m.links = m.selectedItem.ActiveLinks()
	m.editItemLinksMenu.ResetMenu()
	m.editItemLinksMenu.AddGroup("new")
	m.editItemLinksMenu.AddGroupWithItems("links", len(m.links))
}
