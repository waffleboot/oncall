package tea

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *TeaModel) updateItemLink(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItemLinks
			return m, m.getItem
		case "enter":
			return m, func() tea.Msg {
				m.selectedLink.Link = m.linkInput.Value()
				m.selectedItem.UpdateItemLink(m.selectedLink)
				if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
					return fmt.Errorf("update item: %w", err)
				}
				m.currentScreen = screenItemLinks
				return m.getItem()
			}
		}
	}

	m.linkInput, cmd = m.linkInput.Update(msg)
	return m, cmd
}

func (m *TeaModel) viewItemLink() string {
	return m.linkInput.View()
}

func (m *TeaModel) resetItemLink() {
	m.linkInput = textinput.New()
	m.linkInput.Placeholder = "link"
	m.linkInput.Focus()
	m.linkInput.CharLimit = 256
	m.linkInput.Width = 256
	m.linkInput.SetValue(m.selectedLink.Link)
}
