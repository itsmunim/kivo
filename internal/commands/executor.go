package commands

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/itsmunim/kivo/internal/resp"
	"github.com/itsmunim/kivo/internal/store"
)

// Appender appends a write command to durable storage.
// It is implemented by *persistence.AOF.
type Appender interface {
	Write(args []string) error
}

// Executor runs Redis commands against the engine and, for write commands,
// appends them to the Appender afterwards. Both the RESP TCP server and the
// web console share one Executor, so every path that mutates data is
// persisted exactly once — commands can no longer bypass the AOF.
type Executor struct {
	engine   *store.Engine
	registry *Registry
	appender Appender
}

// NewExecutor creates an Executor. appender may be nil to disable persistence.
// A typed-nil pointer (e.g. (*persistence.AOF)(nil) from a nil argument)
// is normalized to a nil Appender so Write is never called on a nil receiver.
func NewExecutor(engine *store.Engine, registry *Registry, appender Appender) *Executor {
	return &Executor{engine: engine, registry: registry, appender: normalizeAppender(appender)}
}

func normalizeAppender(a Appender) Appender {
	if a == nil {
		return nil
	}
	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}
	return a
}

// Execute runs a command from already-parsed RESP arguments; the first
// element is the command name (case-insensitive). It returns the RESP
// result of the command and, for write commands, appends the original
// arguments to the Appender. This is the hot path for the TCP server:
// no intermediate []string conversion exists.
func (x *Executor) Execute(args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewError("ERR empty command")
	}

	name := strings.ToUpper(args[0].String())
	handler, ok := x.registry.Get(name)
	if !ok {
		return resp.NewError(fmt.Sprintf("ERR unknown command '%s'", args[0].String()))
	}

	result := handler(x.engine, args[1:])

	if x.appender != nil && isWriteCommand(name) {
		if err := x.appender.Write(stringArgs(args)); err != nil {
			fmt.Printf("aof write error: %v\n", err)
		}
	}
	return result
}

// ExecuteStrings runs a command given as string arguments (used by the web
// console). It converts to RESP values and delegates to Execute.
func (x *Executor) ExecuteStrings(args []string) resp.Value {
	if len(args) == 0 {
		return resp.NewError("ERR empty command")
	}
	values := make([]resp.Value, len(args))
	for i, s := range args {
		values[i] = resp.NewBulkString(s)
	}
	return x.Execute(values)
}

// stringArgs converts parsed RESP values back to strings for AOF persistence,
// which stores commands as RESP arrays over plain string args.
func stringArgs(args []resp.Value) []string {
	out := make([]string, len(args))
	for i, a := range args {
		out[i] = a.String()
	}
	return out
}

// isWriteCommand returns true if the (already upper-cased) command modifies data.
func isWriteCommand(cmd string) bool {
	switch cmd {
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
