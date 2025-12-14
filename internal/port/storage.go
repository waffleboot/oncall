package port

import "github.com/waffleboot/oncall/internal/model"

type Storage interface {
	GetItems() ([]model.Item, error)
	UpdateItem(item model.Item) error
	DeleteItem(itemID int) error

	CloseJournal() error
}
