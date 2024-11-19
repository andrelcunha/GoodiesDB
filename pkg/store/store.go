package store

import (
	"sync"
	"time"
)

type Store struct {
	data    map[string]string
	expires map[string]time.Time
	mu      sync.RWMutex // Allow multiple read, single write
}

// NewStore creates a new store
func NewStore() *Store {
	return &Store{
		data:    make(map[string]string),
		expires: make(map[string]time.Time),
	}
}

// Set sets the value for a key
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get gets the value for a key
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isExpired(key) {
		return "", false
	}
	value, ok := s.data[key]
	return value, ok
}

func (s *Store) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Exists checks if a key exists
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isExpired(key) {
		return false
	}
	_, ok := s.data[key]
	return ok
}

// SetNx sets the value for a key if the key does not exist
func (s *Store) SetNX(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; exists {
		return false
	}
	s.data[key] = value
	return true
}

// Expire sets the expiration time for a key
func (s *Store) Expire(key string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; exists {
		s.expires[key] = time.Now().Add(ttl)
		return true
	}
	return false
}

// isExpired checks if a key has expired
func (s *Store) isExpired(key string) bool {
	if exp, exists := s.expires[key]; exists {
		if time.Now().After(exp) {
			delete(s.data, key)
			delete(s.expires, key)
			return true
		}
	}
	return false
}
