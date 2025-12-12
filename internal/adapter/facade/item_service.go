package facade

import (
	"fmt"
	"time"

	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type ItemService struct {
	storage port.Storage
	gen     port.IDGenerator
}

func NewItemService(storage port.Storage, gen port.IDGenerator) *ItemService {
	return &ItemService{storage: storage, gen: gen}
}

func (s *ItemService) CreateItem() model.Item {
	return model.Item{
		ID:   s.gen.GenerateID(),
		Type: model.ItemTypeAsk,
	}
}

func (s *ItemService) AddItem(item model.Item) error {
	if err := s.storage.AddItem(item); err != nil {
		return fmt.Errorf("create item: %w", err)
	}

	return nil
}

func (s *ItemService) GetItem(itemID int) (model.Item, error) {
	item, err := s.storage.GetItem(itemID)
	if err != nil {
		return model.Item{}, fmt.Errorf("get item: %w", err)
	}

	return item, nil
}

func (s *ItemService) DeleteItem(item model.Item) error {
	if err := s.storage.DeleteItem(item.ID); err != nil {
		return fmt.Errorf("delete item: %w", err)
	}
	return nil
}

func (s *ItemService) SleepItem(item model.Item) error {
	item.Sleep(time.Now())

	if err := s.storage.UpdateItem(item); err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

func (s *ItemService) AwakeItem(item model.Item) error {
	item.Awake()

	if err := s.storage.UpdateItem(item); err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

func (s *ItemService) CloseItem(item model.Item) error {
	item.Close(time.Now())

	if err := s.storage.UpdateItem(item); err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}

func (s *ItemService) GetItems() ([]model.Item, error) {
	return s.storage.GetItems()
}

func (s *ItemService) SetItemType(item model.Item, itemType model.ItemType) error {
	item.Type = itemType

	if err := s.storage.UpdateItem(item); err != nil {
		return fmt.Errorf("update item: %w", err)
	}

	return nil
}
