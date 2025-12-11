package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"

	"github.com/waffleboot/oncall/internal/model"
)

type (
	Config struct {
		Filename string
	}
	Storage struct {
		config Config
		lastID int
		items  []model.Item
	}
	storedData struct {
		LastID int          `json:"last_id"`
		Items  []storedItem `json:"items"`
	}
	storedItem struct {
		ID int `json:"id"`
	}
)

func NewStorage(config Config) (*Storage, error) {
	s := &Storage{config: config}
	if err := s.loadData(); err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}
	return s, nil
}

func (s *Storage) GenerateID() int {
	s.lastID++
	return s.lastID
}

func (s *Storage) AddItem(newItem model.Item) error {
	storedItems := make([]storedItem, len(s.items)+1)
	for i := range s.items {
		storedItems[i].fromDomain(s.items[i])
	}

	storedItems[len(storedItems)-1].fromDomain(newItem)

	if err := s.saveData(storedData{
		LastID: s.lastID,
		Items:  storedItems,
	}); err != nil {
		return fmt.Errorf("save items: %w", err)
	}

	s.items = append(s.items, newItem)

	return nil
}

func (s *Storage) DeleteItem(item model.Item) error {
	newItems := slices.DeleteFunc(s.items, func(it model.Item) bool {
		return it.ID == item.ID
	})

	storedItems := make([]storedItem, len(newItems))
	for i := range newItems {
		storedItems[i].fromDomain(newItems[i])
	}

	if err := s.saveData(storedData{
		LastID: s.lastID,
		Items:  storedItems,
	}); err != nil {
		return fmt.Errorf("save items: %w", err)
	}

	s.items = newItems

	return nil
}

func (s *Storage) GetItems() []model.Item {
	return s.items
}

func (s *Storage) loadData() error {
	f, err := os.Open(s.config.Filename)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	var data storedData

	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	s.lastID = data.LastID
	s.items = make([]model.Item, 0, len(data.Items))

	for i := range data.Items {
		s.items = append(s.items, data.Items[i].toDomain())
	}

	return nil
}

func (s *Storage) saveData(data storedData) error {
	f, err := os.Create(s.config.Filename)
	if err != nil {
		return fmt.Errorf("os create: %w", err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	return nil
}

func (s *storedItem) fromDomain(item model.Item) {
	s.ID = item.ID
}

func (s *storedItem) toDomain() model.Item {
	return model.Item{ID: s.ID}
}
