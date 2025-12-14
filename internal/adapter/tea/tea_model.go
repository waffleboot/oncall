package tea

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

const (
	screenAllItems screen = "all_items"
)

type (
	TeaModelConfig struct {
		ItemService port.ItemService
	}
	TeaModel struct {
		config        TeaModelConfig
		screens       []screen
		items         []model.Item
		allItemsModel allItemsModel
	}
	allItemsModel struct{}
)

func NewTeaModel(config TeaModelConfig) *TeaModel {
	m := &TeaModel{config: config}
	m.screenPush(screenAllItems)
	return m
}

func (m *TeaModel) Init() tea.Cmd {
	return m.getItems
}

func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.currentScreen() {
	case screenAllItems:
		return m.updateAllItems(msg)
	}
	return m, nil
}

func (m *TeaModel) View() string {
	switch m.currentScreen() {
	case screenAllItems:
		return m.viewAllItems()
	}
	return ""
}

func (m *TeaModel) viewAllItems() string {
	var s strings.Builder
	for i := range m.items {
		s.WriteString(fmt.Sprintf("#%d\n", m.items[i].ID))
	}
	return s.String()
}

func (m *TeaModel) updateAllItems(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case []model.Item:
		m.items = msg
	}
	return m, nil
}

func (m TeaModel) getItems() tea.Msg {
	items, err := m.config.ItemService.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}
	return items
}
