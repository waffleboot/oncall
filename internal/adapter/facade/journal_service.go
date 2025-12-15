package facade

import (
	"errors"
	"fmt"
	"os"

	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type JournalService struct {
	storage port.Storage
}

func NewJournalService(storage port.Storage) *JournalService {
	return &JournalService{storage: storage}
}

func (s *JournalService) PrintJournal() (err error) {
	f, err := os.Create("journal.txt")
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	items, err := s.storage.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}

	m := make(map[model.ItemType][]model.Item)

	for _, item := range items {
		m[item.Type] = append(m[item.Type], item)
	}

	_, _ = fmt.Fprintln(f, "date")
	_, _ = fmt.Fprintln(f)

	for _, itemType := range []model.ItemType{
		model.ItemTypeInc,
		model.ItemTypeAdhoc,
		model.ItemTypeAsk,
		model.ItemTypeAlert,
	} {
		items := m[itemType]
		if len(items) == 0 {
			continue
		}

		_, _ = fmt.Fprintln(f, itemType.String())

		for _, item := range items {
			if links := item.PrintedLinks(); len(links) > 0 {
				for _, link := range links {
					_, _ = fmt.Fprintln(f)
					_, _ = fmt.Fprintln(f, link.Address)
				}
			}
		}
	}

	return nil
}

func (s *JournalService) CloseJournal() error {
	return s.storage.CloseJournal()
}
