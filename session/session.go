package session

import "errors"

var (
	ErrSessionNotExist      = errors.New("session not exists")
	ErrKeyNotExistInSession = errors.New("key not exists in session")
)

type Session interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Del(key string) error
	Save() error
}

type Manager interface {
	Init(addr string, options ...string) (err error)
	CreateSession() (session Session, err error)
	Get(sessionId string) (session Session, err error)
}