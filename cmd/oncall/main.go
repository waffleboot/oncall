package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/waffleboot/oncall/internal/adapter/facade"
	"github.com/waffleboot/oncall/internal/adapter/storage"
	"github.com/waffleboot/oncall/internal/adapter/tea"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() (err error) {
	s, err := storage.NewStorage(storage.Config{Filename: "oncall.json"})
	if err != nil {
		return fmt.Errorf("new storage: %w", err)
	}

	log, err := getLogger()
	if err != nil {
		return fmt.Errorf("get logger: %w", err)
	}
	defer func() {
		err = errors.Join(err, log.Sync())
	}()

	itemService := facade.NewItemService(s, s)

	journalService := facade.NewJournalService(s)

	p := tea.NewController(tea.WithService(itemService, journalService, log))
	if err := p.Run(); err != nil {
		return fmt.Errorf("tea run: %w", err)
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
