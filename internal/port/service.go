package port

import "github.com/waffleboot/oncall/internal/model"

type ItemService interface {
	CreateItem() model.Item
	UpdateItem(model.Item) error
	GetItems() ([]model.Item, error)
	SleepItem(model.Item) error
	AwakeItem(model.Item) error
}
