package model

type User struct {
	Nick string
	Name string
}

func (u User) MenuItem() string {
	return u.Name
}

func (u User) ToPrint() string {
	return u.Nick
}
