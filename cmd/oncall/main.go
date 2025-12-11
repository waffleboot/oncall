package main

import (
	"fmt"
	"os"

	"github.com/waffleboot/oncall/internal/adapter/facade"
	"github.com/waffleboot/oncall/internal/adapter/tea"
)

func main() {
	p := tea.NewController(tea.WithOnCallService(facade.NewOnCallService()))
	if err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
