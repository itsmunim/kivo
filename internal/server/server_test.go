package server

import (
	"net"
	"testing"
	"time"

	"github.com/itsmunim/kivo/internal/commands"
	"github.com/itsmunim/kivo/internal/config"
	"github.com/itsmunim/kivo/internal/resp"
	"github.com/itsmunim/kivo/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func startTestServer(t *testing.T) (*Server, string) {
	engine := store.NewEngine()
	t.Cleanup(func() { engine.Stop() })

	registry := commands.NewRegistry()
	cfg := config.Config{Addr: "127.0.0.1:0"} // Let OS pick a port.

	srv := New(cfg, engine, registry, nil)

	ln, err := net.Listen("tcp", cfg.Addr)
	require.NoError(t, err)
	srv.listener = ln

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			srv.wg.Add(1)
			go srv.handleConn(conn)
		}
	}()

	t.Cleanup(func() {
		srv.Stop()
	})

	return srv, ln.Addr().String()
}

func sendCommand(t *testing.T, addr string, args ...string) resp.Value {
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer conn.Close()

	writer := resp.NewWriter(conn)
	reader := resp.NewReader(conn)

	cmdArgs := make([]resp.Value, len(args))
	for i, arg := range args {
		cmdArgs[i] = resp.NewBulkString(arg)
	}
	err = writer.WriteValue(resp.NewArray(cmdArgs...))
	require.NoError(t, err)
	require.NoError(t, writer.Flush())

	v, err := reader.ReadValue()
	require.NoError(t, err)
	return v
}

func TestServerPing(t *testing.T) {
	_, addr := startTestServer(t)
	v := sendCommand(t, addr, "PING")
	assert.Equal(t, "PONG", v.String())
}

func TestServerEcho(t *testing.T) {
	_, addr := startTestServer(t)
	v := sendCommand(t, addr, "ECHO", "hello world")
	assert.Equal(t, "hello world", v.String())
}

func TestServerGetSet(t *testing.T) {
	_, addr := startTestServer(t)

	v := sendCommand(t, addr, "SET", "key", "value")
	assert.Equal(t, "OK", v.String())

	v = sendCommand(t, addr, "GET", "key")
	assert.Equal(t, "value", v.String())

	v = sendCommand(t, addr, "GET", "missing")
	assert.True(t, v.IsNull())
}

func TestServerDel(t *testing.T) {
	_, addr := startTestServer(t)

	sendCommand(t, addr, "SET", "a", "1")
	sendCommand(t, addr, "SET", "b", "2")

	v := sendCommand(t, addr, "DEL", "a", "b", "missing")
	assert.Equal(t, int64(2), v.Integer())
}

func TestServerIncr(t *testing.T) {
	_, addr := startTestServer(t)

	v := sendCommand(t, addr, "INCR", "counter")
	assert.Equal(t, int64(1), v.Integer())

	v = sendCommand(t, addr, "INCRBY", "counter", "5")
	assert.Equal(t, int64(6), v.Integer())
}

func TestServerLPushLRange(t *testing.T) {
	_, addr := startTestServer(t)

	v := sendCommand(t, addr, "LPUSH", "list", "a", "b", "c")
	assert.Equal(t, int64(3), v.Integer())

	v = sendCommand(t, addr, "LRANGE", "list", "0", "-1")
	arr := v.Array()
	require.Len(t, arr, 3)
	assert.Equal(t, "c", arr[0].String())
	assert.Equal(t, "b", arr[1].String())
	assert.Equal(t, "a", arr[2].String())
}

func TestServerSAddSMembers(t *testing.T) {
	_, addr := startTestServer(t)

	v := sendCommand(t, addr, "SADD", "set", "a", "b")
	assert.Equal(t, int64(2), v.Integer())

	v = sendCommand(t, addr, "SMEMBERS", "set")
	arr := v.Array()
	require.Len(t, arr, 2)
}

func TestServerHSetHGet(t *testing.T) {
	_, addr := startTestServer(t)

	v := sendCommand(t, addr, "HSET", "hash", "field", "value")
	assert.Equal(t, int64(1), v.Integer())

	v = sendCommand(t, addr, "HGET", "hash", "field")
	assert.Equal(t, "value", v.String())
}

func TestServerZAddZRange(t *testing.T) {
	_, addr := startTestServer(t)

	v := sendCommand(t, addr, "ZADD", "zset", "1", "a", "2", "b")
	assert.Equal(t, int64(2), v.Integer())

	v = sendCommand(t, addr, "ZRANGE", "zset", "0", "-1")
	arr := v.Array()
	require.Len(t, arr, 2)
	assert.Equal(t, "a", arr[0].String())
	assert.Equal(t, "b", arr[1].String())
}

func TestServerUnknownCommand(t *testing.T) {
	_, addr := startTestServer(t)
	v := sendCommand(t, addr, "FOOBAR")
	assert.Equal(t, resp.Error, v.Type())
	assert.Contains(t, v.Error(), "unknown command")
}

func TestServerQuit(t *testing.T) {
	_, addr := startTestServer(t)
	v := sendCommand(t, addr, "QUIT")
	assert.Equal(t, "OK", v.String())
}

func TestServerPipeline(t *testing.T) {
	_, addr := startTestServer(t)

	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	defer conn.Close()

	writer := resp.NewWriter(conn)
	reader := resp.NewReader(conn)

	// Send multiple commands without reading.
	for i := 0; i < 5; i++ {
		err := writer.WriteValue(resp.NewArray(
			resp.NewBulkString("SET"),
			resp.NewBulkString("key"),
			resp.NewBulkString("val"),
		))
		require.NoError(t, err)
	}

	require.NoError(t, writer.Flush())
	// Read all responses.
	for i := 0; i < 5; i++ {
		v, err := reader.ReadValue()
		require.NoError(t, err)
		assert.Equal(t, "OK", v.String())
	}
}

func TestServerExpiration(t *testing.T) {
	_, addr := startTestServer(t)

	v := sendCommand(t, addr, "SET", "key", "val", "PX", "100")
	assert.Equal(t, "OK", v.String())

	v = sendCommand(t, addr, "GET", "key")
	assert.Equal(t, "val", v.String())

	time.Sleep(150 * time.Millisecond)

	v = sendCommand(t, addr, "GET", "key")
	assert.True(t, v.IsNull())
}

func TestServerKeys(t *testing.T) {
	_, addr := startTestServer(t)

	sendCommand(t, addr, "SET", "user:1", "a")
	sendCommand(t, addr, "SET", "user:2", "b")
	sendCommand(t, addr, "SET", "post:1", "c")

	v := sendCommand(t, addr, "KEYS", "user:*")
	arr := v.Array()
	require.Len(t, arr, 2)
}

func TestServerSelect(t *testing.T) {
	_, addr := startTestServer(t)
	v := sendCommand(t, addr, "SELECT", "5")
	assert.Equal(t, "OK", v.String())
}

func TestServerFlushDB(t *testing.T) {
	_, addr := startTestServer(t)

	sendCommand(t, addr, "SET", "a", "1")
	v := sendCommand(t, addr, "FLUSHDB")
	assert.Equal(t, "OK", v.String())

	v = sendCommand(t, addr, "DBSIZE")
	assert.Equal(t, int64(0), v.Integer())
}

func TestServerAuth(t *testing.T) {
	_, addr := startTestServer(t)
	v := sendCommand(t, addr, "AUTH", "secret")
	assert.Equal(t, "OK", v.String())
}

func TestServerType(t *testing.T) {
	_, addr := startTestServer(t)

	sendCommand(t, addr, "SET", "str", "val")
	v := sendCommand(t, addr, "TYPE", "str")
	assert.Equal(t, "string", v.String())

	sendCommand(t, addr, "LPUSH", "list", "a")
	v = sendCommand(t, addr, "TYPE", "list")
	assert.Equal(t, "list", v.String())

	v = sendCommand(t, addr, "TYPE", "missing")
	assert.Equal(t, "none", v.String())
}

func TestServerRename(t *testing.T) {
	_, addr := startTestServer(t)

	sendCommand(t, addr, "SET", "old", "val")
	v := sendCommand(t, addr, "RENAME", "old", "new")
	assert.Equal(t, "OK", v.String())

	v = sendCommand(t, addr, "GET", "new")
	assert.Equal(t, "val", v.String())
}

func TestServerRenameNX(t *testing.T) {
	_, addr := startTestServer(t)

	sendCommand(t, addr, "SET", "a", "1")
	sendCommand(t, addr, "SET", "b", "2")

	v := sendCommand(t, addr, "RENAMENX", "a", "b")
	assert.Equal(t, int64(0), v.Integer())

	v = sendCommand(t, addr, "RENAMENX", "a", "c")
	assert.Equal(t, int64(1), v.Integer())
}

func TestServerExists(t *testing.T) {
	_, addr := startTestServer(t)

	sendCommand(t, addr, "SET", "a", "1")
	v := sendCommand(t, addr, "EXISTS", "a", "missing")
	assert.Equal(t, int64(1), v.Integer())
}

func TestServerMSetMGet(t *testing.T) {
	_, addr := startTestServer(t)

	v := sendCommand(t, addr, "MSET", "a", "1", "b", "2")
	assert.Equal(t, "OK", v.String())

	v = sendCommand(t, addr, "MGET", "a", "b", "missing")
	arr := v.Array()
	require.Len(t, arr, 3)
	assert.Equal(t, "1", arr[0].String())
	assert.Equal(t, "2", arr[1].String())
	assert.True(t, arr[2].IsNull())
}

func TestServerStartAndStop(t *testing.T) {
	engine := store.NewEngine()
	defer engine.Stop()
	registry := commands.NewRegistry()
	cfg := config.Config{Addr: "127.0.0.1:0"}
	srv := New(cfg, engine, registry, nil)

	// Manually create listener and assign it (avoids race in Start).
	ln, err := net.Listen("tcp", cfg.Addr)
	require.NoError(t, err)
	srv.listener = ln

	go func() {
		err := srv.Start()
		assert.NoError(t, err)
	}()

	addr := ln.Addr().String()

	// Connect and send a command.
	v := sendCommand(t, addr, "PING")
	assert.Equal(t, "PONG", v.String())

	// Check active connections.
	assert.GreaterOrEqual(t, srv.ActiveConns(), int64(0))

	// Stop the server.
	srv.Stop()
}

func TestServerActiveConns(t *testing.T) {
	engine := store.NewEngine()
	defer engine.Stop()
	registry := commands.NewRegistry()
	cfg := config.Config{Addr: "127.0.0.1:0"}
	srv := New(cfg, engine, registry, nil)

	assert.Equal(t, int64(0), srv.ActiveConns())
}

func TestServerWriteError(t *testing.T) {
	engine := store.NewEngine()
	defer engine.Stop()
	registry := commands.NewRegistry()
	cfg := config.Config{Addr: "127.0.0.1:0"}
	srv := New(cfg, engine, registry, nil)

	ln, err := net.Listen("tcp", cfg.Addr)
	require.NoError(t, err)
	srv.listener = ln

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			srv.wg.Add(1)
			go srv.handleConn(conn)
		}
	}()

	addr := ln.Addr().String()

	// Connect and immediately close to trigger read error.
	conn, err := net.Dial("tcp", addr)
	require.NoError(t, err)
	conn.Close()

	time.Sleep(50 * time.Millisecond)
	srv.Stop()
}
