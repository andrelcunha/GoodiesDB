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

		dbIndex, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("Invalid database index: %s", parts[1])
			continue
		}

		switch parts[0] {

		case "SET":
			if len(parts) == 4 {
				s.Set(dbIndex, parts[2], parts[3])
			}

		case "DEL":
			if len(parts) == 3 {
				s.Del(dbIndex, parts[2])
			}

		case "SETNX":
			if len(parts) == 4 {
				s.SetNX(dbIndex, parts[2], parts[3])
			}

		case "EXPIRE":
			if len(parts) == 4 {
				key := parts[2]
				ttl, err := strconv.Atoi(parts[3])
				if err == nil {
					duration := time.Duration(ttl) * time.Second
					s.Expire(dbIndex, key, duration)
				}
			}

		case "LPUSH":
			if len(parts) >= 4 {
				s.LPush(dbIndex, parts[2], parts[3:]...)
			}

		case "RPUSH":
			if len(parts) >= 4 {
				s.RPush(dbIndex, parts[2], parts[3:]...)
			}

		case "LPOP":
			if len(parts) == 4 {
				count, err := strconv.Atoi(parts[3])
				if err == nil {
					s.LPop(dbIndex, parts[2], &count)
				}
			}

		case "RPOP":
			if len(parts) == 4 {
				count, err := strconv.Atoi(parts[3])
				if err == nil {
					s.RPop(dbIndex, parts[2], &count)
				}
			}

		default:
			log.Printf("Unknown command: %s", cmd)
		}
	}

	return scanner.Err()
}
