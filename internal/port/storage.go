package port

import "github.com/waffleboot/oncall/internal/model"

type Storage interface {
	GetItem(itemID int) (model.Item, error)
	UpdateItem(item model.Item) error
	DeleteItem(itemID int) error
	GetItems() ([]model.Item, error)

	CloseJournal() error
}
