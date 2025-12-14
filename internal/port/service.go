package port

import "github.com/waffleboot/oncall/internal/model"

type ItemService interface {
	GetItems() ([]model.Item, error)
	CreateItem() model.Item
	UpdateItem(model.Item) error
	CloseItem(model.Item) error
	SleepItem(model.Item) error
	AwakeItem(model.Item) error
	DeleteItem(model.Item) error
}

type JournalService interface {
	CloseJournal() error
}
