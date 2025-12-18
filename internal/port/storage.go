package port

import (
	"github.com/waffleboot/oncall/internal/model"
)

type Storage interface {
	GetJournal() (model.Journal, error)
	SaveJournal(model.Journal) error
	CloseJournal(model.Journal) error
}
