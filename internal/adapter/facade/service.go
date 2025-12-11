package facade

import (
	"fmt"

	"github.com/waffleboot/oncall/internal/port"
)

type Service struct {
	storage port.Storage
}

func NewService(storage port.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Items() []string {
	return s.storage.Items()
}

func (s *Service) AddItem(item string) error {
	if err := s.storage.AddItem(item); err != nil {
		return fmt.Errorf("add item: %w", err)
	}

	return nil
}

func (s *Service) DeleteItem(item string) error {
	if err := s.storage.DeleteItem(item); err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}
