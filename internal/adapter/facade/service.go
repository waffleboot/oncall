package facade

import (
	"fmt"
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type ItemService struct {
	items []model.Item
}

func NewItemService() *ItemService {
	return &ItemService{items: []model.Item{{ID: 1}, {ID: 2}}}
}

func (s *ItemService) CreateItem() model.Item {
	return model.Item{ID: len(s.items) + 1}
}

func (s *ItemService) UpdateItem(item model.Item) error {
	for i := range s.items {
		if s.items[i].ID == item.ID {
			s.items[i] = item
			return nil
		}
	}
	s.items = append(s.items, item)
	return nil
}

func (s *ItemService) GetItems() ([]model.Item, error) {
	return s.items, nil
}

func (s *ItemService) SleepItem(item model.Item) error {
	for i := range s.items {
		if s.items[i].ID == item.ID {
			s.items[i].Sleep(time.Now())
			return nil
		}
	}
	return fmt.Errorf("not found")
}

func (s *ItemService) AwakeItem(item model.Item) error {
	for i := range s.items {
		if s.items[i].ID == item.ID {
			s.items[i].Awake()
			return nil
		}
	}
	return fmt.Errorf("not found")
}

func (s *ItemService) CloseItem(item model.Item) error {
	return nil
}

func (s *ItemService) DeleteItem(item model.Item) error {
	return nil
}
