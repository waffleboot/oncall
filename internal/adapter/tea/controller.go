package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type (
	controller struct {
		startModel func() *startModel
		editModel  func(item model.Item, prev tea.Model) *editModel
		errorModel func(message string, prev tea.Model) *errorModel
	}
	option func(*controller)
)

func NewController(opts ...option) *controller {
	c := &controller{}
	for i := range opts {
		opts[i](c)
	}
	return c
}

func (c *controller) Run() error {
	p := tea.NewProgram(c.startModel())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tea run: %w", err)
	}
	return nil
}

func WithService(service port.Service, builder port.ItemBuilder) func(c *controller) {
	return func(controller *controller) {
		controller.startModel = func() *startModel {
			return NewStartModel(controller, service, builder)
		}
		controller.editModel = func(item model.Item, prev tea.Model) *editModel {
			return NewEditModel(controller, service, item, prev)
		}
		controller.errorModel = NewErrorModel
	}
}
