package facade

import (
	"fmt"
	"io"
	"time"

	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port"
)

type JournalService struct {
	itemService port.ItemService
}

func NewJournalService(itemService port.ItemService) *JournalService {
	return &JournalService{itemService: itemService}
}

func (s *JournalService) PrintJournal(w io.Writer, at time.Time) error {
	write := func(format string, args ...any) {
		_, _ = fmt.Fprintf(w, format+"\n", args...)
	}

	write("# %s", time.Now().Format(time.DateOnly))

	items := s.itemService.GetItems()

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

		write("")

		switch itemType {
		case model.ItemTypeInc:
			write("## Инциденты")
		case model.ItemTypeAdhoc:
			write("## ADHOC")
		case model.ItemTypeAsk:
			write("## Обращения")
		case model.ItemTypeAlert:
			write("## Алерты")
		}

		for i, item := range items {
			write("\n%d) %s\n", i+1, item.ToPrint())

			if len(item.Description) > 0 {
				write("")
				write(item.Description)
			}

			newline := true
			if vms := item.PrintedVMs(); len(vms) > 0 {
				for _, vm := range vms {
					if newline {
						write("")
						newline = false
					}
					write(vm.ToPrint())
					newline = vm.HasNode()
				}
			}

			if nodes := item.PrintedNodes(); len(nodes) > 0 {
				write("")
				for _, node := range nodes {
					write(node.ToPrint())
				}
			}

			if links := item.PrintedLinks(); len(links) > 0 {
				write("")
				for _, link := range links {
					write(link.ToPrint())
				}
			}

			if notes := item.PrintedNotes(); len(notes) > 0 {
				for _, note := range notes {
					write("")
					write(note.ToPrint())
				}
			}
		}
	}

	return nil
}
