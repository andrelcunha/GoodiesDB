package main

import (
	"fmt"
	"time"

	"com.github.andrelcunha.go-redis-clone/pkg/persistence"
	"com.github.andrelcunha.go-redis-clone/pkg/persistence/aof"
	"com.github.andrelcunha.go-redis-clone/pkg/server"
	"com.github.andrelcunha.go-redis-clone/pkg/store"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file, using default values.")
	}
	// Set up configuration
	config := server.NewConfig()
	config.LoadFromEnv()

	fmt.Println("Starting Redis Clone Server...")

	// Set up AOF
	aofChan := make(chan string)

	// Set up the store
	s := store.NewStore(aofChan)

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
	go aof.AOFWriter(aofChan, "appendonly.aof")

	// Start the server
	if err := srv.Start(":6379"); err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
