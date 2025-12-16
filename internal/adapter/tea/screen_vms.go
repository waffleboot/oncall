package tea

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"go.uber.org/zap"
)

func (m *TeaModel) updateVMs(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.menuVMs.ProcessMsg(msg) {
		return m, nil
	}

	newVM := func() tea.Msg {
		vm := m.selectedItem.CreateVM()
		m.log.Info("create vm", zap.Int("id", vm.ID))
		if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
			return fmt.Errorf("update item: %w", err)
		}
		return vm
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.currentScreen = screenItem
			return m, m.getItem
		case "d":
			if g, p := m.menuVMs.GetGroup(); g == "vms" {
				return m, func() tea.Msg {
					m.selectedItem.DeleteVM(m.vms[p], time.Now())
					if err := m.itemService.UpdateItem(m.selectedItem); err != nil {
						return fmt.Errorf("update item: %w", err)
					}
					return m.getItem()
				}
			}
		case "n":
			return m, newVM
		case "enter", " ":
			switch g, p := m.menuVMs.GetGroup(); g {
			case "exit":
				m.currentScreen = screenItem
				return m, m.getItem
			case "new":
				return m, newVM
			case "vms":
				return m, func() tea.Msg { return m.vms[p] }
			}
		}
	case model.Item:
		m.selectedItem = msg
		m.resetVMs("")
	case model.VM:
		m.selectedVM = msg
		m.currentScreen = screenVM
		m.resetVM()
	}

	return m, nil
}

func (m *TeaModel) viewVMs() string {
	return m.menuVMs.GenerateMenu()
}

func (m *TeaModel) resetVMs(toGroup string) {
	m.vms = m.selectedItem.ActiveVMs()
	m.menuVMs.ResetMenu()
	m.menuVMs.AddGroup("exit")
	m.menuVMs.AddGroup("new")
	m.menuVMs.AddGroupWithItems("vms", len(m.vms))
	if toGroup != "" {
		m.menuVMs.JumpToGroup(toGroup)
	} else {
		m.menuVMs.AdjustCursor()
	}
}
