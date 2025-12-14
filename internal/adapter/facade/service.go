package facade

import (
	"errors"
	"fmt"
	"slices"
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
	i, err := s.getItem(item.ID)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}
	s.items[i].Sleep(time.Now())
	return nil
}

func (s *ItemService) AwakeItem(item model.Item) error {
	i, err := s.getItem(item.ID)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}
	s.items[i].Awake()
	return nil
}

func (s *ItemService) CloseItem(item model.Item) error {
	i, err := s.getItem(item.ID)
	if err != nil {
		return fmt.Errorf("get item: %w", err)
	}
	s.items[i].Close(time.Now())
	return nil
}

func (s *ItemService) DeleteItem(item model.Item) error {
	s.items = slices.DeleteFunc(s.items, func(it model.Item) bool {
		return it.ID == item.ID
	})
	return nil
}

func (s *ItemService) getItem(itemID int) (int, error) {
	for i := range s.items {
		if s.items[i].ID == itemID {
			return i, nil
		}
	}
	return 0, errors.New("not found")
}
