package port

type NumGenerator interface {
	GenerateNum() (int, error)
}
