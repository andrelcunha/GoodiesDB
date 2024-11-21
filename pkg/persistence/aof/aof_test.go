package aof

import (
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"com.github.andrelcunha.go-redis-clone/pkg/store"
)

func TestRebuildStoreFromAOF(t *testing.T) {
	aofFilename := "test_appendonly.aof"
	os.Remove(aofFilename)
	aofChan := make(chan string, 100)

	// Start the AOF writer
	go AOFWriter(aofChan, aofFilename)

	// Initialize the store with AOF logging
	s := store.NewStore(aofChan)

	dbIndex := 0

	// Set and expire commands
	s.Set(dbIndex, "Key1", "Value1")
	s.Set(dbIndex, "Key2", "Value2")
	s.Expire(dbIndex, "Key1", 3*time.Second)

	// SETNX command
	s.SetNX(dbIndex, "Key3", "Value3")   // Should succeed
	s.SetNX(dbIndex, "Key3", "NewValue") // Should fail because Key3 already exists

	// List commands
	s.LPush(dbIndex, "List1", "Value1", "Value2", "Value3")
	s.RPush(dbIndex, "List1", "Value4")
	s.LPop(dbIndex, "List1", nil)
	s.RPop(dbIndex, "List1", nil)

	// List trimming commands
	s.LTrim(dbIndex, "List1", 1, 2)

	// Rename command
	s.Rename(dbIndex, "Key2", "RenamedKey")

	// Give some time for commands to be written to AOF
	time.Sleep(1 * time.Second)

	// Rebuild state from AOF
	newStore := store.NewStore(aofChan)
	err := RebuildStoreFromAOF(newStore, aofFilename)
	if err != nil {
		t.Fatalf("Failed to rebuild state from AOF: %v", err)
	}

	// Verify Key2 has been renamed to RenamedKey
	value, ok := newStore.Get(dbIndex, "RenamedKey")
	if !ok || value != "Value2" {
		t.Errorf("Expected Value2 for RenamedKey, got %s", value)
		t.Fail()
	}

	// Verify List1 contents
	list, _ := newStore.LRange(dbIndex, "List1", 0, -1)
	expectedList := []string{"Value2"}
	if len(list) != len(expectedList) {
		t.Errorf("Expected list length to be %d, got %d", len(expectedList), len(list))
		t.Fail()
	}
	for i, v := range list {
		if v != expectedList[i] {
			t.Errorf("Expected list[%d] to be %s, got %s", i, expectedList[i], v)
			t.Fail()
		}
	}

	// Wait for the key to expire
	time.Sleep(4 * time.Second)

	// Verify Key1 exists after it expires
	if newStore.Exists(dbIndex, "Key1") {
		t.Errorf("Expected Key1 to be expired after waiting more than 3 seconds")
		t.Fail()
	}

	// Clean up the AOF file
	os.Remove(aofFilename)
}

// Test aofRename
func TestAofRename(t *testing.T) {
	cmd := "RENAME 0 Key1 newName"
	parts, s, dbIndex := prepareCmdTest(cmd)

	s.Set(dbIndex, "Key1", "value1")

	aofRename(parts, s, dbIndex)
	value, ok := s.Get(dbIndex, "newName")
	if !ok || value != "value1" {
		t.Fatalf("Expeted 'value1, got %s", value)
	}

}

func prepareCmdTest(cmd string) ([]string, *store.Store, int) {
	aofChan := make(chan string, 100)
	s := store.NewStore(aofChan)

	parts := strings.Split(cmd, " ")

	dbIndex, _ := strconv.Atoi(parts[1])
	return parts, s, dbIndex
}
