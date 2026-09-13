package server

import (
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/itsmunim/kivo/internal/commands"
	"github.com/itsmunim/kivo/internal/config"
	"github.com/itsmunim/kivo/internal/persistence"
	"github.com/itsmunim/kivo/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStopWithIdleConnection verifies Stop() does not hang when a client
// connection is still open (e.g. an interactive redis-cli sitting idle),
// and that buffered AOF data is flushed to disk during shutdown.
func TestStopWithIdleConnection(t *testing.T) {
	engine := store.NewEngine()
	defer engine.Stop()

	registry := commands.NewRegistry()

	dir := t.TempDir()
	aofPath := filepath.Join(dir, "appendonly.aof")
	aof, err := persistence.NewAOF(aofPath, "everysec")
	require.NoError(t, err)
	defer aof.Close()

	cfg := config.Config{Addr: "127.0.0.1:0"}
	srv := New(cfg, engine, registry, aof)

	ln, err := net.Listen("tcp", cfg.Addr)
	require.NoError(t, err)
	srv.listener = ln

	go func() { _ = srv.Start() }()
	addr := ln.Addr().String()

	// A real client command.
	val := sendCommand(t, addr, "SET", "foo2", "bar2")
	require.Equal(t, "OK", val.String())

	// An idle connection that stays open after the command.
	idle, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer idle.Close()

	// Stop() must return promptly despite the open idle connection.
	done := make(chan struct{})
	go func() {
		srv.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Stop() hung with an idle connection open (regression)")
	}

	// The engine must still hold the key after replay-free shutdown.
	v, ok := engine.Get("foo2")
	assert.True(t, ok)
	assert.Equal(t, "bar2", v)

	// The AOF file must contain the SET command flushed during Stop().
	data, err := os.ReadFile(aofPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "SET")
	assert.Contains(t, string(data), "foo2")
	assert.Contains(t, string(data), "bar2")
}

// TestAOFWritesLowerCaseCommands is a regression test for the bug where
// clients that send lowercase commands (rdcli, ioredis, redis-py, go-redis
// all send lowercase) never triggered the AOF write hook, because
// isWriteCommand compared the raw command name against uppercase literals.
// As a result the store updated fine but the AOF stayed empty and data
// was lost on restart.
func TestAOFWritesLowerCaseCommands(t *testing.T) {
	engine := store.NewEngine()
	defer engine.Stop()

	registry := commands.NewRegistry()

	dir := t.TempDir()
	aofPath := filepath.Join(dir, "appendonly.aof")
	aof, err := persistence.NewAOF(aofPath, "always")
	require.NoError(t, err)
	defer aof.Close()

	cfg := config.Config{Addr: "127.0.0.1:0"}
	srv := New(cfg, engine, registry, aof)

	ln, err := net.Listen("tcp", cfg.Addr)
	require.NoError(t, err)
	srv.listener = ln

	go func() { _ = srv.Start() }()
	addr := ln.Addr().String()

	// Send the command in lowercase, exactly like rdcli and most clients do.
	val := sendCommand(t, addr, "set", "foo2", "bar2")
	require.Equal(t, "OK", val.String())

	// Confirm the store has it (proves the command executed).
	v, ok := engine.Get("foo2")
	require.True(t, ok)
	assert.Equal(t, "bar2", v)

	// The AOF must contain the lowercase write as well.
	aof.Sync()
	data, err := os.ReadFile(aofPath)
	require.NoError(t, err)
	assert.Contains(t, string(data), "set")
	assert.Contains(t, string(data), "foo2")
	assert.Contains(t, string(data), "bar2")
}
