package store

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Store struct {
	Data    []map[string]interface{}
	Expires []map[string]time.Time
	mu      sync.RWMutex
	aofChan chan string
}

// NewStore creates a new store
func NewStore(aofChan chan string) *Store {
	data := make([]map[string]interface{}, 16)
	expires := make([]map[string]time.Time, 16)
	for i := range data {
		data[i] = make(map[string]interface{})
		expires[i] = make(map[string]time.Time)
	}
	return &Store{
		Data:    data,
		Expires: expires,
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
func (s *Store) Set(dbIndex int, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Data[dbIndex][key] = value
	s.aofChan <- fmt.Sprintf("SET %d %s %s", dbIndex, key, value)
}

// Get gets the value for a key
func (s *Store) Get(dbIndex int, key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isExpired(dbIndex, key) {
		return "", false
	}
	value, ok := s.Data[dbIndex][key].(string)
	return value, ok
}

func (s *Store) Del(dbIndex int, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Data[dbIndex], key)
	delete(s.Expires[dbIndex], key)
	s.aofChan <- fmt.Sprintf("DEL %d %s", dbIndex, key)
}

// Exists checks if a key exists
func (s *Store) Exists(dbIndex int, key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.isExpired(dbIndex, key) {
		return false
	}
	_, ok := s.Data[dbIndex][key]
	return ok
}

// SetNx sets the value for a key if the key does not exist
func (s *Store) SetNX(dbIndex int, key, value string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.Data[dbIndex][key]; exists {
		return false
	}
	s.Data[dbIndex][key] = value
	s.aofChan <- fmt.Sprintf("SET %d %s %s", dbIndex, key, value)
	return true
}

// Expire sets the expiration time for a key
func (s *Store) Expire(dbIndex int, key string, ttl time.Duration) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.Data[dbIndex][key]; exists {
		s.Expires[dbIndex][key] = time.Now().Add(ttl)
		s.aofChan <- fmt.Sprintf("EXPIRE %d %s %d", dbIndex, key, int(ttl.Seconds()))
		return true
	}
	return false
}

// isExpired checks if a key has expired
func (s *Store) isExpired(dbIndex int, key string) bool {
	if exp, exists := s.Expires[dbIndex][key]; exists {
		if time.Now().After(exp) {
			delete(s.Data[dbIndex], key)
			delete(s.Expires[dbIndex], key)
			s.aofChan <- fmt.Sprintf("DEL %d %s", dbIndex, key)
			return true
		}
	}
	return false
}

// Incr increments the value for a key
func (s *Store) Incr(dbIndex int, key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.Data[dbIndex][key].(string)
	if !ok {
		value = "0"
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value is not an integer or out of range")
	}

	intValue++
	s.Data[dbIndex][key] = strconv.Itoa(intValue)
	s.aofChan <- fmt.Sprintf("INCR %d %s", dbIndex, key)
	return intValue, nil
}

// Decr decrements the value for a key
func (s *Store) Decr(dbIndex int, key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.Data[dbIndex][key].(string)
	if !ok {
		value = "0"
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value is not an integer or out of range")
	}

	intValue--
	s.Data[dbIndex][key] = strconv.Itoa(intValue)
	s.aofChan <- fmt.Sprintf("DECR %d %s", dbIndex, key)
	return intValue, nil
}

// TTL Retrieve the remaining time to live for a key
func (s *Store) TTL(dbIndex int, key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Data[dbIndex][key]; !ok {
		return -2, nil
	}

	if _, ok := s.Expires[dbIndex][key]; !ok {
		return -1, nil
	}

	ttl := s.Expires[dbIndex][key].Sub(time.Now())
	return int(ttl.Seconds()), nil
}
