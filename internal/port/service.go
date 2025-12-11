package port

type Service interface {
	AddItem() error
	Items() []string
}
