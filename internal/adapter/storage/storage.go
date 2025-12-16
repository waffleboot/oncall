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
		items   []storedItem
	}
	storedData struct {
		LastNum int          `json:"last_num,omitempty"`
		Items   []storedItem `json:"items,omitempty"`
	}
	storedItem struct {
		ID          uuid.UUID    `json:"id"`
		Num         int          `json:"num"`
		SleepAt     time.Time    `json:"sleepAt,omitempty"`
		ClosedAt    time.Time    `json:"closedAt,omitempty"`
		DeletedAt   time.Time    `json:"deletedAt,omitempty"`
		Links       []storedLink `json:"links,omitempty"`
		VMs         []vm         `json:"vms,omitempty"`
		Type        string       `json:"type,omitempty"`
		Title       string       `json:"title,omitempty"`
		Description string       `json:"description,omitempty"`
	}
	storedLink struct {
		ID          int       `json:"id"`
		Link        string    `json:"link,omitempty"`
		Public      bool      `json:"public"`
		DeletedAt   time.Time `json:"deleted_at,omitempty"`
		Description string    `json:"description,omitempty"`
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

func (s *Storage) DeleteItem(itemID uuid.UUID) error {
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

func (s *storedLink) fromDomain(link model.Link) {
	s.ID = link.ID
	s.Link = link.Address
	s.Public = link.Public
	s.DeletedAt = link.DeletedAt
	s.Description = link.Description
}

func (s *storedLink) toDomain() model.Link {
	return model.Link{
		ID:          s.ID,
		Address:     s.Link,
		Public:      s.Public,
		DeletedAt:   s.DeletedAt,
		Description: s.Description,
	}
}
