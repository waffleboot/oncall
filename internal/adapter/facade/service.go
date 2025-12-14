package facade

import "github.com/waffleboot/oncall/internal/model"

type ItemService struct {
}

func NewItemService() *ItemService {
	return &ItemService{}
}

func (s *ItemService) GetItems() ([]model.Item, error) {
	return []model.Item{{ID: 1}, {ID: 2}}, nil
}
