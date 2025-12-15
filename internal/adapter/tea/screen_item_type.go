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
			m.currentScreen = screenEditItem
		case "enter", " ":
			g, _ := m.editItemTypeMenu.GetGroup()
			return m, func() tea.Msg {
				if item, err := m.itemService.SetItemType(m.selectedItem, model.ItemType(g)); err != nil {
					return fmt.Errorf("set item type: %w", err)
				} else {
					return itemUpdatedMsg{item: item}
				}
			}
		}
	case itemUpdatedMsg:
		m.currentScreen = screenEditItem
		m.resetEditItem(msg.item)
	}
	return m, nil
}

func (m *TeaModel) viewItemType() string {
	var s strings.Builder
	s.WriteString("  Тип обращения:\n\n")
	s.WriteString(m.editItemTypeMenu.GenerateMenu())
	return s.String()
}
