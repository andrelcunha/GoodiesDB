package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"com.github.andrelcunha.go-redis-clone/pkg/store"
)

// Server represents a TCP server
type Server struct {
	store *store.Store
}

// NewServer creates a new server
func NewServer(s *store.Store) *Server {
	return &Server{store: s}
}

// Start starts the server
func (s *Server) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection handles a single client connection
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		s.handleCommand(conn, strings.TrimSpace(line))
	}
}

// handleCommand handles a single command
func (s *Server) handleCommand(conn net.Conn, cmd string) {
	parts := strings.Fields(cmd)

	if len(parts) == 0 {
		return
	}

	switch parts[0] {
	case "SET":
		if len(parts) != 3 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'SET' command")
			return
		}
		s.store.Set(parts[1], parts[2])
		fmt.Fprintln(conn, "OK")
	case "GET":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'GET' command")
			return
		}
		value, ok := s.store.Get(parts[1])
		if !ok {
			fmt.Fprintln(conn, "NULL")
		} else {
			fmt.Fprintln(conn, value)
		}
	case "DEL":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'DEL' command")
			return
		}
		s.store.Delete(parts[1])
		fmt.Fprintln(conn, "OK")
	default:
		fmt.Fprintln(conn, "ERR unknown command '"+parts[0]+"'")
	}
}
