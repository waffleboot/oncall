package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

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

	storage := storageAdapter.NewStorage(storageAdapter.Config{Filename: "oncall.json"}, log.Named("storage"))

	itemService, err := facade.NewItemService(storage, storage, log.Named("item_service"))
	if err != nil {
		return fmt.Errorf("new item service: %w", err)
	}

	journalService := facade.NewJournalService(itemService)

	teaModel := teaAdapter.NewTeaModel(itemService, itemService, log)

	if _, err := tea.NewProgram(teaModel).Run(); err != nil {
		return fmt.Errorf("tea run: %w", err)
	}

	if err := teaModel.Err(); err != nil {
		return fmt.Errorf("tea err: %w", err)
	}

	if teaModel.PrintJournal() {
		fmt.Println("----------")

		printJournal := func(w io.Writer) error {
			if err := journalService.PrintJournal(w, time.Now()); err != nil {
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
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)

	log, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return log, nil
}
