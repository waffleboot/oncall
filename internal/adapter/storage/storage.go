package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/waffleboot/oncall/internal/model"
)

type (
	Config struct {
		Filename string
	}
	Storage struct {
		config  Config
		lastNum int
		items   []item
	}
)

func NewStorage(config Config) (*Storage, error) {
	s := &Storage{config: config}
	if err := s.loadData(); err != nil {
		return nil, fmt.Errorf("load items: %w", err)
	}
	return s, nil
}

func (s *Storage) GenerateNum() int {
	s.lastNum++
	return s.lastNum
}

func (s *Storage) GetItem(itemID uuid.UUID) (model.Item, error) {
	for _, item := range s.items {
		if item.ID == itemID {
			return item.toDomain(), nil
		}
	}

	return model.Item{}, errors.New("not found")
}

func (s *Storage) UpdateItem(it model.Item) error {
	var found bool

	for i := range s.items {
		if s.items[i].ID == it.ID {
			s.items[i].fromDomain(it)
			found = true
		}
	}

	if !found {
		var st item
		st.fromDomain(it)
		s.items = append(s.items, st)
	}

	if err := s.saveData(); err != nil {
		return fmt.Errorf("save data: %w", err)
	}

	return nil
}

func (s *Storage) DeleteItem(itemID uuid.UUID) error {
	var found bool

	for i := range s.items {
		if s.items[i].ID == itemID {
			s.items[i].DeletedAt = time.Now()
			found = true
		}
	}

	if !found {
		return nil
	}

	if err := s.saveData(); err != nil {
		return fmt.Errorf("save data: %w", err)
	}

	return nil
}

func (s *Storage) GetItems() ([]model.Item, error) {
	items := make([]model.Item, 0, len(s.items))

	for i := range s.items {
		if s.items[i].NotDeleted() {
			items = append(items, s.items[i].toDomain())
		}
	}

	return items, nil
}

func (s *Storage) CloseJournal() error {
	ts := time.Now().Format("2006-01-02-15-04-05")
	to := fmt.Sprintf("%s.%s", s.config.Filename, ts)

	if err := os.Rename(s.config.Filename, to); err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	s.lastNum = 0
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

	s.lastNum = data.LastNum
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
		LastNum: s.lastNum,
		Items:   s.items,
	}); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	return nil
}

func from(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	t = t.UTC()
	return &t
}

func to[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}
