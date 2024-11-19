package persistence

import (
	"log"
	"os"
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
