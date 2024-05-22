package user

import (
	"errors"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/utils"
	"strconv"
	"sync"
)

var (
	ErrNoUser        = errors.New("no user found")
	ErrBadPass       = errors.New("invalid password")
	ErrUserExist     = errors.New("such a user already exists")
	ErrInHashingPass = errors.New("error in hashing password")
)

type UserMemoryRepository struct {
	previouslyUserID int
	data             map[string]*User
	mu               *sync.RWMutex
}

func NewMemoryRepo() *UserMemoryRepository {
	return &UserMemoryRepository{
		data: map[string]*User{},
		mu:   &sync.RWMutex{},
	}
}

func (repository *UserMemoryRepository) Login(login, password string) (*User, error) {
	u, ok := repository.data[login]
	if !ok {
		return nil, ErrNoUser
	}
	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, ErrInHashingPass
	}
	if u.Password != hashPassword {
		return nil, ErrBadPass
	}
	return u, nil
}

func (repository *UserMemoryRepository) Register(login, password string) (*User, error) {
	repository.mu.Lock()
	defer repository.mu.Unlock()
	if _, ok := repository.data[login]; ok {
		return nil, ErrUserExist
	}
	repository.previouslyUserID++
	u, err := CreateUser(login, password, strconv.Itoa(repository.previouslyUserID))
	if err != nil {
		return nil, err
	}
	repository.data[login] = u
	return u, nil
}
