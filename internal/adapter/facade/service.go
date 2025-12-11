package facade

import (
	"fmt"

	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type Service struct {
	storage port.Storage
}

func NewService(storage port.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) CreateItem() model.Item {
	return model.Item{ID: s.storage.GenerateID()}
}

func (s *Service) AddItem(item model.Item) error {
	if err := s.storage.AddItem(item); err != nil {
		return fmt.Errorf("create item: %w", err)
	}

	return nil
}

func (s *Service) DeleteItem(item model.Item) error {
	if err := s.storage.DeleteItem(item); err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}

func (s *Service) GetItems() []model.Item {
	return s.storage.GetItems()
}
