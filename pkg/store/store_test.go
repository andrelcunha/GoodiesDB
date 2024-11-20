package store

import (
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	aofChan := make(chan string, 100)

	s := NewStore(aofChan)
	s.Set(0, "Key1", "Value1")
	value, ok := s.Get(0, "Key1")
	if !ok {
		t.Fatalf("Failed to get key")
	}
	if value != "Value1" {
		t.Fatalf("Expected Value1, got %s", value)
	}

	s.Del(0, "Key1")
	_, ok = s.Get(0, "Key1")
	if ok {
		t.Fatalf("Expected key1 to be deleted")
	}
}

func TestExists(t *testing.T) {
	aofChan := make(chan string, 100)

	s := NewStore(aofChan)
	s.Set(0, "Key1", "Value1")
	if !s.Exists(0, "Key1") {
		t.Fatalf("Expected Key1 to exist")
	}
	if s.Exists(0, "Key2") {
		t.Fatalf("Expected Key2 to not exist")
	}
}

func TestSetNX(t *testing.T) {
	aofChan := make(chan string, 100)

	s := NewStore(aofChan)
	if !s.SetNX(0, "Key1", "Value1") {
		t.Fatalf("Expected SETNX to succeed for Key1")
	}
	if s.SetNX(0, "Key1", "Value2") {
		t.Fatalf("Expected SETNX to fail for Key1")
	}
	value, ok := s.Get(0, "Key1")
	if !ok || value != "Value1" {
		t.Fatalf("Expected Value1, got %s", value)
	}
}

func TestExpire(t *testing.T) {
	aofChan := make(chan string, 100)

	s := NewStore(aofChan)
	s.Set(0, "Key1", "Value1")
	if !s.Expire(0, "Key1", 1*time.Second) {
		t.Fatalf("Expected Expire to succeed for Key1")
	}

	time.Sleep(2 * time.Second)
	if s.Exists(0, "Key1") {
		t.Fatalf("Expected Key1 to be expired")
	}
}

func TestIncr(t *testing.T) {
	aofChan := make(chan string, 100)
	s := NewStore(aofChan)

	newValue, err := s.Incr(0, "counter")
	if err != nil {
		t.Fatalf("INCR failed: %v", err)
	}
	// test if value is created and set as '0'++
	if newValue != 1 {
		t.Fatalf("expected 1, got %d", newValue)
	}

	// test if value is incremented
	newValue, err = s.Incr(0, "counter")
	if err != nil {
		t.Fatalf("INCR failed: %v", err)
	}
	if newValue != 2 {
		t.Fatalf("expected 2, got %d", newValue)
	}
}

func TesDecr(t *testing.T) {
	aofChan := make(chan string, 100)
	s := NewStore(aofChan)

	newValue, err := s.Decr(0, "counter")
	if err != nil {
		t.Fatalf("DECR failed: %v", err)
	}
	// test if value is created and set as '0'--
	if newValue != -1 {
		t.Fatalf("expected -1, got %d", newValue)
	}

	// test if value is incremented
	newValue, err = s.Incr(0, "counter")
	if err != nil {
		t.Fatalf("DECR failed: %v", err)
	}
	if newValue != -2 {
		t.Fatalf("expected -2, got %d", newValue)
	}
}

// test Ttl
func TestTtl(t *testing.T) {
	aofChan := make(chan string, 100)
	s := NewStore(aofChan)

	s.Set(0, "Key1", "Value1")
	if !s.Expire(0, "Key1", 4*time.Second) {
		t.Fatalf("Expected Expire to succeed for Key1")
	}
	time.Sleep(1 * time.Second)

	// Test that TTL returns the correct remaining time
	ttl, err := s.TTL(0, "Key1")
	if err != nil {
		t.Fatalf("Expected TTL to succeed for Key1")
	}
	if ttl != 2 {
		t.Fatalf("Expected TTL to be 2 seconds, got %v", ttl)
	}

	time.Sleep(3 * time.Second)

	// Test that TTL returns -2 for expired key
	ttl, err = s.TTL(0, "Key1")
	if err != nil {
		t.Fatalf("Expected TTL to succeed for Key1")
	}
	if ttl != 0 {
		t.Fatalf("Expected TTL to be -2, got %v", ttl)
	}

	s.Set(0, "Key2", "Value2")
	ttl, err = s.TTL(0, "Key2")
	if err != nil {
		t.Fatalf("Expected TTL to succeed for Key2")
	}
	if ttl != -1 {
		t.Fatalf("Expected TTL to be -1, got %v", ttl)
	}

	s.Del(0, "Key2")
	ttl, err = s.TTL(0, "Key2")
	if err != nil {
		t.Fatalf("Expected TTL to succeed for Key2")
	}
	if ttl != -2 {
		t.Fatalf("Expected TTL to be -2, got %v", ttl)
	}
}
