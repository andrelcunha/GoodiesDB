// pkg/persistence/aof_test.go
package persistence

import (
	"os"
	"testing"
	"time"

	"com.github.andrelcunha.go-redis-clone/pkg/store"
)

func TestAOFRecovery(t *testing.T) {
	aofFilename := "test_appendonly.aof"
	aofChan := make(chan string, 100)

	// Start the AOF writer
	go AOFWriter(aofChan, aofFilename)

	// Initialize the store with AOF logging
	s := store.NewStore(aofChan)

	s.Set("Key1", "Value1")
	s.Set("Key2", "Value2")
	s.Expire("Key1", 3*time.Second)

	// Give some time for commands to be written to AOF
	time.Sleep(1 * time.Second)

	// Rebuild state from AOF
	newStore := store.NewStore(aofChan)
	err := RebuildStoreFromAOF(newStore, aofFilename)
	if err != nil {
		t.Fatalf("Failed to rebuild state from AOF: %v", err)
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
		t.Fatalf("Expected Key1 to be expired after waiting more than 3 seconds")
	}

	// Clean up the AOF file
	os.Remove(aofFilename)
}
