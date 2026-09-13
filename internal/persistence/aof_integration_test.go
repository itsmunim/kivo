package persistence

import (
	"path/filepath"
	"testing"

	"github.com/itsmunim/kivo/internal/commands"
	"github.com/itsmunim/kivo/internal/resp"
	"github.com/itsmunim/kivo/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAOFEndToEndPersistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "appendonly.aof")

	// Phase 1: Create AOF with "always" sync, write SET, close
	aof1, err := NewAOF(path, "always")
	require.NoError(t, err)

	err = aof1.Write([]string{"SET", "foo2", "bar2"})
	require.NoError(t, err)

	err = aof1.Close()
	require.NoError(t, err)

	// Phase 2: Create fresh engine, replay AOF
	engine := store.NewEngineWithMaxMemory(0)
	defer engine.Stop()
	registry := commands.NewRegistry()

	aof2, err := NewAOF(path, "always")
	require.NoError(t, err)
	defer aof2.Close()

	count := 0
	err = aof2.Replay(func(args []string) error {
		handler, ok := registry.Get(args[0])
		if !ok {
			return nil // skip unknown commands
		}
		cmdArgs := make([]resp.Value, len(args)-1)
		for i := 1; i < len(args); i++ {
			cmdArgs[i-1] = resp.NewBulkString(args[i])
		}
		handler(engine, cmdArgs)
		count++
		return nil
	})
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Phase 3: Verify key exists
	val, ok := engine.Get("foo2")
	assert.True(t, ok)
	assert.Equal(t, "bar2", val)
}
