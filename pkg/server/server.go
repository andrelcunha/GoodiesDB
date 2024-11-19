package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

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
		s.store.Del(parts[1])
		fmt.Fprintln(conn, "OK")

	case "EXISTS":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'EXISTS' command")
			return
		}
		exists := s.store.Exists(parts[1])
		if exists {
			fmt.Fprintln(conn, 1)
		} else {
			fmt.Fprintln(conn, 0)
		}

	case "SETNX":
		if len(parts) != 3 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'SETNX' command")
			return
		}
		if s.store.SetNX(parts[1], parts[2]) {
			fmt.Fprintln(conn, 1)
		} else {
			fmt.Fprintln(conn, 0)
		}

	case "EXPIRE":
		if len(parts) != 3 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'EXPIRE' command")
			return
		}
		key := parts[1]
		ttl, err := strconv.Atoi(parts[2])
		if err != nil {
			fmt.Fprintln(conn, "ERR invalid TTL")
			return
		}
		duration := time.Duration(ttl) * time.Second
		if s.store.Expire(key, duration) {
			fmt.Fprintln(conn, 1)
		} else {
			fmt.Fprintln(conn, 0)
		}

	default:
		fmt.Fprintln(conn, "ERR unknown command '"+parts[0]+"'")
	}
}
