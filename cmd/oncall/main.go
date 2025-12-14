package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/waffleboot/oncall/internal/adapter/facade"
	"github.com/waffleboot/oncall/internal/adapter/tea"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() (err error) {
	itemService := facade.NewItemService()

	teaModel := tea.NewTeaModel(tea.TeaModelConfig{
		ItemService: itemService,
	})

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
