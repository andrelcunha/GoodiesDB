package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"com.github.andrelcunha.go-redis-clone/pkg/store"
)

// Server represents a TCP server
type Server struct {
	store                    *store.Store
	config                   *Config
	mu                       sync.Mutex
	authenticatedConnections map[net.Conn]bool
	connectionDbs            map[net.Conn]int
}

// NewServer creates a new server
func NewServer(store *store.Store, config *Config) *Server {
	return &Server{
		store:                    store,
		config:                   config,
		authenticatedConnections: make(map[net.Conn]bool),
		connectionDbs:            make(map[net.Conn]int),
	}
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

	//check authentication
	if !s.isAuthenticates(conn) && parts[0] != "AUTH" {
		fmt.Fprintln(conn, "NOAUTH Authentication required.")
		return
	}

	dbIndex := s.getCurrentDb(conn)

	switch parts[0] {

	case "AUTH":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'AUTH' command")
			return
		}
		if parts[1] == s.config.Password {
			s.mu.Lock()
			s.authenticatedConnections[conn] = true
			s.mu.Unlock()
			fmt.Fprintln(conn, "OK")
		} else {
			fmt.Fprintln(conn, "ERR invalid password")
		}

	case "SET":
		if len(parts) != 3 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'SET' command")
			return
		}
		s.store.Set(dbIndex, parts[1], parts[2])
		fmt.Fprintln(conn, "OK")

	case "GET":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'GET' command")
			return
		}
		value, ok := s.store.Get(dbIndex, parts[1])
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
		s.store.Del(dbIndex, parts[1])
		fmt.Fprintln(conn, "OK")

	case "EXISTS":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'EXISTS' command")
			return
		}
		exists := s.store.Exists(dbIndex, parts[1])
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
		if s.store.SetNX(dbIndex, parts[1], parts[2]) {
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
		if s.store.Expire(dbIndex, key, duration) {
			fmt.Fprintln(conn, 1)
		} else {
			fmt.Fprintln(conn, 0)
		}

	case "INCR":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'INCR' command")
			return
		}
		newValue, err := s.store.Incr(dbIndex, parts[1])
		if err != nil {
			fmt.Fprintln(conn, "ERR ", err.Error())
			return
		}
		fmt.Fprintln(conn, newValue)

	case "DECR":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'DECR' command")
			return
		}
		newValue, err := s.store.Decr(dbIndex, parts[1])
		if err != nil {
			fmt.Fprintln(conn, "ERR ", err.Error())
			return
		}
		fmt.Fprintln(conn, newValue)

	case "TTL":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'TTL' command")
			return
		}
		ttl, err := s.store.TTL(dbIndex, parts[1])
		if err != nil {
			fmt.Fprintln(conn, "ERR ", err.Error())
			return
		}
		fmt.Fprintln(conn, ttl)

	case "SELECT":
		if len(parts) != 2 {
			fmt.Fprintln(conn, "ERR wrong number of arguments for 'SELECT' command")
			return
		}
		dbIndex, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Fprintln(conn, "ERR invalid DB index")
			return
		}
		s.selectDb(conn, dbIndex)
		fmt.Fprintln(conn, "OK")

	default:
		fmt.Fprintln(conn, "ERR unknown command '"+parts[0]+"'")
		fmt.Fprintln(conn, "Available commands: AUTH, SET, GET, DEL, EXISTS, SETNX, EXPIRE, INCR, DECR, TTL")
	}
}

func (s *Server) isAuthenticates(conn net.Conn) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.authenticatedConnections[conn]
}

func (s *Server) getCurrentDb(conn net.Conn) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	db, ok := s.connectionDbs[conn]
	if !ok {
		db = 0
		s.connectionDbs[conn] = db
	}
	return db
}

func (s *Server) selectDb(conn net.Conn, dbIndex int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connectionDbs[conn] = dbIndex
}
