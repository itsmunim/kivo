package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/itsmunim/kivo/internal/commands"
	"github.com/itsmunim/kivo/internal/config"
	"github.com/itsmunim/kivo/internal/persistence"
	"github.com/itsmunim/kivo/internal/resp"
	"github.com/itsmunim/kivo/internal/server"
	"github.com/itsmunim/kivo/internal/store"
	"github.com/itsmunim/kivo/internal/webui"
)

func main() {
	fmt.Println("kivo — Redis-compatible in-memory data store")

	cfg := config.Default()

	// CLI flags override defaults.
	flag.StringVar(&cfg.Addr, "addr", cfg.Addr, "TCP listen address (e.g. :6379)")
	flag.Int64Var(&cfg.MaxMemory, "maxmemory", cfg.MaxMemory, "Max memory in bytes (0 = unlimited)")
	flag.BoolVar(&cfg.AOFEnabled, "aof", cfg.AOFEnabled, "Enable AOF persistence")
	flag.StringVar(&cfg.AOFPath, "aof-path", cfg.AOFPath, "Path to AOF file")
	flag.StringVar(&cfg.AOFSync, "aof-sync", cfg.AOFSync, "AOF sync strategy: always, everysec, no")
	flag.BoolVar(&cfg.Verbose, "verbose", cfg.Verbose, "Enable verbose logging")
	flag.BoolVar(&cfg.WebUIEnabled, "webui", cfg.WebUIEnabled, "Enable the web console")
	flag.StringVar(&cfg.WebUIAddr, "webui-addr", cfg.WebUIAddr, "Web console HTTP listen address")
	flag.Parse()

	// Initialize storage engine.
	engine := store.NewEngineWithMaxMemory(cfg.MaxMemory)
	defer engine.Stop()

	// Log warning if maxmemory exceeds available system memory (same as Redis).
	if cfg.MaxMemory > 0 {
		if avail, err := store.AvailableMemory(); err == nil && avail > 0 && cfg.MaxMemory > int64(avail) {
			fmt.Printf("WARNING: maxmemory (%d bytes) exceeds available system memory (%d bytes)\n", cfg.MaxMemory, avail)
		}
	}

	// Initialize command registry.
	registry := commands.NewRegistry()

	// Web console server (embedded React app) on its own HTTP port.
	var web *webui.Server
	if cfg.WebUIEnabled {
		web = webui.New(cfg.WebUIAddr, engine, registry)
		if err := web.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "webui error: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize AOF persistence.
	var aof *persistence.AOF
	if cfg.AOFEnabled {
		var err error
		aof, err = persistence.NewAOF(cfg.AOFPath, cfg.AOFSync)
		if err != nil {
			fmt.Fprintf(os.Stderr, "aof error: %v\n", err)
			os.Exit(1)
		}
		defer aof.Close()

		// Replay AOF on startup.
		aofAbs, err := filepath.Abs(cfg.AOFPath)
		if err != nil {
			aofAbs = cfg.AOFPath
		}
		fmt.Printf("AOF: %s (sync=%s)\n", aofAbs, cfg.AOFSync)
		fmt.Printf("replaying AOF from %s...\n", aofAbs)

		// Replay AOF on startup.
		count := 0
		err = aof.Replay(func(args []string) error {
			handler, ok := registry.Get(args[0])
			if !ok {
				return fmt.Errorf("unknown command: %s", args[0])
			}
			cmdArgs := make([]resp.Value, len(args)-1)
			for i := 1; i < len(args); i++ {
				cmdArgs[i-1] = resp.NewBulkString(args[i])
			}
			handler(engine, cmdArgs)
			count++
			return nil
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "aof replay error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("replayed %d commands from %s\n", count, aofAbs)
	}

	// Create server.
	srv := server.New(cfg, engine, registry, aof)

	// Handle graceful shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine so main can block on shutdown signal.
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.Start()
	}()

	// Block until signal.
	<-sigCh
	fmt.Println("\nshutting down gracefully...")
	// Stop the web console (drains in-flight API calls).
	if web != nil {
		_ = web.Stop()
	}

	// Stop server and wait for all connection handlers to finish.
	// This MUST happen before deferred aof.Close() runs.
	srv.Stop()

	// Wait for server goroutine to finish.
	if err := <-serverErr; err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
