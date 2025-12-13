package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/waffleboot/oncall/internal/model"
)

type (
	Config struct {
		Filename string
	}
	Storage struct {
		config Config
		lastID int
		items  []storedItem
	}
	storedData struct {
		LastID int          `json:"last_id,omitempty"`
		Items  []storedItem `json:"items,omitempty"`
	}
	storedItem struct {
		ID        int       `json:"id"`
		SleepAt   time.Time `json:"sleepAt,omitempty"`
		ClosedAt  time.Time `json:"closedAt,omitempty"`
		DeletedAt time.Time `json:"deletedAt,omitempty"`
		Type      string    `json:"type,omitempty"`
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

func (s *Storage) GetItem(itemID int) (model.Item, error) {
	for i := range s.items {
		if s.items[i].ID == itemID {
			return s.items[i].toDomain(), nil
		}
	}

	return model.Item{}, errors.New("not found")
}

func (s *Storage) UpdateItem(item model.Item) error {
	var found bool

	for i := range s.items {
		if s.items[i].ID == item.ID {
			s.items[i].fromDomain(item)
			found = true
		}
	}

	if !found {
		var st storedItem
		st.fromDomain(item)
		s.items = append(s.items, st)
	}

	if err := s.saveData(); err != nil {
		return fmt.Errorf("save data: %w", err)
	}

	return nil
}

func (s *Storage) DeleteItem(itemID int) error {
	var found bool

	for i := range s.items {
		if s.items[i].ID == itemID {
			s.items[i].DeletedAt = time.Now()
			found = true
		}
	}

	if !found {
		return fmt.Errorf("item not found")
	}

	if err := s.saveData(); err != nil {
		return fmt.Errorf("save data: %w", err)
	}

	return nil
}

func (s *Storage) GetItems() ([]model.Item, error) {
	items := make([]model.Item, 0, len(s.items))

	for i := range s.items {
		if s.items[i].DeletedAt.IsZero() {
			items = append(items, s.items[i].toDomain())
		}
	}

	return items, nil
}

func (s *Storage) CloseJournal() error {
	s.lastID = 0
	s.items = nil
	return s.saveData()
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
	s.items = data.Items

	return nil
}

func (s *Storage) saveData() error {
	f, err := os.Create(s.config.Filename)
	if err != nil {
		return fmt.Errorf("os create: %w", err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")

	if err := enc.Encode(storedData{
		LastID: s.lastID,
		Items:  s.items,
	}); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	return nil
}

func (s *storedItem) fromDomain(item model.Item) {
	s.ID = item.ID
	s.SleepAt = item.SleepAt.UTC()
	s.ClosedAt = item.ClosedAt.UTC()
	s.Type = string(item.Type)
}

func (s *storedItem) toDomain() model.Item {
	return model.Item{
		ID:       s.ID,
		SleepAt:  s.SleepAt,
		ClosedAt: s.ClosedAt,
		Type:     model.ItemType(s.Type),
	}
}
