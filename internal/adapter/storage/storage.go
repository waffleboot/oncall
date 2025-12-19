package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/waffleboot/oncall/internal/model"
	"go.uber.org/zap"
)

type (
	Config struct {
		JournalName string
		Users       []model.User
	}
	Storage struct {
		lastNum int
		config  Config
		log     *zap.Logger
	}
)

func NewStorage(config Config, log *zap.Logger) *Storage {
	return &Storage{config: config, log: log}
}

func (s *Storage) GetJournal() (model.Journal, error) {
	j, err := s.loadJournal()
	if err != nil {
		return model.Journal{}, err
	}

	var user *model.User
	for i := range s.config.Users {
		if s.config.Users[i].Nick == j.Next {
			user = &s.config.Users[i]
			break
		}
	}

	return j.toDomain(user), nil
}

func (s *Storage) SaveJournal(j model.Journal) error {
	s.log.Debug("save journal", zap.Int("items_count", len(j.Items)))
	var st journal
	st.fromDomain(j)
	return s.saveJournal(st)
}

func (s *Storage) CloseJournal(j model.Journal) error {
	ts := time.Now().Format("2006-01-02-15-04-05")
	to := fmt.Sprintf("%s.%s", s.config.JournalName, ts)

	if err := os.Rename(s.config.JournalName, to); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("rename: %w", err)
	}

	return nil
}

func (s *Storage) GenerateNum() (int, error) {
	j, err := s.loadJournal()
	if err != nil {
		return 0, err
	}
	s.lastNum++
	if err := s.saveJournal(j); err != nil {
		return 0, err
	}
	s.log.Debug("generate num", zap.Int("num", s.lastNum))
	return s.lastNum, nil
}

func (s *Storage) loadJournal() (journal, error) {
	f, err := os.Open(s.journalJson())
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return journal{}, nil
		}
		return journal{}, fmt.Errorf("open file: %w", err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	var j journal

	if err := json.NewDecoder(f).Decode(&j); err != nil {
		return journal{}, fmt.Errorf("json decode: %w", err)
	}

	s.log.Debug("journal loaded", zap.Int("items_count", len(j.Items)), zap.Int("last_num", j.LastNum))

	return s.afterLoadJournal(j), nil
}

func (s *Storage) saveJournal(j journal) error {
	if err := os.MkdirAll(s.config.JournalName, 0o775); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	if err := os.MkdirAll(s.filesDir(), 0o775); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	f, err := os.Create(s.journalJson())
	if err != nil {
		return fmt.Errorf("os create: %w", err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	enc := json.NewEncoder(f)
	enc.SetIndent("", " ")

	if err := enc.Encode(s.beforeSaveJournal(j)); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	if err := f.Sync(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	s.log.Debug("journal saved", zap.Int("items_count", len(j.Items)), zap.Int("last_num", j.LastNum))

	return nil
}

func (s *Storage) journalJson() string {
	return filepath.Join(s.config.JournalName, "journal.json")
}

func (s *Storage) beforeSaveJournal(j journal) journal {
	j.LastNum = s.lastNum
	return j
}

func (s *Storage) afterLoadJournal(j journal) journal {
	s.lastNum = j.LastNum
	return j
}

func (s *Storage) filesDir() string {
	return filepath.Join(s.config.JournalName, "files")
}
