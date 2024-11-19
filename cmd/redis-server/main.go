package main

import (
	"fmt"

	"com.github.andrelcunha.go-redis-clone/pkg/server"
	"com.github.andrelcunha.go-redis-clone/pkg/store"
)

func main() {
	fmt.Println("Starting Redis Clone Server...")

	// Start the server
	s := store.NewStore()
	srv := server.NewServer(s)

	if err := srv.Start(":6379"); err != nil {
		fmt.Println("Error starting server: ", err)
		return
	}
}
