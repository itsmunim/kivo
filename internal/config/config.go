package config

// Config holds server configuration.
type Config struct {
	Addr       string // TCP listen address, e.g. ":6379"
	AOFEnabled bool   // Enable AOF persistence
	AOFPath    string // Path to AOF file
	AOFSync    string // "always", "everysec", or "no"
	MaxMemory  int64  // Max memory in bytes (0 = unlimited)
	Verbose    bool   // Enable debug logging
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		Addr:       ":6379",
		AOFEnabled: true,
		AOFPath:    "appendonly.aof",
		AOFSync:    "everysec",
		MaxMemory:  0, // unlimited
	}
}
