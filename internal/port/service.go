package port

import "github.com/waffleboot/oncall/internal/model"

type ItemService interface {
	CreateItem() model.Item

	AddItem(item model.Item) error
	GetItem(itemID int) (model.Item, error)

	SetItemType(item model.Item, itemType model.ItemType) error

	SleepItem(item model.Item) error
	AwakeItem(item model.Item) error
	CloseItem(item model.Item) error
	DeleteItem(item model.Item) error

	GetItems() ([]model.Item, error)
}

type JournalService interface {
	CloseJournal() error
}
