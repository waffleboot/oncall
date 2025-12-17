package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/adapter/facade"
	storageAdapter "github.com/waffleboot/oncall/internal/adapter/storage"
	teaAdapter "github.com/waffleboot/oncall/internal/adapter/tea"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() (err error) {
	log, err := getLogger()
	if err != nil {
		return fmt.Errorf("get logger: %w", err)
	}
	defer func() {
		err = errors.Join(err, log.Sync())
	}()

	storage, err := storageAdapter.NewStorage(storageAdapter.Config{Filename: "oncall.json"})
	if err != nil {
		return fmt.Errorf("new storage: %w", err)
	}

	itemService := facade.NewItemService(storage, storage)

	journalService := facade.NewJournalService(storage)

	teaModel := teaAdapter.NewTeaModel(itemService, journalService, log)

	if _, err := tea.NewProgram(teaModel).Run(); err != nil {
		return fmt.Errorf("tea run: %w", err)
	}

	if err := teaModel.Err(); err != nil {
		return fmt.Errorf("tea err: %w", err)
	}

	if teaModel.PrintJournal() {
		fmt.Println("----------")

		printJournal := func(w io.Writer) error {
			if err := journalService.PrintJournal(w); err != nil {
				return fmt.Errorf("print journal: %w", err)
			}
			return nil
		}

		f, err := os.Create("journal.txt")
		if err != nil {
			return printJournal(os.Stdout)
		}
		defer func() {
			err = errors.Join(f.Close())
		}()

		return printJournal(io.MultiWriter(os.Stdout, f))
	}

	return nil
}

func getLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"debug.log"}
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	log, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return log, nil
}
