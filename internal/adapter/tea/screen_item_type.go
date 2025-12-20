package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateItemType(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuItemType.Update(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "e":
			m.currentScreen = screenItem
		case "a", "1":
			m.menuItemType.JumpToGroup(string(model.ItemTypeAsk))
		case "i", "2":
			m.menuItemType.JumpToGroup(string(model.ItemTypeInc))
		case "r", "3":
			m.menuItemType.JumpToGroup(string(model.ItemTypeAlert))
		case "d", "4":
			m.menuItemType.JumpToGroup(string(model.ItemTypeAdhoc))
		case "enter", " ":
			g, _ := m.menuItemType.GetGroup()
			return m, func() tea.Msg {
				m.selectedItem.Type = model.ItemType(g)
				if _, err := m.itemService.UpdateItem(m.selectedItem); err != nil {
					return fmt.Errorf("update item: %w", err)
				} else {
					return itemUpdatedMsg{}
				}
			}
		}
	case itemUpdatedMsg:
		m.currentScreen = screenItem
		return m, m.getItem
	}
	return m, nil
}

func (m *TeaModel) viewItemType() string {
	var s strings.Builder
	s.WriteString("  Тип обращения:\n\n")
	s.WriteString(m.menuItemType.View())
	return s.String()
}
