package port

import "github.com/waffleboot/oncall/internal/model"

type Service interface {
	AddItem(item model.Item) error
	DeleteItem(item model.Item) error
	GetItems() []model.Item
}

type ItemBuilder interface {
	CreateItem() model.Item
}
