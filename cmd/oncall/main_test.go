package main_test

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/waffleboot/oncall/internal/adapter/facade"
	storageAdapter "github.com/waffleboot/oncall/internal/adapter/storage"
	"github.com/waffleboot/oncall/internal/model"
	"github.com/waffleboot/oncall/internal/port/testutil"
	"go.uber.org/zap/zaptest"
)

//go:embed testdata/report.txt
var report string

func TestReport(t *testing.T) {
	journalName := filepath.Join(t.TempDir(), "oncall")
	storageService := storageAdapter.NewStorage(storageAdapter.Config{
		JournalName: journalName,
	}, zaptest.NewLogger(t))
	itemService, err := facade.NewItemService(storageService, storageService, zaptest.NewLogger(t))
	assert.NoError(t, err)

	itemInc, err := itemService.CreateItem()
	itemInc.Title = "какой-то инцидент"
	itemInc.Description = "здесь дается описание на несколько строк\nс переводом строки\n\nдва раза"

	itemInc.Type = model.ItemTypeInc
	assert.NoError(t, err)

	itemInc.VMs = append(itemInc.VMs, model.VM{
		Name:        "vm-1",
		Node:        "node-1",
		Description: "можно дать описание вм\nна несколько строк",
	})

	itemInc.VMs = append(itemInc.VMs, model.VM{
		Name: "vm-2",
	})

	itemInc.Nodes = append(itemInc.Nodes, model.Node{
		Name:        "node-3",
		Description: "можно дать описание узла\nна несколько строк",
	})

	itemInc.Links = append(itemInc.Links, model.Link{
		Public:  true,
		Address: "http://jira.com",
	})

	itemInc.Links = append(itemInc.Links, model.Link{
		Public:      false,
		Address:     "http://confluence.com",
		Description: "описание работы сервиса",
	})

	itemInc.Notes = append(itemInc.Notes, model.Note{
		Text:   "выдали рекомендацию 1",
		Public: true,
	})

	itemInc.Notes = append(itemInc.Notes, model.Note{
		Text:   "выдали рекомендацию 2",
		Public: true,
	})

	itemInc.Notes = append(itemInc.Notes, model.Note{
		Text:   "этого описания не должно быть",
		Public: false,
	})

	itemInc, err = itemService.UpdateItem(itemInc)
	assert.NoError(t, err)

	journalService := facade.NewJournalService(itemService, testutil.UserService("nick"))

	sb := new(strings.Builder)

	ts := time.Date(2025, 12, 18, 0, 0, 0, 0, time.UTC)

	err = journalService.PrintJournal(sb, ts)
	assert.NoError(t, err)

	assert.Equal(t, report, sb.String())

	if t.Failed() {
		os.WriteFile("journal.txt", []byte(sb.String()), 0o644)
	}
}
