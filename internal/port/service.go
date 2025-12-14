package port

import "github.com/waffleboot/oncall/internal/model"

type ItemService interface {
	GetItems() ([]model.Item, error)
}
