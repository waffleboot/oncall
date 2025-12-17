package tea

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateLinks(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.linksMenu.ProcessMsg(msg) {
		return m, nil
	}

	newLink := func() tea.Msg {
		link := m.selectedItem.CreateLink()
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
			if g, p := m.linksMenu.GetGroup(); g == "links" {
				return m, func() tea.Msg {
					m.selectedItem.DeleteLink(m.links[p], time.Now())
					if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return m.getItem()
				}
			}
		case "n":
			return m, newLink
		case "p":
			if g, p := m.linksMenu.GetGroup(); g == "links" {
				return m, func() tea.Msg {
					link := m.links[p]
					link.Public = !link.Public
					m.selectedItem.UpdateLink(link)
					if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return m.getItem()
				}
			}
		case "enter", " ":
			switch g, p := m.linksMenu.GetGroup(); g {
			case "exit":
				m.currentScreen = screenItem
				return m, m.getItem
			case "new":
				return m, newLink
			case "links":
				m.selectedLink = m.links[p]
				m.currentScreen = screenLink
				m.resetItemLink()
			}
		}
	case itemLinkCreatedMsg:
		m.selectedLink = msg.link
		m.currentScreen = screenLink
		m.resetItemLink()
	case model.Item:
		m.selectedItem = msg
		m.resetLinks("")
	}

	return m, nil
}

func (m *TeaModel) viewLinks() string {
	return m.linksMenu.GenerateMenu()
}

func (m *TeaModel) resetLinks(toGroup string) {
	m.links = m.selectedItem.ActiveLinks()
	m.linksMenu.ResetMenu()
	m.linksMenu.AddGroup("exit")
	m.linksMenu.AddGroup("new")
	m.linksMenu.AddGroupWithItems("links", len(m.links))
	if toGroup != "" {
		m.linksMenu.JumpToGroup(toGroup)
	} else {
		m.linksMenu.AdjustCursor()
	}
}
