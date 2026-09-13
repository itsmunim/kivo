package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	assert.Equal(t, ":6379", cfg.Addr)
	assert.True(t, cfg.AOFEnabled)
	assert.Equal(t, "appendonly.aof", cfg.AOFPath)
	assert.Equal(t, "everysec", cfg.AOFSync)
	assert.Equal(t, int64(0), cfg.MaxMemory)
	assert.False(t, cfg.Verbose)
}
