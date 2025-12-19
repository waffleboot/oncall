package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/waffleboot/oncall/internal/model"
)

type (
	Users struct {
		users []user
	}
	user struct {
		Nick string `json:"nick"`
		Name string `json:"name"`
	}
)

func GetUsers(filename string) (_ []model.User, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer func() {
		err = errors.Join(err, f.Close())
	}()

	var users []user

	if err := json.NewDecoder(f).Decode(&users); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	modelUsers := make([]model.User, 0, len(users))
	for _, user := range users {
		modelUsers = append(modelUsers, user.toDomain())
	}

	return modelUsers, nil
}

func (u user) toDomain() model.User {
	return model.User{
		Nick: u.Nick,
		Name: u.Name,
	}
}
