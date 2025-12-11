package port

type Service interface {
	AddItem(item string) error
	DeleteItem(item string) error
	Items() []string
}
