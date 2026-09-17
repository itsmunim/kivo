// Package persistence implements AOF (Append Only File) persistence.
//
// The AOF format is RESP arrays, one per command:
//
//	*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n
//
// On startup, the AOF is replayed by reading each RESP array and
// executing it through the command registry.
package persistence

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/itsmunim/kivo/internal/resp"
)

// AOF manages the append-only file.
type AOF struct {
	path   string
	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex

	// Sync strategy: "always", "everysec", "no".
	sync   string
	ticker *time.Ticker
	stop   chan struct{}
}

// NewAOF creates a new AOF manager.
func NewAOF(path string, sync string) (*AOF, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open AOF: %w", err)
	}

	aof := &AOF{
		path:   path,
		file:   f,
		writer: bufio.NewWriter(f),
		sync:   sync,
		stop:   make(chan struct{}),
	}

	if sync == "everysec" {
		aof.ticker = time.NewTicker(time.Second)
		go aof.periodicSync()
	}

	return aof, nil
}

// Close flushes and closes the AOF file.
func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.ticker != nil {
		a.ticker.Stop()
		close(a.stop)
	}

	if err := a.writer.Flush(); err != nil {
		return err
	}
	if err := a.file.Sync(); err != nil {
		return err
	}
	return a.file.Close()
}

// Sync flushes the AOF buffer to disk.
func (a *AOF) Sync() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.writer.Flush(); err != nil {
		return err
	}
	return a.file.Sync()
}

// Write appends a command to the AOF.
func (a *AOF) Write(args []string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	values := make([]resp.Value, len(args))
	for i, arg := range args {
		values[i] = resp.NewBulkString(arg)
	}
	cmd := resp.NewArray(values...)

	w := resp.NewWriter(a.writer)
	if err := w.WriteValue(cmd); err != nil {
		return fmt.Errorf("aof write: %w", err)
	}
	// resp.NewWriter buffers internally; push those bytes into a.writer
	// (the AOF's own bufio) before the sync-strategy flush below.
	if err := w.Flush(); err != nil {
		return fmt.Errorf("aof write flush: %w", err)
	}

	if a.sync == "always" {
		if err := a.writer.Flush(); err != nil {
			return fmt.Errorf("aof flush: %w", err)
		}
		if err := a.file.Sync(); err != nil {
			return fmt.Errorf("aof fsync: %w", err)
		}
	}

	return nil
}

// Replay reads the AOF and replays commands through the provided handler.
func (a *AOF) Replay(handler func(args []string) error) error {
	f, err := os.Open(a.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open AOF for replay: %w", err)
	}
	defer f.Close()

	reader := resp.NewReader(f)
	for {
		v, err := reader.ReadValue()
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			}
			return fmt.Errorf("aof replay read: %w", err)
		}

		if v.Type() != resp.Array {
			return fmt.Errorf("aof replay: expected array, got %v", v.Type())
		}

		args := v.Array()
		strArgs := make([]string, len(args))
		for i, arg := range args {
			strArgs[i] = arg.String()
		}

		if err := handler(strArgs); err != nil {
			return fmt.Errorf("aof replay command %v: %w", strArgs, err)
		}
	}
}

func (a *AOF) periodicSync() {
	for {
		select {
		case <-a.ticker.C:
			a.mu.Lock()
			_ = a.writer.Flush()
			_ = a.file.Sync()
			a.mu.Unlock()
		case <-a.stop:
			return
		}
	}
}
