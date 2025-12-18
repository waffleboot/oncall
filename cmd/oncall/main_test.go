package main_test

import (
	_ "embed"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/waffleboot/oncall/internal/adapter/facade"
	storageAdapter "github.com/waffleboot/oncall/internal/adapter/storage"
	"github.com/waffleboot/oncall/internal/model"
	"go.uber.org/zap/zaptest"
)

//go:embed testdata/report.txt
var report string

func TestReport(t *testing.T) {
	journalFile := filepath.Join(t.TempDir(), "oncall.json")
	storageService := storageAdapter.NewStorage(storageAdapter.Config{
		Filename: journalFile,
	}, zaptest.NewLogger(t))
	itemService, err := facade.NewItemService(storageService, storageService, zaptest.NewLogger(t))
	assert.NoError(t, err)

	itemInc, err := itemService.CreateItem()
	assert.NoError(t, err)

	itemInc.Type = model.ItemTypeInc
	itemInc, err = itemService.UpdateItem(itemInc)
	assert.NoError(t, err)

	journalService := facade.NewJournalService(itemService)

	sb := new(strings.Builder)

	ts := time.Date(2025, 12, 18, 0, 0, 0, 0, time.UTC)

	err = journalService.PrintJournal(sb, ts)
	assert.NoError(t, err)

	assert.Equal(t, report, sb.String())
}
