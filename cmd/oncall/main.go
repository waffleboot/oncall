package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/waffleboot/oncall/internal/adapter/facade"
	teaAdapter "github.com/waffleboot/oncall/internal/adapter/tea"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	itemService := facade.NewItemService()

	teaModel := teaAdapter.NewTeaModel(itemService)

	if _, err := tea.NewProgram(teaModel).Run(); err != nil {
		return fmt.Errorf("tea run: %w", err)
	}

	return nil
}
