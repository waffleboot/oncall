package facade

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

var _ port.ItemService = (*ItemService)(nil)

type ItemService struct {
	storage port.Storage
	numGen  port.NumGenerator
}

func NewItemService(storage port.Storage, numGen port.NumGenerator) *ItemService {
	return &ItemService{storage: storage, numGen: numGen}
}

func (s *ItemService) CreateItem() (model.Item, error) {
	item := model.Item{
		ID:   uuid.New(),
		Num:  s.numGen.GenerateNum(),
		Type: model.ItemTypeAsk,
	}
	if err := s.updateItem(item); err != nil {
		return model.Item{}, err
	}
	return item, nil
}

func (s *ItemService) UpdateItem(item model.Item) error {
	return s.updateItem(item)
}

func (s *ItemService) GetItem(id uuid.UUID) (model.Item, error) {
	return s.storage.GetItem(id)
}

func (s *ItemService) GetItems() ([]model.Item, error) {
	return s.storage.GetItems()
}

func (s *ItemService) SleepItem(item model.Item) (model.Item, error) {
	item.Sleep(time.Now())
	if err := s.storage.UpdateItem(item); err != nil {
		return model.Item{}, fmt.Errorf("update item: %w", err)
	}
	return item, nil
}

func (s *ItemService) AwakeItem(item model.Item) (model.Item, error) {
	item.Awake()
	if err := s.storage.UpdateItem(item); err != nil {
		return model.Item{}, fmt.Errorf("update item: %w", err)
	}
	return item, nil
}

func (s *ItemService) CloseItem(item model.Item) error {
	item.Close(time.Now())
	if err := s.storage.UpdateItem(item); err != nil {
		return fmt.Errorf("update item: %w", err)
	}
	return nil
}

func (s *ItemService) DeleteItem(item model.Item) error {
	if err := s.storage.DeleteItem(item.ID); err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}

func (s *ItemService) updateItem(item model.Item) error {
	if err := s.storage.UpdateItem(item); err != nil {
		return fmt.Errorf("update item: %w", err)
	}
	return nil
}

func (s *ItemService) UpdateItemLink(item model.Item, link model.ItemLink) error {
	var found bool
	for i := range item.Links {
		if item.Links[i].ID == link.ID {
			item.Links[i] = link
			found = true
			break
		}
	}

	if !found {
		item.Links = append(item.Links, link)
	}

	if err := s.storage.UpdateItem(item); err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}
