package port

import (
	"github.com/google/uuid"
	"github.com/waffleboot/oncall/internal/model"
)

type Storage interface {
	GetItems() ([]model.Item, error)
	UpdateItem(item model.Item) error
	DeleteItem(itemID uuid.UUID) error

	CloseJournal() error
}
