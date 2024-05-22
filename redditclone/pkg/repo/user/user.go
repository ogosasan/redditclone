package user

import (
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/utils"
)

type User struct {
	Login    string `json:"username"`
	ID       string `json:"id"`
	Password string `json:"-"`
}

type UserRepo interface {
	Login(login, password string) (*User, error)
	Register(login, password string) (*User, error)
}

func CreateUser(login, password, id string) (*User, error) {
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, ErrInHashingPass
	}
	return &User{
		Login:    login,
		Password: hashPassword,
		ID:       id,
	}, nil
}
