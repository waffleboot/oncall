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

func (s *Service) AddItem() error {
	items := s.storage.Items()
	if err := s.storage.AddItem(fmt.Sprintf("item %d", len(items)+1)); err != nil {
		return fmt.Errorf("add item: %w", err)
	}

	return nil
}
