package store

import (
	"testing"
)

func TestStore(t *testing.T) {
	s := NewStore()

	s.Set("Key1", "Value1")
	value, ok := s.Get("Key1")
	if !ok {
		t.Fatalf("Failed to get key")
	}
	if value != "Value1" {
		t.Fatalf("Expected Value1, got %s", value)
	}

	s.Delete("Key1")
	_, ok = s.Get("Key1")
	if ok {
		t.Fatalf("Expected key1 to be deleted")
	}
}

func TestExists(t *testing.T) {
	s := NewStore()
	s.Set("Key1", "Value1")
	if !s.Exists("Key1") {
		t.Fatalf("Expected Key1 to exist")
	}
	if s.Exists("Key2") {
		t.Fatalf("Expected Key2 to not exist")
	}
}

func TestSetNX(t *testing.T) {
	s := NewStore()
	if !s.SetNX("Key1", "Value1") {
		t.Fatalf("Expected SETNX to succeed for Key1")
	}
	if s.SetNX("Key1", "Value2") {
		t.Fatalf("Expected SETNX to fail for Key1")
	}
	value, ok := s.Get("Key1")
	if !ok || value != "Value1" {
		t.Fatalf("Expected Value1, got %s", value)
	}
}
