package session

import (
	uuid "github.com/satori/go.uuid"
	"sync"
)

type MemorySession struct {
	data   map[string]interface{}
	id     string
	rwLock sync.RWMutex
}

func NewMemorySession(id string) *MemorySession {
	s := &MemorySession{
		id:   id,
		data: make(map[string]interface{}, 8),
	}
	return s
}

func (m *MemorySession) Set(key string, value interface{}) (err error) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	m.data[key] = value
	return
}

func (m *MemorySession) Get(key string) (value interface{}, err error) {
	m.rwLock.RLock()
	defer m.rwLock.RUnlock()
	value, ok := m.data[key]
	if !ok {
		err = ErrKeyNotExistInSession
	}
	return
}

func (m *MemorySession) Del(key string) (err error) {
	m.rwLock.Lock()
	defer m.rwLock.Unlock()
	delete(m.data, key)
	return
}

func (m *MemorySession) Save() (err error) {
	return
}


type MemorySessionMgr struct {
	sessionMap map[string]Session
	rwLock     sync.RWMutex
}

func NewMemorySessionMgr() *MemorySessionMgr {
	mgr := &MemorySessionMgr{
		sessionMap: make(map[string]Session, 1024),
	}
	return mgr
}

func (s *MemorySessionMgr) Init(addr string, options ...string) (err error) {
	return
}

func (s *MemorySessionMgr) Get(sessionId string) (session Session, err error) {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	session, ok := s.sessionMap[sessionId]
	if !ok {
		err = ErrSessionNotExist
		return
	}

	return
}

func (s *MemorySessionMgr) CreateSession() (session Session, err error) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	sessionId := uuid.NewV4().String()
	session = NewMemorySession(sessionId)
	s.sessionMap[sessionId] = session
	return
}

