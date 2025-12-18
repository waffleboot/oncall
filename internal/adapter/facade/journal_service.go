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
	_, _ = fmt.Fprintf(w, "# %s\n", time.Now().Format(time.DateOnly))

	items, err := s.storage.GetItems()
	if err != nil {
		return fmt.Errorf("get items: %w", err)
	}

	m := make(map[model.ItemType][]model.Item)

	for _, item := range items {
		m[item.Type] = append(m[item.Type], item)
	}

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
			_, _ = fmt.Fprintln(w, "## Инциденты")
		case model.ItemTypeAdhoc:
			_, _ = fmt.Fprintln(w, "## ADHOC")
		case model.ItemTypeAsk:
			_, _ = fmt.Fprintln(w, "## Обращения")
		case model.ItemTypeAlert:
			_, _ = fmt.Fprintln(w, "## Алерты")
		}

		for i, item := range items {
			_, _ = fmt.Fprintf(w, "\n%d) %s\n", i+1, item.ToPrint())

			if len(item.Description) > 0 {
				_, _ = fmt.Fprintln(w)
				_, _ = fmt.Fprintln(w, item.Description)
			}

			newline := true
			if vms := item.PrintedVMs(); len(vms) > 0 {
				for _, vm := range vms {
					if newline {
						_, _ = fmt.Fprintln(w)
						newline = false
					}
					_, _ = fmt.Fprintln(w, vm.ToPrint())
					newline = vm.HasNode()
				}
			}

			if nodes := item.PrintedNodes(); len(nodes) > 0 {
				_, _ = fmt.Fprintln(w)
				for _, node := range nodes {
					_, _ = fmt.Fprintln(w, node.ToPrint())
				}
			}

			if links := item.PrintedLinks(); len(links) > 0 {
				_, _ = fmt.Fprintln(w)
				for _, link := range links {
					_, _ = fmt.Fprintln(w, link.ToPrint())
				}
			}

			if notes := item.PrintedNotes(); len(notes) > 0 {
				for _, note := range notes {
					_, _ = fmt.Fprintln(w)
					_, _ = fmt.Fprintln(w, note.ToPrint())
				}
			}
		}
	}

	return nil
}

func (s *JournalService) CloseJournal() error {
	return s.storage.CloseJournal()
}
