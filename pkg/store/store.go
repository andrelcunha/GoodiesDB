package store

import "sync"

type Store struct {
	data map[string]string
	mu   sync.RWMutex // Allow multiple read, single write
}

// NewStore creates a new store
func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
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
	value, ok := s.data[key]
	return value, ok
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

// Exists checks if a key exists
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
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
