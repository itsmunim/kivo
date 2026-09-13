// Package server implements the TCP server and connection handling.
package server

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/itsmunim/kivo/internal/commands"
	"github.com/itsmunim/kivo/internal/config"
	"github.com/itsmunim/kivo/internal/persistence"
	"github.com/itsmunim/kivo/internal/resp"
	"github.com/itsmunim/kivo/internal/store"
)

// Server is the TCP server.
type Server struct {
	config   config.Config
	engine   *store.Engine
	registry *commands.Registry
	aof      *persistence.AOF
	listener net.Listener
	conns    sync.Map // map[net.Conn]struct{}
	active   int64    // atomic counter for active connections
	shutdown int32    // atomic: 1 if shutting down
	wg       sync.WaitGroup
}

// New creates a new server.
func New(cfg config.Config, engine *store.Engine, registry *commands.Registry, aof *persistence.AOF) *Server {
	return &Server{
		config:   cfg,
		engine:   engine,
		registry: registry,
		aof:      aof,
	}
}

// Start begins accepting connections.
func (s *Server) Start() error {
	ln := s.listener
	if ln == nil {
		var err error
		ln, err = net.Listen("tcp", s.config.Addr)
		if err != nil {
			return fmt.Errorf("listen %s: %w", s.config.Addr, err)
		}
		s.listener = ln
	}
	fmt.Printf("kivo listening on %s\n", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			if atomic.LoadInt32(&s.shutdown) == 1 {
				return nil
			}
			fmt.Printf("accept error: %v\n", err)
			continue
		}
		s.wg.Add(1)
		go s.handleConn(conn)
	}
}

// Stop initiates graceful shutdown.
func (s *Server) Stop() {
	atomic.StoreInt32(&s.shutdown, 1)
	if s.listener != nil {
		s.listener.Close()
	}
	// Close all active connections to unblock their handlers.
	// Without this, wg.Wait() below would hang forever if a client
	// (e.g. interactive redis-cli) is still connected.
	s.conns.Range(func(k, v interface{}) bool {
		if conn, ok := k.(net.Conn); ok {
			_ = conn.Close()
		}
		return true
	})
	// Wait for all connection handlers to finish.
	s.wg.Wait()
	// Flush AOF to disk before returning. Must happen after handlers
	// finish so no late writes are lost.
	if s.aof != nil {
		_ = s.aof.Sync()
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	s.conns.Store(conn, struct{}{})
	defer s.conns.Delete(conn)

	atomic.AddInt64(&s.active, 1)
	defer atomic.AddInt64(&s.active, -1)

	reader := resp.NewReader(conn)
	writer := resp.NewWriter(conn)

	for {
		v, err := reader.ReadValue()
		if err != nil {
			if err != io.EOF {
				fmt.Printf("read error from %s: %v\n", conn.RemoteAddr(), err)
			}
			return
		}

		if v.Type() != resp.Array {
			_ = writer.WriteValue(resp.NewError("ERR unknown command"))
			continue
		}

		args := v.Array()
		if len(args) == 0 {
			_ = writer.WriteValue(resp.NewError("ERR unknown command"))
			continue
		}

		cmdName := args[0].String()
		cmdArgs := args[1:]

		// Special case: QUIT closes the connection.
		if cmdName == "QUIT" {
			_ = writer.WriteValue(resp.NewSimpleString("OK"))
			return
		}

		handler, ok := s.registry.Get(cmdName)
		if !ok {
			_ = writer.WriteValue(resp.NewError(fmt.Sprintf("ERR unknown command '%s'", cmdName)))
			continue
		}

		result := handler(s.engine, cmdArgs)

		// Write AOF for write commands.
		if s.aof != nil && isWriteCommand(cmdName) {
			cmdArgsStr := make([]string, len(args))
			for i, arg := range args {
				cmdArgsStr[i] = arg.String()
			}
			if err := s.aof.Write(cmdArgsStr); err != nil {
				fmt.Printf("aof write error: %v\n", err)
			}
		}

		if err := writer.WriteValue(result); err != nil {
			fmt.Printf("write error to %s: %v\n", conn.RemoteAddr(), err)
			return
		}
	}
}

// ActiveConns returns the number of active connections.
func (s *Server) ActiveConns() int64 {
	return atomic.LoadInt64(&s.active)
}

// isWriteCommand returns true if the command modifies data.
// The command name is normalized to uppercase for the check.
func isWriteCommand(cmd string) bool {
	switch strings.ToUpper(cmd) {
	case "SET", "DEL", "EXPIRE", "PEXPIRE", "PERSIST",
		"APPEND", "INCR", "DECR", "INCRBY", "DECRBY", "MSET",
		"LPUSH", "RPUSH", "LPOP", "RPOP", "LREM", "LTRIM",
		"SADD", "SREM", "SPOP",
		"HSET", "HDEL",
		"ZADD", "ZREM", "ZINCRBY", "ZREMRANGEBYSCORE", "ZREMRANGEBYRANK",
		"RENAME", "RENAMENX", "FLUSHDB":
		return true
	default:
		return false
	}
}
