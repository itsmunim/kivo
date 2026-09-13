package commands

import (
	"testing"
	"time"

	"github.com/itsmunim/kivo/internal/resp"
	"github.com/itsmunim/kivo/internal/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newEngine(t *testing.T) *store.Engine {
	e := store.NewEngine()
	t.Cleanup(func() { e.Stop() })
	return e
}

func TestPing(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()

	h, _ := r.Get("PING")
	v := h(e, nil)
	assert.Equal(t, "PONG", v.String())

	v = h(e, []resp.Value{resp.NewBulkString("hello")})
	assert.Equal(t, "hello", v.String())
}

func TestGetSet(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()

	set, _ := r.Get("SET")
	get, _ := r.Get("GET")

	v := set(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString("value")})
	assert.Equal(t, "OK", v.String())

	v = get(e, []resp.Value{resp.NewBulkString("key")})
	assert.Equal(t, "value", v.String())

	v = get(e, []resp.Value{resp.NewBulkString("missing")})
	assert.True(t, v.IsNull())
}

func TestSetWithEX(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")

	v := set(e, []resp.Value{
		resp.NewBulkString("key"),
		resp.NewBulkString("val"),
		resp.NewBulkString("EX"),
		resp.NewBulkString("1"),
	})
	assert.Equal(t, "OK", v.String())

	time.Sleep(1100 * time.Millisecond)
	get, _ := r.Get("GET")
	v = get(e, []resp.Value{resp.NewBulkString("key")})
	assert.True(t, v.IsNull())
}

func TestDel(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	del, _ := r.Get("DEL")

	set(e, []resp.Value{resp.NewBulkString("a"), resp.NewBulkString("1")})
	set(e, []resp.Value{resp.NewBulkString("b"), resp.NewBulkString("2")})

	v := del(e, []resp.Value{resp.NewBulkString("a"), resp.NewBulkString("b"), resp.NewBulkString("missing")})
	assert.Equal(t, int64(2), v.Integer())
}

func TestExists(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	exists, _ := r.Get("EXISTS")

	set(e, []resp.Value{resp.NewBulkString("a"), resp.NewBulkString("1")})
	v := exists(e, []resp.Value{resp.NewBulkString("a"), resp.NewBulkString("missing")})
	assert.Equal(t, int64(1), v.Integer())
}

func TestIncrDecr(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	incr, _ := r.Get("INCR")
	decr, _ := r.Get("DECR")
	incrby, _ := r.Get("INCRBY")

	v := incr(e, []resp.Value{resp.NewBulkString("counter")})
	assert.Equal(t, int64(1), v.Integer())

	v = incrby(e, []resp.Value{resp.NewBulkString("counter"), resp.NewBulkString("5")})
	assert.Equal(t, int64(6), v.Integer())

	v = decr(e, []resp.Value{resp.NewBulkString("counter")})
	assert.Equal(t, int64(5), v.Integer())
}

func TestAppendStrLen(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	appendCmd, _ := r.Get("APPEND")
	strlen, _ := r.Get("STRLEN")

	v := appendCmd(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString("hello")})
	assert.Equal(t, int64(5), v.Integer())

	v = appendCmd(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString(" world")})
	assert.Equal(t, int64(11), v.Integer())

	v = strlen(e, []resp.Value{resp.NewBulkString("key")})
	assert.Equal(t, int64(11), v.Integer())
}

func TestMGetMSet(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	mset, _ := r.Get("MSET")
	mget, _ := r.Get("MGET")

	v := mset(e, []resp.Value{
		resp.NewBulkString("a"), resp.NewBulkString("1"),
		resp.NewBulkString("b"), resp.NewBulkString("2"),
	})
	assert.Equal(t, "OK", v.String())

	v = mget(e, []resp.Value{resp.NewBulkString("a"), resp.NewBulkString("missing"), resp.NewBulkString("b")})
	arr := v.Array()
	require.Len(t, arr, 3)
	assert.Equal(t, "1", arr[0].String())
	assert.True(t, arr[1].IsNull())
	assert.Equal(t, "2", arr[2].String())
}

func TestLPushRPush(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	lpush, _ := r.Get("LPUSH")
	rpush, _ := r.Get("RPUSH")
	lrange, _ := r.Get("LRANGE")

	lpush(e, []resp.Value{resp.NewBulkString("list"), resp.NewBulkString("a"), resp.NewBulkString("b")})
	rpush(e, []resp.Value{resp.NewBulkString("list"), resp.NewBulkString("c")})

	v := lrange(e, []resp.Value{resp.NewBulkString("list"), resp.NewBulkString("0"), resp.NewBulkString("-1")})
	arr := v.Array()
	require.Len(t, arr, 3)
	assert.Equal(t, "b", arr[0].String())
	assert.Equal(t, "a", arr[1].String())
	assert.Equal(t, "c", arr[2].String())
}

func TestLPopRPop(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	rpush, _ := r.Get("RPUSH")
	lpop, _ := r.Get("LPOP")
	rpop, _ := r.Get("RPOP")

	rpush(e, []resp.Value{resp.NewBulkString("list"), resp.NewBulkString("a"), resp.NewBulkString("b"), resp.NewBulkString("c")})

	v := lpop(e, []resp.Value{resp.NewBulkString("list")})
	assert.Equal(t, "a", v.String())

	v = rpop(e, []resp.Value{resp.NewBulkString("list")})
	assert.Equal(t, "c", v.String())
}

func TestSAddSMembers(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	sadd, _ := r.Get("SADD")
	smembers, _ := r.Get("SMEMBERS")

	v := sadd(e, []resp.Value{resp.NewBulkString("set"), resp.NewBulkString("a"), resp.NewBulkString("b")})
	assert.Equal(t, int64(2), v.Integer())

	v = smembers(e, []resp.Value{resp.NewBulkString("set")})
	arr := v.Array()
	require.Len(t, arr, 2)
}

func TestHSetHGet(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	hset, _ := r.Get("HSET")
	hget, _ := r.Get("HGET")
	hgetall, _ := r.Get("HGETALL")

	v := hset(e, []resp.Value{
		resp.NewBulkString("hash"),
		resp.NewBulkString("field1"), resp.NewBulkString("val1"),
		resp.NewBulkString("field2"), resp.NewBulkString("val2"),
	})
	assert.Equal(t, int64(2), v.Integer())

	v = hget(e, []resp.Value{resp.NewBulkString("hash"), resp.NewBulkString("field1")})
	assert.Equal(t, "val1", v.String())

	v = hgetall(e, []resp.Value{resp.NewBulkString("hash")})
	arr := v.Array()
	assert.Len(t, arr, 4)
}

func TestZAddZRange(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	zadd, _ := r.Get("ZADD")
	zrange, _ := r.Get("ZRANGE")

	v := zadd(e, []resp.Value{
		resp.NewBulkString("zset"),
		resp.NewBulkString("1"), resp.NewBulkString("a"),
		resp.NewBulkString("2"), resp.NewBulkString("b"),
	})
	assert.Equal(t, int64(2), v.Integer())

	v = zrange(e, []resp.Value{resp.NewBulkString("zset"), resp.NewBulkString("0"), resp.NewBulkString("-1")})
	arr := v.Array()
	require.Len(t, arr, 2)
	assert.Equal(t, "a", arr[0].String())
	assert.Equal(t, "b", arr[1].String())
}

func TestKeys(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	keys, _ := r.Get("KEYS")

	set(e, []resp.Value{resp.NewBulkString("user:1"), resp.NewBulkString("a")})
	set(e, []resp.Value{resp.NewBulkString("user:2"), resp.NewBulkString("b")})

	v := keys(e, []resp.Value{resp.NewBulkString("user:*")})
	arr := v.Array()
	require.Len(t, arr, 2)
}

func TestFlushDB(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	flushdb, _ := r.Get("FLUSHDB")
	dbsize, _ := r.Get("DBSIZE")

	set(e, []resp.Value{resp.NewBulkString("a"), resp.NewBulkString("1")})
	v := flushdb(e, nil)
	assert.Equal(t, "OK", v.String())

	v = dbsize(e, nil)
	assert.Equal(t, int64(0), v.Integer())
}

func TestType(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	lpush, _ := r.Get("LPUSH")
	typ, _ := r.Get("TYPE")

	set(e, []resp.Value{resp.NewBulkString("str"), resp.NewBulkString("val")})
	lpush(e, []resp.Value{resp.NewBulkString("list"), resp.NewBulkString("a")})

	v := typ(e, []resp.Value{resp.NewBulkString("str")})
	assert.Equal(t, "string", v.String())

	v = typ(e, []resp.Value{resp.NewBulkString("list")})
	assert.Equal(t, "list", v.String())

	v = typ(e, []resp.Value{resp.NewBulkString("missing")})
	assert.Equal(t, "none", v.String())
}

func TestRename(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	rename, _ := r.Get("RENAME")
	get, _ := r.Get("GET")

	set(e, []resp.Value{resp.NewBulkString("old"), resp.NewBulkString("val")})
	v := rename(e, []resp.Value{resp.NewBulkString("old"), resp.NewBulkString("new")})
	assert.Equal(t, "OK", v.String())

	v = get(e, []resp.Value{resp.NewBulkString("new")})
	assert.Equal(t, "val", v.String())
}

func TestAuth(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	auth, _ := r.Get("AUTH")

	v := auth(e, []resp.Value{resp.NewBulkString("password")})
	assert.Equal(t, "OK", v.String())
}

func TestSelect(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	selectCmd, _ := r.Get("SELECT")

	v := selectCmd(e, []resp.Value{resp.NewBulkString("5")})
	assert.Equal(t, "OK", v.String())
}

func TestWrongArity(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	get, _ := r.Get("GET")

	v := get(e, []resp.Value{})
	assert.Equal(t, resp.Error, v.Type())
	assert.Contains(t, v.Error(), "wrong number of arguments")
}

// --- Additional command coverage tests ---

func TestExpireAndTTL(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	expire, _ := r.Get("EXPIRE")
	ttl, _ := r.Get("TTL")
	pttl, _ := r.Get("PTTL")
	persist, _ := r.Get("PERSIST")

	set(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString("val")})

	v := expire(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString("3600")})
	assert.Equal(t, int64(1), v.Integer())

	v = ttl(e, []resp.Value{resp.NewBulkString("key")})
	assert.True(t, v.Integer() > 3500 && v.Integer() <= 3600)

	v = pttl(e, []resp.Value{resp.NewBulkString("key")})
	assert.True(t, v.Integer() > 3500000)

	v = persist(e, []resp.Value{resp.NewBulkString("key")})
	assert.Equal(t, int64(1), v.Integer())

	v = ttl(e, []resp.Value{resp.NewBulkString("key")})
	assert.Equal(t, int64(-1), v.Integer())

	v = expire(e, []resp.Value{resp.NewBulkString("missing"), resp.NewBulkString("1")})
	assert.Equal(t, int64(0), v.Integer())
}

func TestPExpire(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	pexpire, _ := r.Get("PEXPIRE")
	get, _ := r.Get("GET")

	set(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString("val")})
	v := pexpire(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString("100")})
	assert.Equal(t, int64(1), v.Integer())

	time.Sleep(150 * time.Millisecond)
	v = get(e, []resp.Value{resp.NewBulkString("key")})
	assert.True(t, v.IsNull())
}

func TestDecrBy(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	decrby, _ := r.Get("DECRBY")

	v := decrby(e, []resp.Value{resp.NewBulkString("counter"), resp.NewBulkString("5")})
	assert.Equal(t, int64(-5), v.Integer())
}

func TestLLenLIndexLRemLTrim(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	rpush, _ := r.Get("RPUSH")
	llen, _ := r.Get("LLEN")
	lindex, _ := r.Get("LINDEX")
	lrem, _ := r.Get("LREM")
	ltrim, _ := r.Get("LTRIM")

	rpush(e, []resp.Value{resp.NewBulkString("list"), resp.NewBulkString("a"), resp.NewBulkString("b"), resp.NewBulkString("c")})

	v := llen(e, []resp.Value{resp.NewBulkString("list")})
	assert.Equal(t, int64(3), v.Integer())

	v = lindex(e, []resp.Value{resp.NewBulkString("list"), resp.NewBulkString("1")})
	assert.Equal(t, "b", v.String())

	v = lrem(e, []resp.Value{resp.NewBulkString("list"), resp.NewBulkString("0"), resp.NewBulkString("b")})
	assert.Equal(t, int64(1), v.Integer())

	v = ltrim(e, []resp.Value{resp.NewBulkString("list"), resp.NewBulkString("0"), resp.NewBulkString("0")})
	assert.Equal(t, "OK", v.String())

	v = llen(e, []resp.Value{resp.NewBulkString("list")})
	assert.Equal(t, int64(1), v.Integer())
}

func TestSetCardPopRandMember(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	sadd, _ := r.Get("SADD")
	scard, _ := r.Get("SCARD")
	spop, _ := r.Get("SPOP")
	srandmember, _ := r.Get("SRANDMEMBER")
	srem, _ := r.Get("SREM")
	sismember, _ := r.Get("SISMEMBER")

	sadd(e, []resp.Value{resp.NewBulkString("set"), resp.NewBulkString("a"), resp.NewBulkString("b"), resp.NewBulkString("c")})

	v := scard(e, []resp.Value{resp.NewBulkString("set")})
	assert.Equal(t, int64(3), v.Integer())

	v = sismember(e, []resp.Value{resp.NewBulkString("set"), resp.NewBulkString("a")})
	assert.Equal(t, int64(1), v.Integer())

	v = spop(e, []resp.Value{resp.NewBulkString("set"), resp.NewBulkString("1")})
	arr := v.Array()
	require.Len(t, arr, 1)

	// SPOP removed one random member; we don't know which.
	v = scard(e, []resp.Value{resp.NewBulkString("set")})
	assert.Equal(t, int64(2), v.Integer())

	v = srandmember(e, []resp.Value{resp.NewBulkString("set"), resp.NewBulkString("1")})
	arr = v.Array()
	require.Len(t, arr, 1)

	// SRANDMEMBER does not remove.
	v = scard(e, []resp.Value{resp.NewBulkString("set")})
	assert.Equal(t, int64(2), v.Integer())

	// SREM may return 0 or 1 depending on whether SPOP removed 'a'.
	v = srem(e, []resp.Value{resp.NewBulkString("set"), resp.NewBulkString("a")})
	assert.True(t, v.Integer() == 0 || v.Integer() == 1)
}

func TestSetUnionInterDiff(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	sadd, _ := r.Get("SADD")
	sunion, _ := r.Get("SUNION")
	sinter, _ := r.Get("SINTER")
	sdiff, _ := r.Get("SDIFF")

	sadd(e, []resp.Value{resp.NewBulkString("s1"), resp.NewBulkString("a"), resp.NewBulkString("b")})
	sadd(e, []resp.Value{resp.NewBulkString("s2"), resp.NewBulkString("b"), resp.NewBulkString("c")})

	v := sunion(e, []resp.Value{resp.NewBulkString("s1"), resp.NewBulkString("s2")})
	arr := v.Array()
	assert.Len(t, arr, 3)

	v = sinter(e, []resp.Value{resp.NewBulkString("s1"), resp.NewBulkString("s2")})
	arr = v.Array()
	assert.Len(t, arr, 1)
	assert.Equal(t, "b", arr[0].String())

	v = sdiff(e, []resp.Value{resp.NewBulkString("s1"), resp.NewBulkString("s2")})
	arr = v.Array()
	assert.Len(t, arr, 1)
	assert.Equal(t, "a", arr[0].String())
}

func TestHashCommands(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	hset, _ := r.Get("HSET")
	hdel, _ := r.Get("HDEL")
	hlen, _ := r.Get("HLEN")
	hexists, _ := r.Get("HEXISTS")
	hkeys, _ := r.Get("HKEYS")
	hvals, _ := r.Get("HVALS")
	hmget, _ := r.Get("HMGET")

	hset(e, []resp.Value{
		resp.NewBulkString("hash"),
		resp.NewBulkString("f1"), resp.NewBulkString("v1"),
		resp.NewBulkString("f2"), resp.NewBulkString("v2"),
	})

	v := hlen(e, []resp.Value{resp.NewBulkString("hash")})
	assert.Equal(t, int64(2), v.Integer())

	v = hexists(e, []resp.Value{resp.NewBulkString("hash"), resp.NewBulkString("f1")})
	assert.Equal(t, int64(1), v.Integer())

	v = hkeys(e, []resp.Value{resp.NewBulkString("hash")})
	arr := v.Array()
	assert.Len(t, arr, 2)

	v = hvals(e, []resp.Value{resp.NewBulkString("hash")})
	arr = v.Array()
	assert.Len(t, arr, 2)

	v = hmget(e, []resp.Value{resp.NewBulkString("hash"), resp.NewBulkString("f1"), resp.NewBulkString("missing")})
	arr = v.Array()
	require.Len(t, arr, 2)
	assert.Equal(t, "v1", arr[0].String())
	assert.True(t, arr[1].IsNull())

	v = hdel(e, []resp.Value{resp.NewBulkString("hash"), resp.NewBulkString("f1")})
	assert.Equal(t, int64(1), v.Integer())

	v = hlen(e, []resp.Value{resp.NewBulkString("hash")})
	assert.Equal(t, int64(1), v.Integer())
}

func TestZSetCommands(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	zadd, _ := r.Get("ZADD")
	zrange, _ := r.Get("ZRANGE")
	zrevrange, _ := r.Get("ZREVRANGE")
	zcard, _ := r.Get("ZCARD")
	zscore, _ := r.Get("ZSCORE")
	zincrby, _ := r.Get("ZINCRBY")
	zcount, _ := r.Get("ZCOUNT")
	zremrangebyscore, _ := r.Get("ZREMRANGEBYSCORE")
	zremrangebyrank, _ := r.Get("ZREMRANGEBYRANK")
	zrangebyscore, _ := r.Get("ZRANGEBYSCORE")
	zrevrangebyscore, _ := r.Get("ZREVRANGEBYSCORE")

	zadd(e, []resp.Value{
		resp.NewBulkString("zset"),
		resp.NewBulkString("1"), resp.NewBulkString("a"),
		resp.NewBulkString("2"), resp.NewBulkString("b"),
		resp.NewBulkString("3"), resp.NewBulkString("c"),
	})

	v := zcard(e, []resp.Value{resp.NewBulkString("zset")})
	assert.Equal(t, int64(3), v.Integer())

	v = zscore(e, []resp.Value{resp.NewBulkString("zset"), resp.NewBulkString("b")})
	assert.Equal(t, "2", v.String())

	v = zincrby(e, []resp.Value{resp.NewBulkString("zset"), resp.NewBulkString("5"), resp.NewBulkString("b")})
	assert.Equal(t, "7", v.String())

	// After ZINCRBY, b's score is 7. Only a(1) and c(3) are in range 1-3.
	v = zcount(e, []resp.Value{resp.NewBulkString("zset"), resp.NewBulkString("1"), resp.NewBulkString("3")})
	assert.Equal(t, int64(2), v.Integer())

	v = zrangebyscore(e, []resp.Value{
		resp.NewBulkString("zset"),
		resp.NewBulkString("1"), resp.NewBulkString("3"),
		resp.NewBulkString("WITHSCORES"),
	})
	arr := v.Array()
	assert.Len(t, arr, 4) // 2 members * 2

	v = zrevrangebyscore(e, []resp.Value{
		resp.NewBulkString("zset"),
		resp.NewBulkString("3"), resp.NewBulkString("1"),
	})
	arr = v.Array()
	assert.Len(t, arr, 2)
	assert.Equal(t, "c", arr[0].String())

	// Reverse by rank: b(7), c(3), a(1)
	v = zrevrange(e, []resp.Value{
		resp.NewBulkString("zset"), resp.NewBulkString("0"), resp.NewBulkString("-1"),
	})
	arr = v.Array()
	assert.Len(t, arr, 3)
	assert.Equal(t, "b", arr[0].String())

	// ZREMRANGEBYSCORE 1-2: only a=1 is in range.
	v = zremrangebyscore(e, []resp.Value{
		resp.NewBulkString("zset"), resp.NewBulkString("1"), resp.NewBulkString("2"),
	})
	assert.Equal(t, int64(1), v.Integer())

	v = zcard(e, []resp.Value{resp.NewBulkString("zset")})
	assert.Equal(t, int64(2), v.Integer())

	// Add more for rank test.
	zadd(e, []resp.Value{
		resp.NewBulkString("zset2"),
		resp.NewBulkString("10"), resp.NewBulkString("x"),
		resp.NewBulkString("20"), resp.NewBulkString("y"),
		resp.NewBulkString("30"), resp.NewBulkString("z"),
	})
	v = zremrangebyrank(e, []resp.Value{
		resp.NewBulkString("zset2"), resp.NewBulkString("0"), resp.NewBulkString("1"),
	})
	assert.Equal(t, int64(2), v.Integer())

	v = zrange(e, []resp.Value{resp.NewBulkString("zset2"), resp.NewBulkString("0"), resp.NewBulkString("-1")})
	arr = v.Array()
	assert.Len(t, arr, 1)
	assert.Equal(t, "z", arr[0].String())
}

func TestRenameNX(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	renamenx, _ := r.Get("RENAMENX")
	rename, _ := r.Get("RENAME")

	set(e, []resp.Value{resp.NewBulkString("old"), resp.NewBulkString("val")})
	set(e, []resp.Value{resp.NewBulkString("existing"), resp.NewBulkString("other")})

	v := renamenx(e, []resp.Value{resp.NewBulkString("old"), resp.NewBulkString("existing")})
	assert.Equal(t, int64(0), v.Integer())

	v = renamenx(e, []resp.Value{resp.NewBulkString("old"), resp.NewBulkString("new")})
	assert.Equal(t, int64(1), v.Integer())

	v = rename(e, []resp.Value{resp.NewBulkString("new"), resp.NewBulkString("existing")})
	assert.Equal(t, "OK", v.String())
}

func TestSetWithPX(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	get, _ := r.Get("GET")

	v := set(e, []resp.Value{
		resp.NewBulkString("key"), resp.NewBulkString("val"),
		resp.NewBulkString("PX"), resp.NewBulkString("50"),
	})
	assert.Equal(t, "OK", v.String())

	time.Sleep(100 * time.Millisecond)
	v = get(e, []resp.Value{resp.NewBulkString("key")})
	assert.True(t, v.IsNull())
}

func TestAuthWrongArity(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	auth, _ := r.Get("AUTH")

	v := auth(e, []resp.Value{})
	assert.Equal(t, resp.Error, v.Type())
}

func TestAppendCommand(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	appendCmd, _ := r.Get("APPEND")
	get, _ := r.Get("GET")

	v := appendCmd(e, []resp.Value{
		resp.NewBulkString("key"),
		resp.NewBulkString("hello"),
	})
	assert.Equal(t, int64(5), v.Integer())

	v = appendCmd(e, []resp.Value{
		resp.NewBulkString("key"),
		resp.NewBulkString(" world"),
	})
	assert.Equal(t, int64(11), v.Integer())

	v = get(e, []resp.Value{resp.NewBulkString("key")})
	assert.Equal(t, "hello world", v.String())
}

func TestStrLenCommand(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	strlen, _ := r.Get("STRLEN")

	set(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString("hello")})

	v := strlen(e, []resp.Value{resp.NewBulkString("key")})
	assert.Equal(t, int64(5), v.Integer())

	v = strlen(e, []resp.Value{resp.NewBulkString("missing")})
	assert.Equal(t, int64(0), v.Integer())
}

func TestDecrCommand(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	decr, _ := r.Get("DECR")

	v := decr(e, []resp.Value{resp.NewBulkString("counter")})
	assert.Equal(t, int64(-1), v.Integer())
}

func TestDecrByCommand(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	decrby, _ := r.Get("DECRBY")

	v := decrby(e, []resp.Value{
		resp.NewBulkString("counter"),
		resp.NewBulkString("5"),
	})
	assert.Equal(t, int64(-5), v.Integer())
}

func TestIncrByCommand(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	incrby, _ := r.Get("INCRBY")

	v := incrby(e, []resp.Value{
		resp.NewBulkString("counter"),
		resp.NewBulkString("10"),
	})
	assert.Equal(t, int64(10), v.Integer())
}

func TestExistsCommand(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	exists, _ := r.Get("EXISTS")

	set(e, []resp.Value{resp.NewBulkString("a"), resp.NewBulkString("1")})

	v := exists(e, []resp.Value{
		resp.NewBulkString("a"),
		resp.NewBulkString("b"),
	})
	assert.Equal(t, int64(1), v.Integer())
}

func TestPersistCommand(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	expire, _ := r.Get("EXPIRE")
	ttl, _ := r.Get("TTL")
	persist, _ := r.Get("PERSIST")

	set(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString("val")})
	expire(e, []resp.Value{resp.NewBulkString("key"), resp.NewBulkString("10")})

	v := persist(e, []resp.Value{resp.NewBulkString("key")})
	assert.Equal(t, int64(1), v.Integer())

	v = ttl(e, []resp.Value{resp.NewBulkString("key")})
	assert.Equal(t, int64(-1), v.Integer())
}

func TestPersistMissingKey(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	persist, _ := r.Get("PERSIST")

	v := persist(e, []resp.Value{resp.NewBulkString("missing")})
	assert.Equal(t, int64(0), v.Integer())
}

func TestRenameCommand(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	rename, _ := r.Get("RENAME")
	get, _ := r.Get("GET")

	set(e, []resp.Value{resp.NewBulkString("old"), resp.NewBulkString("val")})

	v := rename(e, []resp.Value{
		resp.NewBulkString("old"),
		resp.NewBulkString("new"),
	})
	assert.Equal(t, "OK", v.String())

	v = get(e, []resp.Value{resp.NewBulkString("new")})
	assert.Equal(t, "val", v.String())
}

func TestRenameNXCommand(t *testing.T) {
	e := newEngine(t)
	r := NewRegistry()
	set, _ := r.Get("SET")
	renamenx, _ := r.Get("RENAMENX")

	set(e, []resp.Value{resp.NewBulkString("a"), resp.NewBulkString("1")})
	set(e, []resp.Value{resp.NewBulkString("b"), resp.NewBulkString("2")})

	v := renamenx(e, []resp.Value{
		resp.NewBulkString("a"),
		resp.NewBulkString("b"),
	})
	assert.Equal(t, int64(0), v.Integer())

	v = renamenx(e, []resp.Value{
		resp.NewBulkString("a"),
		resp.NewBulkString("c"),
	})
	assert.Equal(t, int64(1), v.Integer())
}
