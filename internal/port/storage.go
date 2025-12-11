package port

type Storage interface {
	AddItem(item string) error
	DeleteItem(item string) error
	Items() []string
}
