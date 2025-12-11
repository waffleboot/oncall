package port

type Storage interface {
	AddItem(item string) error
	Items() []string
}
