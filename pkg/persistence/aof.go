package persistence

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"com.github.andrelcunha.go-redis-clone/pkg/store"
)

// AOFWriter writes commands to a file
func AOFWriter(aofChan chan string, filename string) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open AOF file: %v", err)
	}
	defer file.Close()

	for cmd := range aofChan {
		_, err := file.WriteString(cmd + "\n")
		if err != nil {
			log.Fatalf("Failed to write to AOF file: %v", err)
		}
	}
}

// RebuildStoreFromAOF rebuilds the store from the AOF file
func RebuildStoreFromAOF(s *store.Store, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create scanner to read the AOF file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cmd := scanner.Text()
		parts := strings.Split(cmd, " ")
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {

		case "SET":
			if len(parts) == 3 {
				s.Set(parts[1], parts[2])
			}

		case "DEL":
			if len(parts) == 2 {
				s.Del(parts[1])
			}

		case "SETNX":
			if len(parts) == 3 {
				s.SetNX(parts[1], parts[2])
			}

		case "EXPIRE":
			if len(parts) == 3 {
				key := parts[1]
				ttl, err := strconv.Atoi(parts[2])
				if err == nil {
					duration := time.Duration(ttl) * time.Second
					s.Expire(key, duration)
				}
			}

		default:
			log.Printf("Unknown command: %s", cmd)
		}
	}

	return scanner.Err()
}
