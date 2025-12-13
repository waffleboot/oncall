package tea

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
	"go.uber.org/zap"
)

type (
	Prev       func() (tea.Model, tea.Cmd)
	Controller struct {
		modelStart        func() *ModelStart
		modelEdit         func(item model.Item, next tea.Model) *ModelEdit
		modelCloseJournal func(next tea.Model) *ModelCloseJournal
		modelItemType     func(item model.Item, next tea.Model) *ModelItemType
		modelError        func(message string, next tea.Model) *ModelError
	}
	option func(*Controller)
)

func NewController(opts ...option) *Controller {
	c := &Controller{}
	for i := range opts {
		opts[i](c)
	}
	return c
}

func (c *Controller) Run() error {
	p := tea.NewProgram(c.modelStart())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tea run: %w", err)
	}
	return nil
}

func WithService(itemService port.ItemService, journalService port.JournalService, log *zap.Logger) func(c *Controller) {
	return func(controller *Controller) {
		controller.modelStart = func() *ModelStart {
			return NewModelStart(controller, itemService, log.With(zap.String("model", "start")))
		}
		controller.modelEdit = func(item model.Item, next tea.Model) *ModelEdit {
			return NewModelEdit(controller, itemService, item, next)
		}
		controller.modelError = NewModelError
		controller.modelCloseJournal = func(next tea.Model) *ModelCloseJournal {
			return NewModelCloseJournal(controller, journalService, next)
		}
		controller.modelItemType = func(item model.Item, next tea.Model) *ModelItemType {
			return NewModelItemType(controller, itemService, item, next)
		}
	}
}
