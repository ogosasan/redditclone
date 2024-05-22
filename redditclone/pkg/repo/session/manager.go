package session

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/user"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type SessionsManager struct {
	data              map[string]*Session
	mu                *sync.RWMutex
	previousSessionID int
	secretToken       []byte
}

func NewSessionManager() *SessionsManager {
	return &SessionsManager{
		data:              make(map[string]*Session),
		mu:                &sync.RWMutex{},
		previousSessionID: 0,
		secretToken:       []byte(os.Getenv("SECRET")),
	}
}

func (sm *SessionsManager) CreateToken(user *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(sm.secretToken)
	if err != nil {
		log.Printf("Token is not signed to string %s", err)
		return "", err
	}
	return tokenString, nil
}

func (sm *SessionsManager) CreateSession(user *user.User) (string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.previousSessionID++
	token, err := sm.CreateToken(user)
	if err != nil {
		log.Printf("Session is not created %s", err)
		return "", err
	}
	sm.data[token] = CreateNewSession(strconv.Itoa(sm.previousSessionID), user)
	return token, nil
}

func (sm *SessionsManager) GetSession(inToken string) (*Session, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return sm.secretToken, nil
	}
	token, err := jwt.Parse(inToken, hashSecretGetter)
	if err != nil || !token.Valid {
		log.Printf("Can not parse this token %s", err)
		return nil, err
	}
	sm.mu.RLock()
	curSession, ok := sm.data[inToken]
	sm.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no such session")
	}
	return curSession, nil
}
