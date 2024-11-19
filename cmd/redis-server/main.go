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

	// Start the server
	s := store.NewStore()
	srv := server.NewServer(s)

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

	// Start the server
	if err := srv.Start(":6379"); err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
