package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- String tests ---

func TestSetGet(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("key1", "value1", 0)
	val, ok := e.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)
}

func TestGetMissing(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	_, ok := e.Get("missing")
	assert.False(t, ok)
}

func TestSetWithTTL(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("key1", "value1", 50*time.Millisecond)
	val, ok := e.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)

	// Wait for expiration.
	time.Sleep(100 * time.Millisecond)
	_, ok = e.Get("key1")
	assert.False(t, ok)
}

func TestDel(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("a", "1", 0)
	e.Set("b", "2", 0)
	e.Set("c", "3", 0)

	deleted := e.Del("a", "b", "missing")
	assert.Equal(t, 2, deleted)

	_, ok := e.Get("a")
	assert.False(t, ok)
	_, ok = e.Get("b")
	assert.False(t, ok)
	_, ok = e.Get("c")
	assert.True(t, ok)
}

func TestExists(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("a", "1", 0)
	assert.Equal(t, 1, e.Exists("a"))
	assert.Equal(t, 0, e.Exists("missing"))
}

func TestExpire(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("key1", "value1", 0)
	ok := e.Expire("key1", 50*time.Millisecond)
	require.True(t, ok)

	time.Sleep(100 * time.Millisecond)
	_, exists := e.Get("key1")
	assert.False(t, exists)
}

func TestExpireMissing(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	ok := e.Expire("missing", time.Second)
	assert.False(t, ok)
}

func TestTTL(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("key1", "value1", 500*time.Millisecond)
	ttl := e.TTL("key1")
	assert.True(t, ttl > 400 && ttl <= 500, "TTL should be around 500ms, got %d", ttl)

	// Key with no expiry.
	e.Set("noexpiry", "val", 0)
	assert.Equal(t, int64(-1), e.TTL("noexpiry"))

	// Missing key.
	assert.Equal(t, int64(-2), e.TTL("missing"))
}

func TestPersist(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("key1", "value1", 50*time.Millisecond)
	ok := e.Persist("key1")
	require.True(t, ok)

	time.Sleep(100 * time.Millisecond)
	val, exists := e.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", val)
}

func TestAppend(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	len1 := e.Append("key", "hello")
	assert.Equal(t, 5, len1)

	len2 := e.Append("key", " world")
	assert.Equal(t, 11, len2)

	val, ok := e.Get("key")
	require.True(t, ok)
	assert.Equal(t, "hello world", val)
}

func TestStrLen(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("key", "hello", 0)
	assert.Equal(t, 5, e.StrLen("key"))
	assert.Equal(t, 0, e.StrLen("missing"))
}

func TestIncrBy(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	// Create new key.
	val, err := e.IncrBy("counter", 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), val)

	// Increment existing.
	val, err = e.IncrBy("counter", 5)
	require.NoError(t, err)
	assert.Equal(t, int64(6), val)

	// Decrement.
	val, err = e.IncrBy("counter", -3)
	require.NoError(t, err)
	assert.Equal(t, int64(3), val)
}

func TestIncrByNonInteger(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("key", "notanumber", 0)
	_, err := e.IncrBy("key", 1)
	require.Error(t, err)
}

func TestMGetMSet(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.MSet([]string{"a", "1", "b", "2", "c", "3"})
	vals := e.MGet("a", "b", "missing", "c")
	assert.Equal(t, []string{"1", "2", "", "3"}, vals)
}

// --- List tests ---

func TestLPushRPush(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	len1 := e.LPush("list", "a", "b", "c")
	assert.Equal(t, 3, len1)

	// LPush prepends in reverse order: c, b, a.
	vals := e.LRange("list", 0, -1)
	assert.Equal(t, []string{"c", "b", "a"}, vals)

	len2 := e.RPush("list", "x", "y")
	assert.Equal(t, 5, len2)

	vals = e.LRange("list", 0, -1)
	assert.Equal(t, []string{"c", "b", "a", "x", "y"}, vals)
}

func TestLPopRPop(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.RPush("list", "a", "b", "c")

	val, ok := e.LPop("list")
	require.True(t, ok)
	assert.Equal(t, "a", val)

	val, ok = e.RPop("list")
	require.True(t, ok)
	assert.Equal(t, "c", val)

	assert.Equal(t, 1, e.LLen("list"))
}

func TestLPopEmpty(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	_, ok := e.LPop("missing")
	assert.False(t, ok)
}

func TestLRange(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.RPush("list", "0", "1", "2", "3", "4", "5")

	assert.Equal(t, []string{"1", "2", "3"}, e.LRange("list", 1, 3))
	assert.Equal(t, []string{"0", "1", "2", "3", "4", "5"}, e.LRange("list", 0, -1))
	assert.Equal(t, []string{"4", "5"}, e.LRange("list", -2, -1))
}

func TestLLen(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	assert.Equal(t, 0, e.LLen("missing"))
	e.RPush("list", "a", "b")
	assert.Equal(t, 2, e.LLen("list"))
}

func TestLIndex(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.RPush("list", "a", "b", "c")

	val, ok := e.LIndex("list", 0)
	require.True(t, ok)
	assert.Equal(t, "a", val)

	val, ok = e.LIndex("list", -1)
	require.True(t, ok)
	assert.Equal(t, "c", val)

	_, ok = e.LIndex("list", 10)
	assert.False(t, ok)
}

func TestLRem(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.RPush("list", "a", "b", "a", "c", "a")

	removed := e.LRem("list", 0, "a")
	assert.Equal(t, 3, removed)
	assert.Equal(t, []string{"b", "c"}, e.LRange("list", 0, -1))
}

func TestLRemHead(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.RPush("list", "a", "b", "a", "c", "a")

	removed := e.LRem("list", 2, "a")
	assert.Equal(t, 2, removed)
	assert.Equal(t, []string{"b", "c", "a"}, e.LRange("list", 0, -1))
}

func TestLRemTail(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.RPush("list", "a", "b", "a", "c", "a")

	removed := e.LRem("list", -2, "a")
	assert.Equal(t, 2, removed)
	assert.Equal(t, []string{"a", "b", "c"}, e.LRange("list", 0, -1))
}

func TestLTrim(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.RPush("list", "0", "1", "2", "3", "4")
	e.LTrim("list", 1, 3)
	assert.Equal(t, []string{"1", "2", "3"}, e.LRange("list", 0, -1))
}

func TestLTrimDeletesEmpty(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.RPush("list", "a")
	e.LTrim("list", 1, 5)
	assert.Equal(t, 0, e.LLen("list"))
}

// --- Set tests ---

func TestSAddSMembers(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	added := e.SAdd("set", "a", "b", "c", "a")
	assert.Equal(t, 3, added)

	members := e.SMembers("set")
	assert.Equal(t, []string{"a", "b", "c"}, members)
}

func TestSRem(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("set", "a", "b", "c")
	removed := e.SRem("set", "a", "missing")
	assert.Equal(t, 1, removed)
	assert.Equal(t, []string{"b", "c"}, e.SMembers("set"))
}

func TestSRemDeletesEmpty(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("set", "a")
	e.SRem("set", "a")
	assert.Equal(t, 0, e.SCard("set"))
}

func TestSIsMember(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("set", "a", "b")
	assert.True(t, e.SIsMember("set", "a"))
	assert.False(t, e.SIsMember("set", "c"))
	assert.False(t, e.SIsMember("missing", "a"))
}

func TestSCard(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	assert.Equal(t, 0, e.SCard("missing"))
	e.SAdd("set", "a", "b", "c")
	assert.Equal(t, 3, e.SCard("set"))
}

func TestSPop(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("set", "a", "b", "c", "d", "e")
	result := e.SPop("set", 2)
	require.Len(t, result, 2)
	assert.Equal(t, 3, e.SCard("set"))
}

func TestSPopAll(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("set", "a", "b")
	result := e.SPop("set", 5)
	assert.Len(t, result, 2)
	assert.Equal(t, 0, e.SCard("set"))
}

func TestSRandMember(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("set", "a", "b", "c")
	result := e.SRandMember("set", 2)
	assert.Len(t, result, 2)
	assert.Equal(t, 3, e.SCard("set")) // Not removed.
}

func TestSUnion(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("s1", "a", "b", "c")
	e.SAdd("s2", "b", "c", "d")
	result := e.SUnion("s1", "s2")
	assert.Equal(t, []string{"a", "b", "c", "d"}, result)
}

func TestSInter(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("s1", "a", "b", "c")
	e.SAdd("s2", "b", "c", "d")
	result := e.SInter("s1", "s2")
	assert.Equal(t, []string{"b", "c"}, result)
}

func TestSInterEmpty(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("s1", "a")
	e.SAdd("s2", "b")
	result := e.SInter("s1", "s2")
	assert.Nil(t, result)
}

func TestSDiff(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.SAdd("s1", "a", "b", "c")
	e.SAdd("s2", "b", "c", "d")
	result := e.SDiff("s1", "s2")
	assert.Equal(t, []string{"a"}, result)
}

// --- Hash tests ---

func TestHSetHGet(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	added := e.HSet("hash", []string{"field1", "value1", "field2", "value2"})
	assert.Equal(t, 2, added)

	val, ok := e.HGet("hash", "field1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)

	_, ok = e.HGet("hash", "missing")
	assert.False(t, ok)
}

func TestHSetUpdatesExisting(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.HSet("hash", []string{"field", "old"})
	added := e.HSet("hash", []string{"field", "new"})
	assert.Equal(t, 0, added)

	val, _ := e.HGet("hash", "field")
	assert.Equal(t, "new", val)
}

func TestHGetAll(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.HSet("hash", []string{"a", "1", "b", "2"})
	result := e.HGetAll("hash")
	// Order is not guaranteed, so just check length and presence.
	assert.Len(t, result, 4)
}

func TestHDel(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.HSet("hash", []string{"a", "1", "b", "2", "c", "3"})
	deleted := e.HDel("hash", "a", "missing")
	assert.Equal(t, 1, deleted)
	assert.Equal(t, 2, e.HLen("hash"))
}

func TestHDelDeletesEmpty(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.HSet("hash", []string{"a", "1"})
	e.HDel("hash", "a")
	assert.Equal(t, 0, e.HLen("hash"))
}

func TestHLen(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	assert.Equal(t, 0, e.HLen("missing"))
	e.HSet("hash", []string{"a", "1", "b", "2"})
	assert.Equal(t, 2, e.HLen("hash"))
}

func TestHExists(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.HSet("hash", []string{"a", "1"})
	assert.True(t, e.HExists("hash", "a"))
	assert.False(t, e.HExists("hash", "b"))
}

func TestHKeysHVals(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.HSet("hash", []string{"a", "1", "b", "2"})
	assert.Equal(t, []string{"a", "b"}, e.HKeys("hash"))
	assert.Equal(t, []string{"1", "2"}, e.HVals("hash"))
}

func TestHMGet(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.HSet("hash", []string{"a", "1", "b", "2"})
	result := e.HMGet("hash", "a", "missing", "b")
	assert.Equal(t, []string{"1", "", "2"}, result)
}

// --- Sorted Set tests ---

func TestZAddZRange(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	added := e.ZAdd("zset", map[string]float64{"b": 2.0, "a": 1.0, "c": 3.0})
	assert.Equal(t, 3, added)

	result := e.ZRange("zset", 0, -1, false)
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestZRangeWithScores(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.ZAdd("zset", map[string]float64{"b": 2.0, "a": 1.0})
	result := e.ZRange("zset", 0, -1, true)
	assert.Equal(t, []string{"a", "1", "b", "2"}, result)
}

func TestZRevRange(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.ZAdd("zset", map[string]float64{"a": 1.0, "b": 2.0, "c": 3.0})
	result := e.ZRevRange("zset", 0, -1, false)
	assert.Equal(t, []string{"c", "b", "a"}, result)
}

func TestZRangeByScore(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.ZAdd("zset", map[string]float64{"a": 1.0, "b": 2.0, "c": 3.0, "d": 4.0})
	result := e.ZRangeByScore("zset", 2.0, 3.5, false, 0, -1)
	assert.Equal(t, []string{"b", "c"}, result)
}

func TestZRangeByScoreWithLimit(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.ZAdd("zset", map[string]float64{"a": 1.0, "b": 2.0, "c": 3.0, "d": 4.0})
	result := e.ZRangeByScore("zset", 1.0, 4.0, false, 1, 2)
	assert.Equal(t, []string{"b", "c"}, result)
}

func TestZRem(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.ZAdd("zset", map[string]float64{"a": 1.0, "b": 2.0})
	removed := e.ZRem("zset", "a", "missing")
	assert.Equal(t, 1, removed)
	assert.Equal(t, []string{"b"}, e.ZRange("zset", 0, -1, false))
}

func TestZCard(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	assert.Equal(t, 0, e.ZCard("missing"))
	e.ZAdd("zset", map[string]float64{"a": 1.0, "b": 2.0})
	assert.Equal(t, 2, e.ZCard("zset"))
}

func TestZScore(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.ZAdd("zset", map[string]float64{"a": 1.5})
	score, ok := e.ZScore("zset", "a")
	require.True(t, ok)
	assert.Equal(t, 1.5, score)

	_, ok = e.ZScore("zset", "missing")
	assert.False(t, ok)
}

func TestZIncrBy(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	newScore := e.ZIncrBy("zset", 2.0, "a")
	assert.Equal(t, 2.0, newScore)

	newScore = e.ZIncrBy("zset", 3.0, "a")
	assert.Equal(t, 5.0, newScore)
}

func TestZCount(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.ZAdd("zset", map[string]float64{"a": 1.0, "b": 2.0, "c": 3.0})
	assert.Equal(t, 2, e.ZCount("zset", 1.5, 3.0))
}

func TestZRemRangeByScore(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.ZAdd("zset", map[string]float64{"a": 1.0, "b": 2.0, "c": 3.0})
	removed := e.ZRemRangeByScore("zset", 1.5, 2.5)
	assert.Equal(t, 1, removed)
	assert.Equal(t, []string{"a", "c"}, e.ZRange("zset", 0, -1, false))
}

func TestZRemRangeByRank(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.ZAdd("zset", map[string]float64{"a": 1.0, "b": 2.0, "c": 3.0, "d": 4.0})
	removed := e.ZRemRangeByRank("zset", 1, 2)
	assert.Equal(t, 2, removed)
	assert.Equal(t, []string{"a", "d"}, e.ZRange("zset", 0, -1, false))
}

// --- Key-level tests ---

func TestKeys(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("user:1", "a", 0)
	e.Set("user:2", "b", 0)
	e.Set("post:1", "c", 0)

	assert.Equal(t, []string{"user:1", "user:2"}, e.Keys("user:*"))
	assert.Equal(t, []string{"post:1"}, e.Keys("post:?"))
}

func TestFlushDB(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("a", "1", 0)
	e.Set("b", "2", 0)
	e.FlushDB()
	assert.Equal(t, 0, e.DBSize())
}

func TestDBSize(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	assert.Equal(t, 0, e.DBSize())
	e.Set("a", "1", 0)
	e.Set("b", "2", 0)
	assert.Equal(t, 2, e.DBSize())
}

func TestType(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("str", "val", 0)
	e.RPush("list", "a")
	e.SAdd("set", "a")
	e.HSet("hash", []string{"f", "v"})
	e.ZAdd("zset", map[string]float64{"a": 1.0})

	assert.Equal(t, "string", e.Type("str"))
	assert.Equal(t, "list", e.Type("list"))
	assert.Equal(t, "set", e.Type("set"))
	assert.Equal(t, "hash", e.Type("hash"))
	assert.Equal(t, "zset", e.Type("zset"))
	assert.Equal(t, "none", e.Type("missing"))
}

func TestRename(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("old", "val", 0)
	err := e.Rename("old", "new")
	require.NoError(t, err)

	_, ok := e.Get("old")
	assert.False(t, ok)
	val, ok := e.Get("new")
	assert.True(t, ok)
	assert.Equal(t, "val", val)
}

func TestRenameNX(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("old", "val", 0)
	e.Set("existing", "other", 0)

	ok, err := e.RenameNX("old", "existing")
	require.NoError(t, err)
	assert.False(t, ok)

	ok, err = e.RenameNX("old", "new")
	require.NoError(t, err)
	assert.True(t, ok)
}

// --- Expiration tests ---

func TestLazyExpiration(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("key", "val", 50*time.Millisecond)
	time.Sleep(60 * time.Millisecond)

	// Lazy expiration on read.
	_, ok := e.Get("key")
	assert.False(t, ok)
}

func TestActiveExpiration(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	// Set many keys with short TTL.
	for i := 0; i < 100; i++ {
		e.Set("key"+string(rune(i)), "val", 50*time.Millisecond)
	}

	// Wait for active expiration to run.
	time.Sleep(300 * time.Millisecond)

	// Most should be expired by now.
	assert.Less(t, e.DBSize(), 50)
}

func TestExpiredKeyDeletedFromType(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	e.Set("key", "val", 50*time.Millisecond)
	time.Sleep(60 * time.Millisecond)

	assert.Equal(t, "none", e.Type("key"))
	assert.Equal(t, int64(-2), e.TTL("key"))
}

// --- Concurrent tests ---

func TestConcurrentSetGet(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	// Run many goroutines doing sets and gets.
	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- struct{}{} }()
			for j := 0; j < 100; j++ {
				key := "key" + string(rune(id))
				e.Set(key, "value", 0)
				e.Get(key)
			}
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}

// --- Glob matching tests ---

func TestMatchGlob(t *testing.T) {
	tests := []struct {
		s       string
		pattern string
		want    bool
	}{
		{"hello", "hello", true},
		{"hello", "h*o", true},
		{"hello", "h?llo", true},
		{"hello", "*", true},
		{"hello", "world", false},
		{"hello", "h?ll", false},
		{"", "", true},
		{"", "*", true},
		{"abc", "a*c", true},
		{"abcdef", "a*d*f", true},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, matchGlob(tt.s, tt.pattern), "matchGlob(%q, %q)", tt.s, tt.pattern)
	}
}

// --- Memory limit tests ---

func TestMaxMemoryCanWrite(t *testing.T) {
	e := NewEngineWithMaxMemory(10)
	defer e.Stop()

	assert.True(t, e.CanWrite(5))
	assert.False(t, e.CanWrite(15))

	// Write 8 bytes: "key" (3) + "value" (5).
	e.Set("key", "value", 0)

	// 8 + 3 = 11 > 10, so should be false.
	assert.False(t, e.CanWrite(3))
}

func TestMaxMemoryUnlimited(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	// Unlimited memory should always allow writes.
	for i := 0; i < 1000; i++ {
		e.Set(fmt.Sprintf("key%d", i), "value", 0)
	}
}

func TestAvailableMemory(t *testing.T) {
	mem, err := AvailableMemory()
	// Should not error on supported platforms.
	assert.NoError(t, err)
	// On Linux/macOS, should return >0. On others, returns 0.
	assert.GreaterOrEqual(t, mem, uint64(0))
}

func TestEstimateItemSize(t *testing.T) {
	// String
	size := estimateItemSize("key", Item{Value: "value", Typ: TypeString})
	assert.Equal(t, int64(3+5), size) // "key" + "value"

	// List
	size = estimateItemSize("key", Item{Value: []string{"a", "bb"}, Typ: TypeList})
	assert.Equal(t, int64(3+1+2), size)

	// Set
	size = estimateItemSize("key", Item{Value: map[string]struct{}{"a": {}, "bb": {}}, Typ: TypeSet})
	assert.Equal(t, int64(3+1+2), size)

	// Hash
	size = estimateItemSize("key", Item{Value: map[string]string{"f": "v"}, Typ: TypeHash})
	assert.Equal(t, int64(3+1+1), size)

	// ZSet
	size = estimateItemSize("key", Item{Value: map[string]float64{"m": 1.0}, Typ: TypeZSet})
	assert.Equal(t, int64(3+1+8), size)
}

func TestCheckMemory(t *testing.T) {
	e := NewEngineWithMaxMemory(10)
	defer e.Stop()

	// Should pass when under limit.
	err := e.checkMemory(5)
	assert.NoError(t, err)

	// Fill memory.
	e.Set("key", "value", 0) // 8 bytes

	// Should fail when over limit.
	err = e.checkMemory(5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "OOM")
}

func TestCheckMemoryUnlimited(t *testing.T) {
	e := NewEngine()
	defer e.Stop()

	// Unlimited should always allow.
	err := e.checkMemory(1000000)
	assert.NoError(t, err)
}
