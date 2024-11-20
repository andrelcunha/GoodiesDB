package store

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Store struct {
	Data    map[string]string
	Expires map[string]time.Time
	mu      sync.RWMutex
	aofChan chan string
}

// NewStore creates a new store
func NewStore(aofChan chan string) *Store {
	return &Store{
		Data:    make(map[string]string),
		Expires: make(map[string]time.Time),
		aofChan: aofChan,
	}
}

// Lock locks the store for writing
func (s *Store) Lock() {
	s.mu.Lock()
}

// Unlock unlocks the store
func (s *Store) Unlock() {
	s.mu.Unlock()
}

// RLock locks the store for reading
func (s *Store) RLock() {
	s.mu.RLock()
}

// RUnlock unlocks the store
func (s *Store) RUnlock() {
	s.mu.RUnlock()
}

// Set sets the value for a key
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[key] = value
	s.aofChan <- fmt.Sprintf("SET %s %s", key, value)
}

// Get gets the value for a key
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isExpired(key) {
		return "", false
	}
	value, ok := s.Data[key]
	return value, ok
}

func (s *Store) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Data, key)
	delete(s.Expires, key)
	s.aofChan <- fmt.Sprintf("DEL %s", key)
}

// Exists checks if a key exists
func (s *Store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isExpired(key) {
		return false
	}
	_, ok := s.Data[key]
	return ok
}

// SetNx sets the value for a key if the key does not exist
func (s *Store) SetNX(key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.Data[key]; exists {
		return false
	}
	s.Data[key] = value
	s.aofChan <- fmt.Sprintf("SET %s %s", key, value)
	return true
}

// Expire sets the expiration time for a key
func (s *Store) Expire(key string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.Data[key]; exists {
		s.Expires[key] = time.Now().Add(ttl)
		s.aofChan <- fmt.Sprintf("EXPIRE %s %d", key, int(ttl.Seconds()))
		return true
	}
	return false
}

// isExpired checks if a key has expired
func (s *Store) isExpired(key string) bool {
	if exp, exists := s.Expires[key]; exists {
		if time.Now().After(exp) {
			delete(s.Data, key)
			delete(s.Expires, key)
			s.aofChan <- fmt.Sprintf("DEL %s", key)
			return true
		}
	}
	return false
}

// Incr increments the value for a key
func (s *Store) Incr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.Data[key]
	if !ok {
		value = "0"
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value is not an integer or out of range")
	}

	intValue++
	s.Data[key] = strconv.Itoa(intValue)
	s.aofChan <- fmt.Sprintf("INCR %s", key)
	return intValue, nil
}

// Decr decrements the value for a key
func (s *Store) Decr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.Data[key]
	if !ok {
		value = "0"
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value is not an integer or out of range")
	}

	intValue--
	s.Data[key] = strconv.Itoa(intValue)
	s.aofChan <- fmt.Sprintf("DECR %s", key)
	return intValue, nil
}

// TTL Retrieve the remaining time to live for a key
func (s *Store) TTL(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Data[key]; !ok {
		return -2, nil
	}

	if _, ok := s.Expires[key]; !ok {
		return -1, nil
	}

	ttl := s.Expires[key].Sub(time.Now())
	return int(ttl.Seconds()), nil
}
