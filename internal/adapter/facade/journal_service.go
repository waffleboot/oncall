package facade

import "github.com/waffleboot/oncall/internal/port"

type JournalService struct {
	storage port.Storage
}

func NewJournalService(storage port.Storage) *JournalService {
	return &JournalService{storage: storage}
}

func (s *JournalService) CloseJournal() error {
	return s.storage.CloseJournal()
}
