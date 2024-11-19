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
