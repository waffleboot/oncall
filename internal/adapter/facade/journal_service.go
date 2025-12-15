package facade

import (
	"fmt"
	"io"
	"time"

	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type JournalService struct {
	storage port.Storage
}

func NewJournalService(storage port.Storage) *JournalService {
	return &JournalService{storage: storage}
}

func (s *JournalService) PrintJournal(w io.Writer) (err error) {
	items, err := s.storage.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}

	m := make(map[model.ItemType][]model.Item)

	for _, item := range items {
		m[item.Type] = append(m[item.Type], item)
	}

	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, time.Now().Format(time.DateOnly))

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

		_, _ = fmt.Fprintln(w)

		switch itemType {
		case model.ItemTypeInc:
			_, _ = fmt.Fprintln(w, "# Инциденты")
		case model.ItemTypeAdhoc:
			_, _ = fmt.Fprintln(w, "# ADHOC")
		case model.ItemTypeAsk:
			_, _ = fmt.Fprintln(w, "# Обращения")
		case model.ItemTypeAlert:
			_, _ = fmt.Fprintln(w, "# Алерты")
		}

		for i, item := range items {

			_, _ = fmt.Fprintf(w, "\n%d) %s #%d", i+1, item.TitleForView(), item.Num)

			switch {
			case item.IsActive():
				_, _ = fmt.Fprintf(w, " (in progress)\n")
			case item.IsSleep():
				_, _ = fmt.Fprintf(w, " (in progress, in waiting)\n")
			}

			if len(item.Description) > 0 {
				_, _ = fmt.Fprintln(w)
				_, _ = fmt.Fprintln(w, item.Description)
			}

			if links := item.PrintedLinks(); len(links) > 0 {
				for _, link := range links {
					_, _ = fmt.Fprintln(w)
					_, _ = fmt.Fprintln(w, link.Address)
				}
			}
		}
	}

	return nil
}

func (s *JournalService) CloseJournal() error {
	return s.storage.CloseJournal()
}
