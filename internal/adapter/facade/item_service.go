package facade

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
	"go.uber.org/zap"
)

type ItemService struct {
	numGen  port.NumGenerator
	storage port.Storage
	journal model.Journal
	log     *zap.Logger
}

func NewItemService(storage port.Storage, numGen port.NumGenerator, log *zap.Logger) (*ItemService, error) {
	journal, err := storage.GetJournal()
	if err != nil {
		return nil, fmt.Errorf("get journal: %w", err)
	}
	return &ItemService{storage: storage, journal: journal, numGen: numGen, log: log}, nil
}

func (s *ItemService) CreateItem() (model.Item, error) {
	num, err := s.numGen.GenerateNum()
	if err != nil {
		return model.Item{}, fmt.Errorf("generate num: %w", err)
	}

	item := s.journal.CreateItem(num)

	if err := s.saveJournal(); err != nil {
		return model.Item{}, err
	}

	s.log.Debug("item created", zap.Int("num", num))

	return item, nil
}

func (s *ItemService) UpdateItem(item model.Item) (model.Item, error) {
	return s.updateItem(item)
}

func (s *ItemService) DeleteItem(item model.Item) (model.Item, error) {
	item.Delete()
	return s.updateItem(item)
}

func (s *ItemService) GetItem(id uuid.UUID) (model.Item, error) {
	return s.journal.GetItem(id)
}

func (s *ItemService) GetItems() []model.Item {
	items := make([]model.Item, 0, len(s.journal.Items))
	for _, item := range s.journal.Items {
		if !item.IsDeleted() {
			items = append(items, item)
		}
	}
	return items
}

func (s *ItemService) CloseItem(item model.Item) (model.Item, error) {
	item.Close()
	return s.updateItem(item)
}

func (s *ItemService) SleepItem(item model.Item) (model.Item, error) {
	item.Sleep()
	return s.updateItem(item)
}

func (s *ItemService) AwakeItem(item model.Item) (model.Item, error) {
	item.Awake()
	return s.updateItem(item)
}

func (s *ItemService) CloseJournal() error {
	if err := s.storage.CloseJournal(s.journal); err != nil {
		return fmt.Errorf("close journal: %w", err)
	}
	s.journal = model.NewJournal()
	return nil
}

func (s *ItemService) updateItem(item model.Item) (_ model.Item, err error) {
	item, err = s.journal.UpdateItem(item)
	if err != nil {
		return model.Item{}, fmt.Errorf("update item: %w", err)
	}
	if err := s.saveJournal(); err != nil {
		return model.Item{}, err
	}
	return item, nil
}

func (s *ItemService) saveJournal() error {
	if err := s.storage.SaveJournal(s.journal); err != nil {
		return fmt.Errorf("save journal: %w", err)
	}
	s.log.Debug("journal saved", zap.Int("items_count", len(s.journal.Items)))
	return nil
}
