package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
)

type (
	Config struct {
		Filename string
	}
	Storage struct {
		config Config
		items  []string
	}
	storedItem struct {
		Name string `json:"name"`
	}
)

func NewStorage(config Config) (_ *Storage, err error) {
	s := &Storage{config: config}
	if err := s.loadItems(); err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}
	return s, nil
}

func (s *Storage) AddItem(item string) (err error) {
	items := make([]storedItem, len(s.items)+1)
	for i := range s.items {
		items[i].fromDomain(s.items[i])
	}
	items[len(items)-1].fromDomain(item)

	if err := s.saveItems(items); err != nil {
		return fmt.Errorf("save items: %w", err)
	}

	s.items = append(s.items, item)

	return nil
}

func (s *Storage) Items() []string {
	return s.items
}

func (s *Storage) loadItems() error {
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

	var items []storedItem

	if err := json.NewDecoder(f).Decode(&items); err != nil {
		return fmt.Errorf("json decode: %w", err)
	}

	s.items = make([]string, 0, len(items))

	for i := range items {
		s.items = append(s.items, items[i].toDomain())
	}

	return nil
}

func (s *Storage) saveItems(items []storedItem) error {
	f, err := os.Create(s.config.Filename)
	if err != nil {
		return fmt.Errorf("os create: %w", err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")

	if err := enc.Encode(items); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	return nil
}

func (s *storedItem) fromDomain(item string) {
	s.Name = item
}

func (s *storedItem) toDomain() string {
	return s.Name
}
