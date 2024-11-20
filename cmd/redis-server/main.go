package main

import (
	"fmt"
	"time"

	"com.github.andrelcunha.go-redis-clone/pkg/persistence"
	"com.github.andrelcunha.go-redis-clone/pkg/server"
	"com.github.andrelcunha.go-redis-clone/pkg/store"
)

func main() {
	fmt.Println("Starting Redis Clone Server...")

	// Set up AOF
	aofChan := make(chan string)

	// Set up the store
	s := store.NewStore(aofChan)

	// Set up configuration
	config := server.NewConfig()

	// Start the server
	srv := server.NewServer(s, config)

	//Load snapshot on startup
	if err := persistence.LoadSnapshot(s, "snapshot.gob"); err != nil {
		fmt.Println("No snapshot found, starting with empty store.")
	}

	// Periodic snapshot saving
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			if err := persistence.SaveSnapshot(s, "snapshot.gob"); err != nil {
				fmt.Println("Error saving snapshot: ", err)
			} else {
				fmt.Println("Snapshot saved.")
			}
		}
	}()

	// Start the AOF writer
	go persistence.AOFWriter(aofChan, "appendonly.aof")

	// Start the server
	if err := srv.Start(":6379"); err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
