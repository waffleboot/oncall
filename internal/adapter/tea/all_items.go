package tea

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
)

type allItemsModel struct {
	menu *Menu
}

func (m *TeaModel) updateAllItems(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.allItemsModel.menu.ProcessMsg(msg) {
		return m, nil
	}

	switch msg := msg.(type) {
	case []model.Item:
		m.allItemsModel.menu.ResetMenu()
		m.allItemsModel.menu.AddGroup("exit")
		m.allItemsModel.menu.AddGroup("new")
		m.allItemsModel.menu.AddGroup("close_journal")
		m.allItemsModel.menu.AddGroup("print_journal")
		m.allItemsModel.menu.AddGroupWithItems("items", len(m.items))
		m.allItemsModel.menu.AdjustCursor()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, tea.Quit
		case "1":
			m.selectedItem = 0
			m.screenPush(screenEditItem)
		case "2":
			m.selectedItem = 1
			m.screenPush(screenEditItem)
		}
	}
	return m, nil
}

func (m *TeaModel) viewAllItems() string {
	return m.allItemsModel.menu.GenerateMenu()
}
