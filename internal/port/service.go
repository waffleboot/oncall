package port

import (
	"github.com/waffleboot/oncall/internal/model"
)

type ItemService interface {
	GetItem(id int) (model.Item, error)
	GetItems() []model.Item
	CreateItem() (model.Item, error)
	UpdateItem(model.Item) (model.Item, error)
	DeleteItem(model.Item) (model.Item, error)
	AwakeItem(model.Item) (model.Item, error)
	SleepItem(model.Item) (model.Item, error)
	CloseItem(model.Item) (model.Item, error)
}

type UserService interface {
	GetUser() *model.User
	SetUser(model.User) error
}

type JournalService interface {
	CloseJournal() error
}
