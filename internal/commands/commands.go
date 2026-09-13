package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/itsmunim/kivo/internal/resp"
	"github.com/itsmunim/kivo/internal/store"
)

// Handler is a command handler function.
type Handler func(e *store.Engine, args []resp.Value) resp.Value

// Registry maps command names (uppercase) to handlers.
type Registry struct {
	handlers map[string]Handler
}

// NewRegistry creates a command registry with all built-in commands.
func NewRegistry() *Registry {
	r := &Registry{handlers: make(map[string]Handler)}
	r.registerAll()
	return r
}

// Get looks up a command handler by name.
func (r *Registry) Get(name string) (Handler, bool) {
	h, ok := r.handlers[strings.ToUpper(name)]
	return h, ok
}

func (r *Registry) registerAll() {
	// Connection
	r.handlers["PING"] = cmdPing
	r.handlers["ECHO"] = cmdEcho
	r.handlers["QUIT"] = cmdQuit
	r.handlers["AUTH"] = cmdAuth
	r.handlers["SELECT"] = cmdSelect // no-op for single DB

	// Strings
	r.handlers["GET"] = cmdGet
	r.handlers["SET"] = cmdSet
	r.handlers["DEL"] = cmdDel
	r.handlers["EXISTS"] = cmdExists
	r.handlers["EXPIRE"] = cmdExpire
	r.handlers["PEXPIRE"] = cmdPExpire
	r.handlers["TTL"] = cmdTTL
	r.handlers["PTTL"] = cmdPTTL
	r.handlers["PERSIST"] = cmdPersist
	r.handlers["MGET"] = cmdMGet
	r.handlers["MSET"] = cmdMSet
	r.handlers["INCR"] = cmdIncr
	r.handlers["DECR"] = cmdDecr
	r.handlers["INCRBY"] = cmdIncrBy
	r.handlers["DECRBY"] = cmdDecrBy
	r.handlers["APPEND"] = cmdAppend
	r.handlers["STRLEN"] = cmdStrLen

	// Lists
	r.handlers["LPUSH"] = cmdLPush
	r.handlers["RPUSH"] = cmdRPush
	r.handlers["LPOP"] = cmdLPop
	r.handlers["RPOP"] = cmdRPop
	r.handlers["LRANGE"] = cmdLRange
	r.handlers["LLEN"] = cmdLLen
	r.handlers["LINDEX"] = cmdLIndex
	r.handlers["LREM"] = cmdLRem
	r.handlers["LTRIM"] = cmdLTrim

	// Sets
	r.handlers["SADD"] = cmdSAdd
	r.handlers["SREM"] = cmdSRem
	r.handlers["SMEMBERS"] = cmdSMembers
	r.handlers["SISMEMBER"] = cmdSIsMember
	r.handlers["SCARD"] = cmdSCard
	r.handlers["SPOP"] = cmdSPop
	r.handlers["SRANDMEMBER"] = cmdSRandMember
	r.handlers["SUNION"] = cmdSUnion
	r.handlers["SINTER"] = cmdSInter
	r.handlers["SDIFF"] = cmdSDiff

	// Hashes
	r.handlers["HSET"] = cmdHSet
	r.handlers["HGET"] = cmdHGet
	r.handlers["HGETALL"] = cmdHGetAll
	r.handlers["HDEL"] = cmdHDel
	r.handlers["HLEN"] = cmdHLen
	r.handlers["HEXISTS"] = cmdHExists
	r.handlers["HKEYS"] = cmdHKeys
	r.handlers["HVALS"] = cmdHVals
	r.handlers["HMGET"] = cmdHMGet

	// Sorted Sets
	r.handlers["ZADD"] = cmdZAdd
	r.handlers["ZREM"] = cmdZRem
	r.handlers["ZRANGE"] = cmdZRange
	r.handlers["ZREVRANGE"] = cmdZRevRange
	r.handlers["ZRANGEBYSCORE"] = cmdZRangeByScore
	r.handlers["ZREVRANGEBYSCORE"] = cmdZRevRangeByScore
	r.handlers["ZCARD"] = cmdZCard
	r.handlers["ZSCORE"] = cmdZScore
	r.handlers["ZINCRBY"] = cmdZIncrBy
	r.handlers["ZCOUNT"] = cmdZCount
	r.handlers["ZREMRANGEBYSCORE"] = cmdZRemRangeByScore
	r.handlers["ZREMRANGEBYRANK"] = cmdZRemRangeByRank

	// Keys
	r.handlers["KEYS"] = cmdKeys
	r.handlers["FLUSHDB"] = cmdFlushDB
	r.handlers["DBSIZE"] = cmdDBSize
	r.handlers["TYPE"] = cmdType
	r.handlers["RENAME"] = cmdRename
	r.handlers["RENAMENX"] = cmdRenameNX
}

// --- Helper functions ---

func wrongArity(cmd string) resp.Value {
	return resp.NewError(fmt.Sprintf("ERR wrong number of arguments for '%s' command", cmd))
}

func wrongType() resp.Value {
	return resp.NewError("WRONGTYPE Operation against a key holding the wrong kind of value")
}

func intValue(n int) resp.Value {
	return resp.NewInteger(int64(n))
}

func bulkStrings(ss []string) []resp.Value {
	vs := make([]resp.Value, len(ss))
	for i, s := range ss {
		vs[i] = resp.NewBulkString(s)
	}
	return vs
}

// --- Connection commands ---

func cmdPing(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) == 0 {
		return resp.NewSimpleString("PONG")
	}
	if len(args) == 1 {
		return resp.NewBulkString(args[0].String())
	}
	return wrongArity("PING")
}

func cmdEcho(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("ECHO")
	}
	return resp.NewBulkString(args[0].String())
}

func cmdQuit(e *store.Engine, args []resp.Value) resp.Value {
	return resp.NewSimpleString("OK")
}

func cmdAuth(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("AUTH")
	}
	// Simple password auth: no password configured means auth is disabled.
	// TODO: wire up actual password checking from config.
	return resp.NewSimpleString("OK")
}

func cmdSelect(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("SELECT")
	}
	// Single DB only: SELECT is a no-op.
	return resp.NewSimpleString("OK")
}

// --- String commands ---

func cmdGet(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("GET")
	}
	val, ok := e.Get(args[0].String())
	if !ok {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString(val)
}

func cmdSet(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return wrongArity("SET")
	}
	key := args[0].String()
	value := args[1].String()
	var ttl time.Duration

	// Parse options: EX seconds, PX milliseconds.
	for i := 2; i < len(args); i++ {
		op := strings.ToUpper(args[i].String())
		switch op {
		case "EX":
			if i+1 >= len(args) {
				return resp.NewError("ERR syntax error")
			}
			sec, err := strconv.Atoi(args[i+1].String())
			if err != nil {
				return resp.NewError("ERR value is not an integer or out of range")
			}
			ttl = time.Duration(sec) * time.Second
			i++
		case "PX":
			if i+1 >= len(args) {
				return resp.NewError("ERR syntax error")
			}
			ms, err := strconv.Atoi(args[i+1].String())
			if err != nil {
				return resp.NewError("ERR value is not an integer or out of range")
			}
			ttl = time.Duration(ms) * time.Millisecond
			i++
		default:
			return resp.NewError("ERR syntax error")
		}
	}

	e.Set(key, value, ttl)
	return resp.NewSimpleString("OK")
}

func cmdDel(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return wrongArity("DEL")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = arg.String()
	}
	return intValue(e.Del(keys...))
}

func cmdExists(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return wrongArity("EXISTS")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = arg.String()
	}
	return intValue(e.Exists(keys...))
}

func cmdExpire(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("EXPIRE")
	}
	sec, err := strconv.Atoi(args[1].String())
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	if e.Expire(args[0].String(), time.Duration(sec)*time.Second) {
		return resp.NewInteger(1)
	}
	return resp.NewInteger(0)
}

func cmdPExpire(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("PEXPIRE")
	}
	ms, err := strconv.Atoi(args[1].String())
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	if e.Expire(args[0].String(), time.Duration(ms)*time.Millisecond) {
		return resp.NewInteger(1)
	}
	return resp.NewInteger(0)
}

func cmdTTL(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("TTL")
	}
	ttl := e.TTL(args[0].String())
	if ttl < 0 {
		return resp.NewInteger(ttl) // -1 or -2.
	}
	return resp.NewInteger(ttl / 1000) // Convert ms to seconds.
}

func cmdPTTL(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("PTTL")
	}
	return resp.NewInteger(e.TTL(args[0].String()))
}

func cmdPersist(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("PERSIST")
	}
	if e.Persist(args[0].String()) {
		return resp.NewInteger(1)
	}
	return resp.NewInteger(0)
}

func cmdMGet(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return wrongArity("MGET")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = arg.String()
	}
	vals := e.MGet(keys...)
	results := make([]resp.Value, len(vals))
	for i, v := range vals {
		if v == "" {
			results[i] = resp.NewNullBulkString()
		} else {
			results[i] = resp.NewBulkString(v)
		}
	}
	return resp.NewArray(results...)
}

func cmdMSet(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 2 || len(args)%2 != 0 {
		return wrongArity("MSET")
	}
	pairs := make([]string, len(args))
	for i, arg := range args {
		pairs[i] = arg.String()
	}
	e.MSet(pairs)
	return resp.NewSimpleString("OK")
}

func cmdIncr(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("INCR")
	}
	val, err := e.IncrBy(args[0].String(), 1)
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	return resp.NewInteger(val)
}

func cmdDecr(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("DECR")
	}
	val, err := e.IncrBy(args[0].String(), -1)
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	return resp.NewInteger(val)
}

func cmdIncrBy(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("INCRBY")
	}
	delta, err := strconv.ParseInt(args[1].String(), 10, 64)
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	val, err := e.IncrBy(args[0].String(), delta)
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	return resp.NewInteger(val)
}

func cmdDecrBy(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("DECRBY")
	}
	delta, err := strconv.ParseInt(args[1].String(), 10, 64)
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	val, err := e.IncrBy(args[0].String(), -delta)
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	return resp.NewInteger(val)
}

func cmdAppend(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("APPEND")
	}
	return intValue(e.Append(args[0].String(), args[1].String()))
}

func cmdStrLen(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("STRLEN")
	}
	return intValue(e.StrLen(args[0].String()))
}

// --- List commands ---

func cmdLPush(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return wrongArity("LPUSH")
	}
	key := args[0].String()
	values := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		values[i-1] = args[i].String()
	}
	return intValue(e.LPush(key, values...))
}

func cmdRPush(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return wrongArity("RPUSH")
	}
	key := args[0].String()
	values := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		values[i-1] = args[i].String()
	}
	return intValue(e.RPush(key, values...))
}

func cmdLPop(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("LPOP")
	}
	val, ok := e.LPop(args[0].String())
	if !ok {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString(val)
}

func cmdRPop(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("RPOP")
	}
	val, ok := e.RPop(args[0].String())
	if !ok {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString(val)
}

func cmdLRange(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 3 {
		return wrongArity("LRANGE")
	}
	start, err1 := strconv.Atoi(args[1].String())
	stop, err2 := strconv.Atoi(args[2].String())
	if err1 != nil || err2 != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	vals := e.LRange(args[0].String(), start, stop)
	return resp.NewArray(bulkStrings(vals)...)
}

func cmdLLen(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("LLEN")
	}
	return intValue(e.LLen(args[0].String()))
}

func cmdLIndex(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("LINDEX")
	}
	index, err := strconv.Atoi(args[1].String())
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	val, ok := e.LIndex(args[0].String(), index)
	if !ok {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString(val)
}

func cmdLRem(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 3 {
		return wrongArity("LREM")
	}
	count, err := strconv.Atoi(args[1].String())
	if err != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	return intValue(e.LRem(args[0].String(), count, args[2].String()))
}

func cmdLTrim(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 3 {
		return wrongArity("LTRIM")
	}
	start, err1 := strconv.Atoi(args[1].String())
	stop, err2 := strconv.Atoi(args[2].String())
	if err1 != nil || err2 != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	e.LTrim(args[0].String(), start, stop)
	return resp.NewSimpleString("OK")
}

// --- Set commands ---

func cmdSAdd(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return wrongArity("SADD")
	}
	key := args[0].String()
	members := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		members[i-1] = args[i].String()
	}
	return intValue(e.SAdd(key, members...))
}

func cmdSRem(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return wrongArity("SREM")
	}
	key := args[0].String()
	members := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		members[i-1] = args[i].String()
	}
	return intValue(e.SRem(key, members...))
}

func cmdSMembers(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("SMEMBERS")
	}
	members := e.SMembers(args[0].String())
	return resp.NewArray(bulkStrings(members)...)
}

func cmdSIsMember(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("SISMEMBER")
	}
	if e.SIsMember(args[0].String(), args[1].String()) {
		return resp.NewInteger(1)
	}
	return resp.NewInteger(0)
}

func cmdSCard(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("SCARD")
	}
	return intValue(e.SCard(args[0].String()))
}

func cmdSPop(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 1 || len(args) > 2 {
		return wrongArity("SPOP")
	}
	count := 1
	if len(args) == 2 {
		var err error
		count, err = strconv.Atoi(args[1].String())
		if err != nil {
			return resp.NewError("ERR value is not an integer or out of range")
		}
	}
	members := e.SPop(args[0].String(), count)
	if members == nil {
		return resp.NewNullArray()
	}
	return resp.NewArray(bulkStrings(members)...)
}

func cmdSRandMember(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 1 || len(args) > 2 {
		return wrongArity("SRANDMEMBER")
	}
	count := 1
	if len(args) == 2 {
		var err error
		count, err = strconv.Atoi(args[1].String())
		if err != nil {
			return resp.NewError("ERR value is not an integer or out of range")
		}
	}
	members := e.SRandMember(args[0].String(), count)
	if members == nil {
		return resp.NewNullArray()
	}
	return resp.NewArray(bulkStrings(members)...)
}

func cmdSUnion(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return wrongArity("SUNION")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = arg.String()
	}
	return resp.NewArray(bulkStrings(e.SUnion(keys...))...)
}

func cmdSInter(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return wrongArity("SINTER")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = arg.String()
	}
	return resp.NewArray(bulkStrings(e.SInter(keys...))...)
}

func cmdSDiff(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 1 {
		return wrongArity("SDIFF")
	}
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = arg.String()
	}
	return resp.NewArray(bulkStrings(e.SDiff(keys...))...)
}

// --- Hash commands ---

func cmdHSet(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 3 || (len(args)-1)%2 != 0 {
		return wrongArity("HSET")
	}
	key := args[0].String()
	pairs := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		pairs[i-1] = args[i].String()
	}
	return intValue(e.HSet(key, pairs))
}

func cmdHGet(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("HGET")
	}
	val, ok := e.HGet(args[0].String(), args[1].String())
	if !ok {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString(val)
}

func cmdHGetAll(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("HGETALL")
	}
	pairs := e.HGetAll(args[0].String())
	return resp.NewArray(bulkStrings(pairs)...)
}

func cmdHDel(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return wrongArity("HDEL")
	}
	key := args[0].String()
	fields := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		fields[i-1] = args[i].String()
	}
	return intValue(e.HDel(key, fields...))
}

func cmdHLen(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("HLEN")
	}
	return intValue(e.HLen(args[0].String()))
}

func cmdHExists(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("HEXISTS")
	}
	if e.HExists(args[0].String(), args[1].String()) {
		return resp.NewInteger(1)
	}
	return resp.NewInteger(0)
}

func cmdHKeys(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("HKEYS")
	}
	return resp.NewArray(bulkStrings(e.HKeys(args[0].String()))...)
}

func cmdHVals(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("HVALS")
	}
	return resp.NewArray(bulkStrings(e.HVals(args[0].String()))...)
}

func cmdHMGet(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return wrongArity("HMGET")
	}
	key := args[0].String()
	fields := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		fields[i-1] = args[i].String()
	}
	vals := e.HMGet(key, fields...)
	results := make([]resp.Value, len(vals))
	for i, v := range vals {
		if v == "" {
			results[i] = resp.NewNullBulkString()
		} else {
			results[i] = resp.NewBulkString(v)
		}
	}
	return resp.NewArray(results...)
}

// --- Sorted Set commands ---

func cmdZAdd(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 3 || (len(args)-1)%2 != 0 {
		return wrongArity("ZADD")
	}
	key := args[0].String()
	members := make(map[string]float64)
	for i := 1; i < len(args); i += 2 {
		score, err := strconv.ParseFloat(args[i].String(), 64)
		if err != nil {
			return resp.NewError("ERR value is not a valid float")
		}
		members[args[i+1].String()] = score
	}
	return intValue(e.ZAdd(key, members))
}

func cmdZRem(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 2 {
		return wrongArity("ZREM")
	}
	key := args[0].String()
	members := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		members[i-1] = args[i].String()
	}
	return intValue(e.ZRem(key, members...))
}

func cmdZRange(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 3 || len(args) > 4 {
		return wrongArity("ZRANGE")
	}
	start, err1 := strconv.Atoi(args[1].String())
	stop, err2 := strconv.Atoi(args[2].String())
	if err1 != nil || err2 != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	withScores := false
	if len(args) == 4 && strings.ToUpper(args[3].String()) == "WITHSCORES" {
		withScores = true
	}
	return resp.NewArray(bulkStrings(e.ZRange(args[0].String(), start, stop, withScores))...)
}

func cmdZRevRange(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) < 3 || len(args) > 4 {
		return wrongArity("ZREVRANGE")
	}
	start, err1 := strconv.Atoi(args[1].String())
	stop, err2 := strconv.Atoi(args[2].String())
	if err1 != nil || err2 != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	withScores := false
	if len(args) == 4 && strings.ToUpper(args[3].String()) == "WITHSCORES" {
		withScores = true
	}
	return resp.NewArray(bulkStrings(e.ZRevRange(args[0].String(), start, stop, withScores))...)
}

func cmdZRangeByScore(e *store.Engine, args []resp.Value) resp.Value {
	return zRangeByScoreGeneric(e, args, false)
}

func cmdZRevRangeByScore(e *store.Engine, args []resp.Value) resp.Value {
	return zRangeByScoreGeneric(e, args, true)
}

func zRangeByScoreGeneric(e *store.Engine, args []resp.Value, reverse bool) resp.Value {
	if len(args) < 3 {
		return wrongArity("ZRANGEBYSCORE")
	}
	key := args[0].String()

	// For ZRANGEBYSCORE: args are min, max.
	// For ZREVRANGEBYSCORE: args are max, min (Redis convention).
	var s1, s2 float64
	var err1, err2 error
	s1, err1 = strconv.ParseFloat(args[1].String(), 64)
	s2, err2 = strconv.ParseFloat(args[2].String(), 64)
	if err1 != nil || err2 != nil {
		return resp.NewError("ERR min or max is not a float")
	}

	withScores := false
	offset := 0
	count := -1
	for i := 3; i < len(args); i++ {
		op := strings.ToUpper(args[i].String())
		switch op {
		case "WITHSCORES":
			withScores = true
		case "LIMIT":
			if i+2 >= len(args) {
				return resp.NewError("ERR syntax error")
			}
			var err error
			offset, err = strconv.Atoi(args[i+1].String())
			if err != nil {
				return resp.NewError("ERR value is not an integer or out of range")
			}
			count, err = strconv.Atoi(args[i+2].String())
			if err != nil {
				return resp.NewError("ERR value is not an integer or out of range")
			}
			i += 2
		default:
			return resp.NewError("ERR syntax error")
		}
	}

	if reverse {
		// s1 is max, s2 is min for ZREVRANGEBYSCORE.
		return resp.NewArray(bulkStrings(e.ZRevRangeByScore(key, s1, s2, withScores, offset, count))...)
	}
	return resp.NewArray(bulkStrings(e.ZRangeByScore(key, s1, s2, withScores, offset, count))...)
}

func cmdZCard(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("ZCARD")
	}
	return intValue(e.ZCard(args[0].String()))
}

func cmdZScore(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("ZSCORE")
	}
	score, ok := e.ZScore(args[0].String(), args[1].String())
	if !ok {
		return resp.NewNullBulkString()
	}
	return resp.NewBulkString(strconv.FormatFloat(score, 'f', -1, 64))
}

func cmdZIncrBy(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 3 {
		return wrongArity("ZINCRBY")
	}
	increment, err := strconv.ParseFloat(args[1].String(), 64)
	if err != nil {
		return resp.NewError("ERR value is not a valid float")
	}
	score := e.ZIncrBy(args[0].String(), increment, args[2].String())
	return resp.NewBulkString(strconv.FormatFloat(score, 'f', -1, 64))
}

func cmdZCount(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 3 {
		return wrongArity("ZCOUNT")
	}
	min, err1 := strconv.ParseFloat(args[1].String(), 64)
	max, err2 := strconv.ParseFloat(args[2].String(), 64)
	if err1 != nil || err2 != nil {
		return resp.NewError("ERR min or max is not a float")
	}
	return intValue(e.ZCount(args[0].String(), min, max))
}

func cmdZRemRangeByScore(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 3 {
		return wrongArity("ZREMRANGEBYSCORE")
	}
	min, err1 := strconv.ParseFloat(args[1].String(), 64)
	max, err2 := strconv.ParseFloat(args[2].String(), 64)
	if err1 != nil || err2 != nil {
		return resp.NewError("ERR min or max is not a float")
	}
	return intValue(e.ZRemRangeByScore(args[0].String(), min, max))
}

func cmdZRemRangeByRank(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 3 {
		return wrongArity("ZREMRANGEBYRANK")
	}
	start, err1 := strconv.Atoi(args[1].String())
	stop, err2 := strconv.Atoi(args[2].String())
	if err1 != nil || err2 != nil {
		return resp.NewError("ERR value is not an integer or out of range")
	}
	return intValue(e.ZRemRangeByRank(args[0].String(), start, stop))
}

// --- Key commands ---

func cmdKeys(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("KEYS")
	}
	return resp.NewArray(bulkStrings(e.Keys(args[0].String()))...)
}

func cmdFlushDB(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 0 {
		return wrongArity("FLUSHDB")
	}
	e.FlushDB()
	return resp.NewSimpleString("OK")
}

func cmdDBSize(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 0 {
		return wrongArity("DBSIZE")
	}
	return intValue(e.DBSize())
}

func cmdType(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 1 {
		return wrongArity("TYPE")
	}
	return resp.NewSimpleString(e.Type(args[0].String()))
}

func cmdRename(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("RENAME")
	}
	if err := e.Rename(args[0].String(), args[1].String()); err != nil {
		return resp.NewError("ERR no such key")
	}
	return resp.NewSimpleString("OK")
}

func cmdRenameNX(e *store.Engine, args []resp.Value) resp.Value {
	if len(args) != 2 {
		return wrongArity("RENAMENX")
	}
	ok, err := e.RenameNX(args[0].String(), args[1].String())
	if err != nil {
		return resp.NewError("ERR no such key")
	}
	if ok {
		return resp.NewInteger(1)
	}
	return resp.NewInteger(0)
}
