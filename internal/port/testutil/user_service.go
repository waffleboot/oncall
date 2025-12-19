package testutil

import "github.com/waffleboot/oncall/internal/model"

type UserService string

func (s UserService) GetUser() *model.User {
	return &model.User{Nick: string(s)}
}

func (s UserService) SetUser(model.User) error {
	return nil
}
