package persistence

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAOFWriteAndReplay(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "appendonly.aof")

	aof, err := NewAOF(path, "no")
	require.NoError(t, err)

	// Write some commands.
	err = aof.Write([]string{"SET", "key1", "value1"})
	require.NoError(t, err)
	err = aof.Write([]string{"SET", "key2", "value2"})
	require.NoError(t, err)
	err = aof.Write([]string{"DEL", "key1"})
	require.NoError(t, err)

	require.NoError(t, aof.Close())

	// Replay.
	var replayed [][]string
	aof2, err := NewAOF(path, "no")
	require.NoError(t, err)
	defer aof2.Close()

	err = aof2.Replay(func(args []string) error {
		replayed = append(replayed, args)
		return nil
	})
	require.NoError(t, err)

	require.Len(t, replayed, 3)
	assert.Equal(t, []string{"SET", "key1", "value1"}, replayed[0])
	assert.Equal(t, []string{"SET", "key2", "value2"}, replayed[1])
	assert.Equal(t, []string{"DEL", "key1"}, replayed[2])
}

func TestAOFReplayEmpty(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.aof")

	aof, err := NewAOF(path, "no")
	require.NoError(t, err)
	defer aof.Close()

	var replayed [][]string
	err = aof.Replay(func(args []string) error {
		replayed = append(replayed, args)
		return nil
	})
	require.NoError(t, err)
	assert.Empty(t, replayed)
}

func TestAOFWriteFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "appendonly.aof")

	aof, err := NewAOF(path, "no")
	require.NoError(t, err)

	err = aof.Write([]string{"SET", "mykey", "myvalue"})
	require.NoError(t, err)
	require.NoError(t, aof.Close())

	// Read raw file.
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	expected := "*3\r\n$3\r\nSET\r\n$5\r\nmykey\r\n$7\r\nmyvalue\r\n"
	assert.Equal(t, expected, string(data))
}

func TestAOFSyncAlways(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "appendonly.aof")

	aof, err := NewAOF(path, "always")
	require.NoError(t, err)
	defer aof.Close()

	// Write should flush immediately.
	err = aof.Write([]string{"SET", "key", "val"})
	require.NoError(t, err)

	// File should exist and have content even before close.
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}


func TestAOFCorruptedReplay(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "appendonly.aof")

	// Write truncated RESP data (incomplete bulk string).
	err := os.WriteFile(path, []byte("*2\r\n$3\r\nSET\r\n$5\r\nhel"), 0644)
	require.NoError(t, err)

	aof, err := NewAOF(path, "no")
	require.NoError(t, err)
	defer aof.Close()

	var replayed [][]string
	err = aof.Replay(func(args []string) error {
		replayed = append(replayed, args)
		return nil
	})
	assert.Error(t, err)
}


func TestAOFEverysecSync(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "appendonly.aof")

	aof, err := NewAOF(path, "everysec")
	require.NoError(t, err)
	defer aof.Close()

	// Write and let periodic sync run.
	err = aof.Write([]string{"SET", "key", "val"})
	require.NoError(t, err)

	// Wait for the ticker to fire.
	time.Sleep(1200 * time.Millisecond)

	// File should have content.
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(data), "SET")
}
