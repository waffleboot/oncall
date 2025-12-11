package main

import (
	"fmt"
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

func run() error {
	s, err := storage.NewStorage(storage.Config{Filename: "oncall.json"})
	if err != nil {
		return fmt.Errorf("new storage: %w", err)
	}

	p := tea.NewController(tea.WithService(facade.NewService(s, s)))
	if err := p.Run(); err != nil {
		return fmt.Errorf("tea run: %w", err)
	}

	return nil
}
