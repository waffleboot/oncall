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
		startModel        func() (*StartModel, error)
		editModel         func(item model.Item, prev Prev) *EditModel
		errorModel        func(message string, prev Prev) *ErrorModel
		closeJournalModel func(prev Prev) *CloseJournalModel
		itemTypeModel     func(item model.Item, prev Prev) *ItemTypeModel
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
	start, err := c.startModel()
	if err != nil {
		return fmt.Errorf("start model: %w", err)
	}

	p := tea.NewProgram(start)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("tea run: %w", err)
	}
	return nil
}

func WithService(itemService port.ItemService, journalService port.JournalService, log *zap.Logger) func(c *Controller) {
	return func(controller *Controller) {
		controller.startModel = func() (*StartModel, error) {
			return NewStartModel(controller, itemService)
		}
		controller.editModel = func(item model.Item, prev Prev) *EditModel {
			return NewEditModel(controller, itemService, item, prev)
		}
		controller.errorModel = NewErrorModel
		controller.closeJournalModel = func(prev Prev) *CloseJournalModel {
			return NewCloseJournalModel(controller, journalService, prev)
		}
		controller.itemTypeModel = func(item model.Item, prev Prev) *ItemTypeModel {
			return NewItemTypeModel(controller, itemService, item, prev)
		}
	}
}
