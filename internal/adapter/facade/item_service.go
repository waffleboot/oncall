package facade

import (
	"fmt"
	"time"

	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type ItemService struct {
	storage port.Storage
	idGen   port.IDGenerator
}

func NewItemService(storage port.Storage, idGen port.IDGenerator) *ItemService {
	return &ItemService{storage: storage, idGen: idGen}
}

func (s *ItemService) CreateItem() model.Item {
	return model.Item{
		ID:   s.idGen.GenerateID(),
		Type: model.ItemTypeAsk,
	}
}

func (s *ItemService) UpdateItem(item model.Item) error {
	if err := s.storage.UpdateItem(item); err != nil {
		return fmt.Errorf("update item: %w", err)
	}
	return nil
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
