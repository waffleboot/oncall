package port

import (
	"github.com/google/uuid"
	"github.com/waffleboot/oncall/internal/model"
)

type ItemService interface {
	GetItem(id uuid.UUID) (model.Item, error)
	GetItems() []model.Item
	CreateItem() (model.Item, error)
	UpdateItem(model.Item) (model.Item, error)
	DeleteItem(model.Item) (model.Item, error)
	AwakeItem(model.Item) (model.Item, error)
	SleepItem(model.Item) (model.Item, error)
	CloseItem(model.Item) (model.Item, error)
}

type JournalService interface {
	CloseJournal() error
}
