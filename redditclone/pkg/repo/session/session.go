package session

import (
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/user"
)

type Session struct {
	ID   string
	User *user.User
}

func CreateNewSession(id string, user *user.User) *Session {
	return &Session{
		ID:   id,
		User: user,
	}
}
