package service

import (
	"errors"
	"sync"
)

type UserStore interface {
	Save(user *User) error
	Find(username string) (*User, error)
}

type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

func (s *InMemoryUserStore) Save(user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.users[user.Username] != nil {
		return errors.New("")
	}

	s.users[user.Username] = user.Clone()
	return nil
}

func (s *InMemoryUserStore) Find(username string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user := s.users[username]
	if user == nil {
		return nil, nil
	}

	return user.Clone(), nil
}
