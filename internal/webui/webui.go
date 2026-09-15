// Package webui serves an embedded React single-page application that gives
// kivo a browser-based console: a key browser plus a live Redis terminal.
//
// The React app is built with Vite (see /webui) and its output is embedded
// into the binary at compile time via go:embed, so a single kivo binary
// serves both the RESP TCP port and the web console on its own HTTP port.
package webui

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/itsmunim/kivo/internal/commands"
	"github.com/itsmunim/kivo/internal/resp"
	"github.com/itsmunim/kivo/internal/store"
)

//go:embed dist
var distFS embed.FS

// Server is the HTTP server for the web console.
type Server struct {
	addr     string
	engine   *store.Engine
	executor *commands.Executor
	httpSrv  *http.Server
	mux      *http.ServeMux
}

// New creates a web console server bound to a TCP address (e.g. ":3001").
func New(addr string, engine *store.Engine, executor *commands.Executor) *Server {
	s := &Server{addr: addr, engine: engine, executor: executor, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/keys", s.handleKeys)
	s.mux.HandleFunc("/api/command", s.handleCommand)
	s.mux.HandleFunc("/api/info", s.handleInfo)
	s.mux.Handle("/", s.staticHandler())
}

// Handler exposes the HTTP handler (useful for tests).
func (s *Server) Handler() http.Handler { return s.mux }

// Start binds the listener and serves the console in the background.
// It fails fast (synchronously) if the address is already in use.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("webui: listen %s: %w", s.addr, err)
	}
	s.httpSrv = &http.Server{Handler: s.mux}
	display := s.addr
	if strings.HasPrefix(display, ":") {
		display = "localhost" + display
	}
	fmt.Printf("kivo Web UI: http://%s\n", display)
	go func() { _ = s.httpSrv.Serve(ln) }()
	return nil
}

// Stop gracefully shuts down the web server.
func (s *Server) Stop() error {
	if s.httpSrv == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.httpSrv.Shutdown(ctx)
}

// ---------- API handlers ----------

type keyInfo struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	TTLMs int64  `json:"ttl_ms"`
}

type keysResponse struct {
	Count int       `json:"count"`
	Keys  []keyInfo `json:"keys"`
}

func (s *Server) handleKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	keys := s.engine.Keys("*")
	out := make([]keyInfo, 0, len(keys))
	for _, k := range keys {
		out = append(out, keyInfo{Key: k, Type: s.engine.Type(k), TTLMs: s.engine.TTL(k)})
	}
	writeJSON(w, keysResponse{Count: len(out), Keys: out})
}

type infoResponse struct {
	Server string `json:"server"`
	DBSize int    `json:"db_size"`
}

func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, infoResponse{Server: "kivo", DBSize: s.engine.DBSize()})
}

type commandRequest struct {
	Command string `json:"command"`
}

type commandResponse struct {
	Command    string         `json:"command"`
	DurationMs float64        `json:"duration_ms"`
	Result     map[string]any `json:"result"`
}

func (s *Server) handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req commandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body: "+err.Error(), http.StatusBadRequest)
		return
	}

	args := splitCommand(req.Command)
	if len(args) == 0 {
		writeJSON(w, commandResponse{Command: req.Command, Result: errorResult("ERR empty command")})
		return
	}

	start := time.Now()
	result := s.executor.Execute(args)
	durationMs := float64(time.Since(start).Microseconds()) / 1000.0

	writeJSON(w, commandResponse{
		Command:    req.Command,
		DurationMs: durationMs,
		Result:     valueToJSON(result),
	})
}

// ---------- helpers ----------

func errorResult(msg string) map[string]any {
	return map[string]any{"kind": "error", "value": msg}
}

// valueToJSON converts a RESP value to a JSON object the web terminal can render.
func valueToJSON(v resp.Value) map[string]any {
	switch v.Type() {
	case resp.SimpleString:
		return map[string]any{"kind": "simple", "value": v.String()}
	case resp.Error:
		return map[string]any{"kind": "error", "value": v.Error()}
	case resp.Integer:
		return map[string]any{"kind": "integer", "value": v.Integer()}
	case resp.BulkString:
		if v.IsNull() {
			return map[string]any{"kind": "bulk", "value": nil}
		}
		return map[string]any{"kind": "bulk", "value": v.String()}
	case resp.Array:
		if v.IsNull() {
			return map[string]any{"kind": "array", "value": nil}
		}
		items := make([]map[string]any, 0, len(v.Array()))
		for _, el := range v.Array() {
			items = append(items, valueToJSON(el))
		}
		return map[string]any{"kind": "array", "value": items}
	default:
		return map[string]any{"kind": "unknown", "value": v.String()}
	}
}

// splitCommand splits a raw command line into tokens, honoring double quotes.
// Example: `SET foo "hello world"` -> ["SET", "foo", "hello world"]
func splitCommand(s string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false
	has := false
	flush := func() {
		if has {
			tokens = append(tokens, cur.String())
			cur.Reset()
			has = false
		}
	}
	for _, r := range s {
		switch {
		case r == '"':
			inQuote = !inQuote
			has = true
		case (r == ' ' || r == '\t' || r == '\n' || r == '\r') && !inQuote:
			flush()
		default:
			cur.WriteRune(r)
			has = true
		}
	}
	flush()
	return tokens
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// staticHandler serves the embedded SPA with a fallback to index.html so
// deep links work without server-side routing.
func (s *Server) staticHandler() http.Handler {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		panic(fmt.Sprintf("webui: embedded dist missing: %v", err))
	}
	fileServer := http.FileServer(http.FS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/")
		if p == "" {
			p = "index.html"
		}
		f, err := sub.Open(p)
		if err != nil {
			// Unknown path -> serve the SPA entry point.
			r.URL.Path = "/"
		} else {
			f.Close()
		}
		fileServer.ServeHTTP(w, r)
	})
}
