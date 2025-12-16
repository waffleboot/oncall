package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

func (m *TeaModel) updateItemType(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.editItemTypeMenu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItem
		case "enter", " ":
			g, _ := m.editItemTypeMenu.GetGroup()
			return m, func() tea.Msg {
				m.selectedItem.Type = model.ItemType(g)
				if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
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
	s.WriteString(m.editItemTypeMenu.GenerateMenu())
	return s.String()
}
