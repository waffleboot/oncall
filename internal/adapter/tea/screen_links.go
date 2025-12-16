package tea

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateItemLinks(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editItemLinksMenu.ProcessMsg(msg) {
		return m, nil
	}

	newLink := func() tea.Msg {
		link := m.selectedItem.CreateItemLink()
		if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
			return fmt.Errorf("update item link: %w", err)
		}
		return itemLinkCreatedMsg{link: link}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItem
			return m, m.getItem
		case "d":
			if g, p := m.editItemLinksMenu.GetGroup(); g == "links" {
				return m, func() tea.Msg {
					m.selectedItem.DeleteItemLink(m.links[p], time.Now())
					if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item link: %w", err)
					}
					return m.getItem()
				}
			}
		case "n":
			return m, newLink
		case "enter", " ":
			switch g, p := m.editItemLinksMenu.GetGroup(); g {
			case "exit":
				m.currentScreen = screenItem
				return m, m.getItem
			case "new":
				return m, newLink
			case "links":
				m.selectedLink = m.links[p]
				m.currentScreen = screenItemLink
				m.resetItemLink()
			}
		}
	case itemLinkCreatedMsg:
		m.selectedLink = msg.link
		m.currentScreen = screenItemLink
		m.resetItemLink()
	case model.Item:
		m.selectedItem = msg
		m.resetItemLinks("")
	}

	return m, nil
}

func (m *TeaModel) viewItemLinks() string {
	return m.editItemLinksMenu.GenerateMenu()
}

func (m *TeaModel) resetItemLinks(toGroup string) {
	m.links = m.selectedItem.ActiveLinks()
	m.editItemLinksMenu.ResetMenu()
	m.editItemLinksMenu.AddGroup("exit")
	m.editItemLinksMenu.AddGroup("new")
	m.editItemLinksMenu.AddGroupWithItems("links", len(m.links))
	if toGroup != "" {
		m.editItemLinksMenu.JumpToGroup(toGroup)
	} else {
		m.editItemLinksMenu.AdjustCursor()
	}
}
