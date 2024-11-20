package persistence

import (
	"os"
	"testing"
	"time"

	"com.github.andrelcunha.go-redis-clone/pkg/store"
)

func TestSaveLoadSnapshot(t *testing.T) {
	// Create a temporary AOF file
	aofFilename := "test_appendonly.aof"
	aofChan := make(chan string, 100)

	// Start the AOF writer
	go AOFWriter(aofChan, aofFilename)

	// Initialize a new store with the AOF file
	s := store.NewStore(aofChan)

	s.Set("Key1", "Value1")
	s.Set("Key2", "Value2")
	s.Expire("Key1", 3*time.Second)

	err := SaveSnapshot(s, "test_snapshot.gob")
	if err != nil {
		t.Fatalf("Failed to save snapshot: %v", err)
	}

	newStore := store.NewStore(aofChan)
	err = LoadSnapshot(newStore, "test_snapshot.gob")
	if err != nil {
		t.Fatalf("Failed to load snapshot: %v", err)
	}

	// Verify Key1 exists before it expires
	value, ok := newStore.Get("Key1")
	if !ok || value != "Value1" {
		t.Fatalf("Expected Value1, got %s", value)
	}

	// Verify Key2 exists before it expires
	value, ok = newStore.Get("Key2")
	if !ok || value != "Value2" {
		t.Fatalf("Expected Value2, got %s", value)
	}

	// Wait for the key to expire
	time.Sleep(4 * time.Second)

	// Verify Key1 exists after it expires
	if newStore.Exists("Key1") {
		t.Fatalf("Expected Key1 to be expered after snapshot load an waiting more than 3 seconds")
	}

	// Clean up the snapshot file
	err = os.Remove("test_snapshot.gob")

	// Clean up the AOF file
	os.Remove(aofFilename)

}
