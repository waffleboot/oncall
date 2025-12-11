package port

import "github.com/waffleboot/oncall/internal/model"

type Storage interface {
	GenerateID() int
	AddItem(item model.Item) error
	DeleteItem(item model.Item) error
	GetItems() []model.Item
}
