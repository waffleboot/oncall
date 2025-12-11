package port

import "github.com/waffleboot/oncall/internal/model"

type Service interface {
	CreateItem() model.Item
	AddItem(item model.Item) error
	DeleteItem(item model.Item) error
	GetItems() []model.Item
}
