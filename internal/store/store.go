// Package store implements the in-memory data engine.
//
// The engine stores typed values (strings, lists, sets, hashes, sorted sets)
// with optional expiration. All operations are thread-safe via a global
// read-write mutex.
//
// Expiration uses two mechanisms:
//  1. Lazy: check on every read; delete if expired.
//  2. Active: background goroutine samples keys periodically.
package store

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DataType represents the type of a stored value.
type DataType int

const (
	TypeString DataType = iota
	TypeList
	TypeSet
	TypeHash
	TypeZSet
)

func (d DataType) String() string {
	switch d {
	case TypeString:
		return "string"
	case TypeList:
		return "list"
	case TypeSet:
		return "set"
	case TypeHash:
		return "hash"
	case TypeZSet:
		return "zset"
	default:
		return "none"
	}
}

// Item is a stored value with metadata.
type Item struct {
	Value     interface{}
	Typ       DataType
	ExpiresAt *time.Time // nil means no expiration
}

// Engine is the core in-memory storage engine.
type Engine struct {
	mu   sync.RWMutex
	data map[string]Item

	// Background expiration ticker.
	ticker *time.Ticker
	stop   chan struct{}

	// maxMemory is the maximum memory in bytes (0 = unlimited).
	maxMemory int64
}

// NewEngine creates a new storage engine with unlimited memory.
func NewEngine() *Engine {
	return NewEngineWithMaxMemory(0)
}

// NewEngineWithMaxMemory creates a new storage engine and starts the active expiration goroutine.
func NewEngineWithMaxMemory(maxMemory int64) *Engine {
	e := &Engine{
		data:      make(map[string]Item),
		stop:      make(chan struct{}),
		maxMemory: maxMemory,
	}
	e.ticker = time.NewTicker(100 * time.Millisecond)
	go e.activeExpiration()
	return e
}

// Stop halts the background expiration goroutine.
func (e *Engine) Stop() {
	e.ticker.Stop()
	close(e.stop)
}

// estimateItemSize returns an approximate memory size for an item in bytes.
func estimateItemSize(key string, item Item) int64 {
	size := int64(len(key))
	switch item.Typ {
	case TypeString:
		size += int64(len(item.Value.(string)))
	case TypeList:
		if l, ok := item.Value.(*List); ok {
			for i := 0; i < l.Len(); i++ {
				v, _ := l.Get(i)
				size += int64(len(v))
			}
		}
	case TypeSet:
		for m := range item.Value.(map[string]struct{}) {
			size += int64(len(m))
		}
	case TypeHash:
		for f, v := range item.Value.(map[string]string) {
			size += int64(len(f) + len(v))
		}
	case TypeZSet:
		for m := range item.Value.(map[string]float64) {
			size += int64(len(m)) + 8 // float64 = 8 bytes
		}
	}
	return size
}

// currentMemoryUsage returns the approximate total memory used by all keys.
func (e *Engine) currentMemoryUsage() int64 {
	var total int64
	for k, item := range e.data {
		total += estimateItemSize(k, item)
	}
	return total
}

// checkMemory returns an error if a write would exceed maxMemory.
// Must be called with write lock held.
func (e *Engine) checkMemory(additionalBytes int64) error {
	if e.maxMemory <= 0 {
		return nil
	}
	if e.currentMemoryUsage()+additionalBytes > e.maxMemory {
		return fmt.Errorf("OOM command not allowed when used memory > 'maxmemory'")
	}
	return nil
}

// CanWrite returns true if the engine can accept a write of approximately
// additionalBytes without exceeding maxMemory.
func (e *Engine) CanWrite(additionalBytes int64) bool {
	if e.maxMemory <= 0 {
		return true
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.currentMemoryUsage()+additionalBytes <= e.maxMemory
}

// AvailableMemory returns the total system memory in bytes.
// Returns 0 if it cannot be determined on the current platform.
func AvailableMemory() (uint64, error) {
	return availableMemory()
}

// ---------- Key-level operations ----------

// Get retrieves a string value by key.
// Returns the value and true if the key exists and is a string.
func (e *Engine) Get(key string) (string, bool) {
	e.mu.RLock()
	item, ok := e.data[key]
	e.mu.RUnlock()
	if !ok {
		return "", false
	}
	if e.isExpired(key, item) {
		return "", false
	}
	if item.Typ != TypeString {
		return "", false
	}
	return item.Value.(string), true
}

// Set stores a string value with optional TTL (<=0 means no expiration).
func (e *Engine) Set(key string, value string, ttl time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()
	item := Item{Value: value, Typ: TypeString}
	if ttl > 0 {
		t := time.Now().Add(ttl)
		item.ExpiresAt = &t
	}
	e.data[key] = item
}

// Del deletes keys and returns the count deleted.
func (e *Engine) Del(keys ...string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	deleted := 0
	for _, key := range keys {
		if _, ok := e.data[key]; ok {
			delete(e.data, key)
			deleted++
		}
	}
	return deleted
}

// Exists returns the number of existing, non-expired keys.
func (e *Engine) Exists(keys ...string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	count := 0
	for _, key := range keys {
		if item, ok := e.data[key]; ok && !e.isExpiredLocked(key, item) {
			count++
		}
	}
	return count
}

// Expire sets a TTL on an existing key. Returns true if the key exists.
func (e *Engine) Expire(key string, ttl time.Duration) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) {
		return false
	}
	t := time.Now().Add(ttl)
	item.ExpiresAt = &t
	e.data[key] = item
	return true
}

// TTL returns the remaining TTL for a key, or -1 if no TTL, or -2 if key does not exist.
func (e *Engine) TTL(key string) int64 {
	e.mu.RLock()
	item, ok := e.data[key]
	e.mu.RUnlock()
	if !ok || e.isExpired(key, item) {
		return -2
	}
	if item.ExpiresAt == nil {
		return -1
	}
	remaining := time.Until(*item.ExpiresAt)
	if remaining <= 0 {
		return -2
	}
	return int64(remaining.Milliseconds())
}

// Persist removes the expiration from a key. Returns true if the key exists.
func (e *Engine) Persist(key string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) {
		return false
	}
	item.ExpiresAt = nil
	e.data[key] = item
	return true
}

// Keys returns all keys matching a glob-like pattern (* = any sequence, ? = any char).
func (e *Engine) Keys(pattern string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	var matches []string
	for key, item := range e.data {
		if e.isExpiredLocked(key, item) {
			continue
		}
		if matchGlob(key, pattern) {
			matches = append(matches, key)
		}
	}
	sort.Strings(matches)
	return matches
}

// FlushDB deletes all keys.
func (e *Engine) FlushDB() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data = make(map[string]Item)
}

// DBSize returns the number of keys.
func (e *Engine) DBSize() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	count := 0
	for key, item := range e.data {
		if !e.isExpiredLocked(key, item) {
			count++
		}
	}
	return count
}

// Type returns the type of a key, or "none" if it does not exist.
func (e *Engine) Type(key string) string {
	e.mu.RLock()
	item, ok := e.data[key]
	e.mu.RUnlock()
	if !ok || e.isExpired(key, item) {
		return "none"
	}
	return item.Typ.String()
}

// Rename renames a key. Returns an error if the source does not exist.
func (e *Engine) Rename(key, newkey string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) {
		return fmt.Errorf("no such key")
	}
	delete(e.data, key)
	e.data[newkey] = item
	return nil
}

// RenameNX renames a key only if the new key does not exist.
func (e *Engine) RenameNX(key, newkey string) (bool, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) {
		return false, fmt.Errorf("no such key")
	}
	if _, exists := e.data[newkey]; exists {
		return false, nil
	}
	delete(e.data, key)
	e.data[newkey] = item
	return true, nil
}

// ---------- String operations ----------

// Append appends a value to a string key. Returns the new length.
// If the key does not exist, it is created.
func (e *Engine) Append(key, value string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) || item.Typ != TypeString {
		// Create new string.
		e.data[key] = Item{Value: value, Typ: TypeString}
		return len(value)
	}
	newVal := item.Value.(string) + value
	item.Value = newVal
	e.data[key] = item
	return len(newVal)
}

// StrLen returns the length of a string value, or 0 if not a string.
func (e *Engine) StrLen(key string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) || item.Typ != TypeString {
		return 0
	}
	return len(item.Value.(string))
}

// IncrBy increments a string key by a delta. Returns the new value.
func (e *Engine) IncrBy(key string, delta int64) (int64, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	item, ok := e.data[key]
	var current int64
	if ok && !e.isExpiredLocked(key, item) {
		if item.Typ != TypeString {
			return 0, fmt.Errorf("value is not an integer")
		}
		var err error
		current, err = parseInt(item.Value.(string))
		if err != nil {
			return 0, fmt.Errorf("value is not an integer")
		}
	}
	current += delta
	e.data[key] = Item{Value: fmt.Sprintf("%d", current), Typ: TypeString}
	return current, nil
}

// MGet returns values for multiple keys. Non-existent keys return empty strings.
func (e *Engine) MGet(keys ...string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]string, len(keys))
	for i, key := range keys {
		item, ok := e.data[key]
		if ok && !e.isExpiredLocked(key, item) && item.Typ == TypeString {
			result[i] = item.Value.(string)
		}
	}
	return result
}

// MSet sets multiple key-value pairs.
func (e *Engine) MSet(pairs []string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i := 0; i < len(pairs)-1; i += 2 {
		e.data[pairs[i]] = Item{Value: pairs[i+1], Typ: TypeString}
	}
}

// ---------- List operations ----------

// getList returns the list stored at key. When create is true and the key
// does not exist, an empty list is returned so callers can build it up.
func (e *Engine) getList(key string, create bool) (*List, bool) {
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) {
		if create {
			return NewList(), true
		}
		return nil, false
	}
	if item.Typ != TypeList {
		return nil, false
	}
	list, ok := item.Value.(*List)
	return list, ok
}

func (e *Engine) saveList(key string, list *List) {
	e.data[key] = Item{Value: list, Typ: TypeList}
}

// LPush pushes values to the head of a list (O(1) per element).
// Multiple values are pushed in argument order, so the last argument
// ends up at the head — matching Redis semantics. Returns the new length.
func (e *Engine) LPush(key string, values ...string) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	list, ok := e.getList(key, true)
	if !ok {
		return 0
	}
	if len(values) > 0 {
		for _, v := range values {
			list.PushFront(v)
		}
		e.saveList(key, list)
	}
	return list.Len()
}

// RPush pushes values to the tail of a list (O(1) per element).
// Returns the new length.
func (e *Engine) RPush(key string, values ...string) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	list, ok := e.getList(key, true)
	if !ok {
		return 0
	}
	if len(values) > 0 {
		for _, v := range values {
			list.PushBack(v)
		}
		e.saveList(key, list)
	}
	return list.Len()
}

// LPop pops from the head of a list. Returns the value and true if it existed.
func (e *Engine) LPop(key string) (string, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	list, ok := e.getList(key, false)
	if !ok || list.Len() == 0 {
		return "", false
	}
	val, _ := list.PopFront()
	if list.Len() == 0 {
		delete(e.data, key)
	} else {
		e.saveList(key, list)
	}
	return val, true
}

// RPop pops from the tail of a list. Returns the value and true if it existed.
func (e *Engine) RPop(key string) (string, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()

	list, ok := e.getList(key, false)
	if !ok || list.Len() == 0 {
		return "", false
	}
	val, _ := list.PopBack()
	if list.Len() == 0 {
		delete(e.data, key)
	} else {
		e.saveList(key, list)
	}
	return val, true
}

// LRange returns a range of elements from a list.
func (e *Engine) LRange(key string, start, stop int) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	list, ok := e.getList(key, false)
	if !ok {
		return nil
	}
	n := list.Len()
	start, stop = normalizeRange(start, stop, n)
	if start > stop {
		return nil
	}
	result := make([]string, stop-start+1)
	for i := start; i <= stop; i++ {
		result[i-start], _ = list.Get(i)
	}
	return result
}

// LLen returns the length of a list.
func (e *Engine) LLen(key string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()

	list, ok := e.getList(key, false)
	if !ok {
		return 0
	}
	return list.Len()
}

// LIndex returns the element at an index (negative counts from the tail).
func (e *Engine) LIndex(key string, index int) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	list, ok := e.getList(key, false)
	if !ok {
		return "", false
	}
	return list.Get(index)
}

// LRem removes elements equal to value from a list.
// count > 0: remove first N from head. count < 0: remove first N from tail.
// count = 0: remove all. Returns the number removed.
func (e *Engine) LRem(key string, count int, value string) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	list, ok := e.getList(key, false)
	if !ok {
		return 0
	}
	all := list.Copy()
	removed := 0
	var kept []string

	switch {
	case count == 0:
		for _, v := range all {
			if v == value {
				removed++
			} else {
				kept = append(kept, v)
			}
		}
	case count > 0:
		for _, v := range all {
			if v == value && removed < count {
				removed++
			} else {
				kept = append(kept, v)
			}
		}
	default:
		// Remove from tail: iterate backwards, keep in reverse, then un-reverse.
		target := -count
		rev := make([]string, 0, len(all))
		for i := len(all) - 1; i >= 0; i-- {
			if all[i] == value && removed < target {
				removed++
			} else {
				rev = append(rev, all[i])
			}
		}
		kept = make([]string, len(rev))
		for i := range rev {
			kept[len(rev)-1-i] = rev[i]
		}
	}

	if len(kept) == 0 {
		delete(e.data, key)
	} else {
		newList := NewList()
		for _, v := range kept {
			newList.PushBack(v)
		}
		e.saveList(key, newList)
	}
	return removed
}

// LTrim trims a list to the specified range.
func (e *Engine) LTrim(key string, start, stop int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	list, ok := e.getList(key, false)
	if !ok {
		return
	}
	n := list.Len()
	start, stop = normalizeRange(start, stop, n)
	if start > stop {
		delete(e.data, key)
		return
	}
	kept := make([]string, stop-start+1)
	for i := start; i <= stop; i++ {
		kept[i-start], _ = list.Get(i)
	}
	newList := NewList()
	for _, v := range kept {
		newList.PushBack(v)
	}
	e.saveList(key, newList)
}

// ---------- Set operations ----------

func (e *Engine) getSet(key string, create bool) (map[string]struct{}, bool) {
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) {
		if create {
			return make(map[string]struct{}), true
		}
		return nil, false
	}
	if item.Typ != TypeSet {
		return nil, false
	}
	return item.Value.(map[string]struct{}), true
}

func (e *Engine) saveSet(key string, set map[string]struct{}) {
	e.data[key] = Item{Value: set, Typ: TypeSet}
}

// SAdd adds members to a set. Returns the number of new members added.
func (e *Engine) SAdd(key string, members ...string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	set, ok := e.getSet(key, true)
	if !ok {
		return 0
	}
	added := 0
	for _, m := range members {
		if _, exists := set[m]; !exists {
			set[m] = struct{}{}
			added++
		}
	}
	e.saveSet(key, set)
	return added
}

// SRem removes members from a set. Returns the number removed.
func (e *Engine) SRem(key string, members ...string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	set, ok := e.getSet(key, false)
	if !ok {
		return 0
	}
	removed := 0
	for _, m := range members {
		if _, exists := set[m]; exists {
			delete(set, m)
			removed++
		}
	}
	if len(set) == 0 {
		delete(e.data, key)
	} else {
		e.saveSet(key, set)
	}
	return removed
}

// SMembers returns all members of a set.
func (e *Engine) SMembers(key string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	set, ok := e.getSet(key, false)
	if !ok {
		return nil
	}
	members := make([]string, 0, len(set))
	for m := range set {
		members = append(members, m)
	}
	sort.Strings(members)
	return members
}

// SIsMember checks if a member exists in a set.
func (e *Engine) SIsMember(key, member string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	set, ok := e.getSet(key, false)
	if !ok {
		return false
	}
	_, exists := set[member]
	return exists
}

// SCard returns the cardinality of a set.
func (e *Engine) SCard(key string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	set, ok := e.getSet(key, false)
	if !ok {
		return 0
	}
	return len(set)
}

// SPop removes and returns random members from a set.
func (e *Engine) SPop(key string, count int) []string {
	e.mu.Lock()
	defer e.mu.Unlock()
	set, ok := e.getSet(key, false)
	if !ok {
		return nil
	}
	members := make([]string, 0, len(set))
	for m := range set {
		members = append(members, m)
	}
	if count >= len(members) {
		delete(e.data, key)
		sort.Strings(members)
		return members
	}
	// Shuffle and pick first 'count'.
	rand.Shuffle(len(members), func(i, j int) {
		members[i], members[j] = members[j], members[i]
	})
	result := members[:count]
	for _, m := range result {
		delete(set, m)
	}
	e.saveSet(key, set)
	sort.Strings(result)
	return result
}

// SRandMember returns random members without removing.
func (e *Engine) SRandMember(key string, count int) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	set, ok := e.getSet(key, false)
	if !ok {
		return nil
	}
	members := make([]string, 0, len(set))
	for m := range set {
		members = append(members, m)
	}
	if count >= len(members) {
		sort.Strings(members)
		return members
	}
	// Shuffle copy and pick.
	shuffled := make([]string, len(members))
	copy(shuffled, members)
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	result := shuffled[:count]
	sort.Strings(result)
	return result
}

// SUnion returns the union of multiple sets.
func (e *Engine) SUnion(keys ...string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	union := make(map[string]struct{})
	for _, key := range keys {
		set, ok := e.getSet(key, false)
		if !ok {
			continue
		}
		for m := range set {
			union[m] = struct{}{}
		}
	}
	result := make([]string, 0, len(union))
	for m := range union {
		result = append(result, m)
	}
	sort.Strings(result)
	return result
}

// SInter returns the intersection of multiple sets.
func (e *Engine) SInter(keys ...string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if len(keys) == 0 {
		return nil
	}
	// Start with the first set.
	first, ok := e.getSet(keys[0], false)
	if !ok {
		return nil
	}
	inter := make(map[string]struct{})
	for m := range first {
		inter[m] = struct{}{}
	}
	for _, key := range keys[1:] {
		set, ok := e.getSet(key, false)
		if !ok {
			return nil
		}
		for m := range inter {
			if _, exists := set[m]; !exists {
				delete(inter, m)
			}
		}
	}
	if len(inter) == 0 {
		return nil
	}
	result := make([]string, 0, len(inter))
	for m := range inter {
		result = append(result, m)
	}
	sort.Strings(result)
	return result
}

// SDiff returns the difference between the first set and all others.
func (e *Engine) SDiff(keys ...string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if len(keys) == 0 {
		return nil
	}
	first, ok := e.getSet(keys[0], false)
	if !ok {
		return nil
	}
	diff := make(map[string]struct{})
	for m := range first {
		diff[m] = struct{}{}
	}
	for _, key := range keys[1:] {
		set, ok := e.getSet(key, false)
		if !ok {
			continue
		}
		for m := range set {
			delete(diff, m)
		}
	}
	result := make([]string, 0, len(diff))
	for m := range diff {
		result = append(result, m)
	}
	sort.Strings(result)
	return result
}

// ---------- Hash operations ----------

func (e *Engine) getHash(key string, create bool) (map[string]string, bool) {
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) {
		if create {
			return make(map[string]string), true
		}
		return nil, false
	}
	if item.Typ != TypeHash {
		return nil, false
	}
	return item.Value.(map[string]string), true
}

func (e *Engine) saveHash(key string, hash map[string]string) {
	e.data[key] = Item{Value: hash, Typ: TypeHash}
}

// HSet sets hash fields. Returns the number of new fields added.
func (e *Engine) HSet(key string, pairs []string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	hash, ok := e.getHash(key, true)
	if !ok {
		return 0
	}
	added := 0
	for i := 0; i < len(pairs)-1; i += 2 {
		field, value := pairs[i], pairs[i+1]
		if _, exists := hash[field]; !exists {
			added++
		}
		hash[field] = value
	}
	e.saveHash(key, hash)
	return added
}

// HGet returns a hash field value.
func (e *Engine) HGet(key, field string) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	hash, ok := e.getHash(key, false)
	if !ok {
		return "", false
	}
	val, exists := hash[field]
	return val, exists
}

// HGetAll returns all fields and values.
func (e *Engine) HGetAll(key string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	hash, ok := e.getHash(key, false)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(hash)*2)
	for field, value := range hash {
		result = append(result, field, value)
	}
	return result
}

// HDel deletes hash fields. Returns the number deleted.
func (e *Engine) HDel(key string, fields ...string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	hash, ok := e.getHash(key, false)
	if !ok {
		return 0
	}
	deleted := 0
	for _, field := range fields {
		if _, exists := hash[field]; exists {
			delete(hash, field)
			deleted++
		}
	}
	if len(hash) == 0 {
		delete(e.data, key)
	} else {
		e.saveHash(key, hash)
	}
	return deleted
}

// HLen returns the number of fields in a hash.
func (e *Engine) HLen(key string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	hash, ok := e.getHash(key, false)
	if !ok {
		return 0
	}
	return len(hash)
}

// HExists checks if a field exists.
func (e *Engine) HExists(key, field string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	hash, ok := e.getHash(key, false)
	if !ok {
		return false
	}
	_, exists := hash[field]
	return exists
}

// HKeys returns all field names.
func (e *Engine) HKeys(key string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	hash, ok := e.getHash(key, false)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(hash))
	for field := range hash {
		result = append(result, field)
	}
	sort.Strings(result)
	return result
}

// HVals returns all field values.
func (e *Engine) HVals(key string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	hash, ok := e.getHash(key, false)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(hash))
	for _, value := range hash {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

// HMGet returns values for multiple fields.
func (e *Engine) HMGet(key string, fields ...string) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	hash, ok := e.getHash(key, false)
	if !ok {
		result := make([]string, len(fields))
		return result
	}
	result := make([]string, len(fields))
	for i, field := range fields {
		if val, exists := hash[field]; exists {
			result[i] = val
		}
	}
	return result
}

// ---------- Sorted Set operations ----------

// zsetMember pairs a member with its score.
type zsetMember struct {
	member string
	score  float64
}

func (e *Engine) getZSet(key string, create bool) (map[string]float64, bool) {
	item, ok := e.data[key]
	if !ok || e.isExpiredLocked(key, item) {
		if create {
			return make(map[string]float64), true
		}
		return nil, false
	}
	if item.Typ != TypeZSet {
		return nil, false
	}
	return item.Value.(map[string]float64), true
}

func (e *Engine) saveZSet(key string, zset map[string]float64) {
	e.data[key] = Item{Value: zset, Typ: TypeZSet}
}

// sortedZSetMembers returns members sorted by score ascending, then member name.
func sortedZSetMembers(zset map[string]float64) []zsetMember {
	members := make([]zsetMember, 0, len(zset))
	for m, s := range zset {
		members = append(members, zsetMember{member: m, score: s})
	}
	sort.Slice(members, func(i, j int) bool {
		if members[i].score != members[j].score {
			return members[i].score < members[j].score
		}
		return members[i].member < members[j].member
	})
	return members
}

// ZAdd adds members with scores. Returns the number of new members.
func (e *Engine) ZAdd(key string, members map[string]float64) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	zset, ok := e.getZSet(key, true)
	if !ok {
		return 0
	}
	added := 0
	for m, s := range members {
		if _, exists := zset[m]; !exists {
			added++
		}
		zset[m] = s
	}
	e.saveZSet(key, zset)
	return added
}

// ZRem removes members. Returns the number removed.
func (e *Engine) ZRem(key string, members ...string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return 0
	}
	removed := 0
	for _, m := range members {
		if _, exists := zset[m]; exists {
			delete(zset, m)
			removed++
		}
	}
	if len(zset) == 0 {
		delete(e.data, key)
	} else {
		e.saveZSet(key, zset)
	}
	return removed
}

// ZRange returns members by rank (0-based, inclusive).
func (e *Engine) ZRange(key string, start, stop int, withScores bool) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return nil
	}
	members := sortedZSetMembers(zset)
	start, stop = normalizeRange(start, stop, len(members))
	if start > stop {
		return nil
	}
	result := make([]string, 0, (stop-start+1)*(map[bool]int{true: 2, false: 1}[withScores]))
	for i := start; i <= stop; i++ {
		result = append(result, members[i].member)
		if withScores {
			result = append(result, fmt.Sprintf("%g", members[i].score))
		}
	}
	return result
}

// ZRevRange returns members by rank in reverse order.
func (e *Engine) ZRevRange(key string, start, stop int, withScores bool) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return nil
	}
	members := sortedZSetMembers(zset)
	// Reverse.
	for i, j := 0, len(members)-1; i < j; i, j = i+1, j-1 {
		members[i], members[j] = members[j], members[i]
	}
	start, stop = normalizeRange(start, stop, len(members))
	if start > stop {
		return nil
	}
	result := make([]string, 0, (stop-start+1)*(map[bool]int{true: 2, false: 1}[withScores]))
	for i := start; i <= stop; i++ {
		result = append(result, members[i].member)
		if withScores {
			result = append(result, fmt.Sprintf("%g", members[i].score))
		}
	}
	return result
}

// ZRangeByScore returns members within a score range.
func (e *Engine) ZRangeByScore(key string, min, max float64, withScores bool, offset, count int) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return nil
	}
	members := sortedZSetMembers(zset)
	var filtered []zsetMember
	for _, m := range members {
		if m.score >= min && m.score <= max {
			filtered = append(filtered, m)
		}
	}
	if offset >= len(filtered) {
		return nil
	}
	end := offset + count
	if end > len(filtered) || count < 0 {
		end = len(filtered)
	}
	result := make([]string, 0, (end-offset)*(map[bool]int{true: 2, false: 1}[withScores]))
	for i := offset; i < end; i++ {
		result = append(result, filtered[i].member)
		if withScores {
			result = append(result, fmt.Sprintf("%g", filtered[i].score))
		}
	}
	return result
}

// ZRevRangeByScore returns members within a score range in reverse order.
func (e *Engine) ZRevRangeByScore(key string, max, min float64, withScores bool, offset, count int) []string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return nil
	}
	members := sortedZSetMembers(zset)
	// Reverse.
	for i, j := 0, len(members)-1; i < j; i, j = i+1, j-1 {
		members[i], members[j] = members[j], members[i]
	}
	var filtered []zsetMember
	for _, m := range members {
		if m.score <= max && m.score >= min {
			filtered = append(filtered, m)
		}
	}
	if offset >= len(filtered) {
		return nil
	}
	end := offset + count
	if end > len(filtered) || count < 0 {
		end = len(filtered)
	}
	result := make([]string, 0, (end-offset)*(map[bool]int{true: 2, false: 1}[withScores]))
	for i := offset; i < end; i++ {
		result = append(result, filtered[i].member)
		if withScores {
			result = append(result, fmt.Sprintf("%g", filtered[i].score))
		}
	}
	return result
}

// ZCard returns the cardinality.
func (e *Engine) ZCard(key string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return 0
	}
	return len(zset)
}

// ZScore returns a member's score.
func (e *Engine) ZScore(key, member string) (float64, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return 0, false
	}
	score, exists := zset[member]
	return score, exists
}

// ZIncrBy increments a member's score. Returns the new score.
func (e *Engine) ZIncrBy(key string, increment float64, member string) float64 {
	e.mu.Lock()
	defer e.mu.Unlock()
	zset, ok := e.getZSet(key, true)
	if !ok {
		return 0
	}
	zset[member] += increment
	e.saveZSet(key, zset)
	return zset[member]
}

// ZCount returns the number of members within a score range.
func (e *Engine) ZCount(key string, min, max float64) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return 0
	}
	count := 0
	for _, score := range zset {
		if score >= min && score <= max {
			count++
		}
	}
	return count
}

// ZRemRangeByScore removes members within a score range. Returns count removed.
func (e *Engine) ZRemRangeByScore(key string, min, max float64) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return 0
	}
	removed := 0
	for m, score := range zset {
		if score >= min && score <= max {
			delete(zset, m)
			removed++
		}
	}
	if len(zset) == 0 {
		delete(e.data, key)
	} else {
		e.saveZSet(key, zset)
	}
	return removed
}

// ZRemRangeByRank removes members by rank range. Returns count removed.
func (e *Engine) ZRemRangeByRank(key string, start, stop int) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	zset, ok := e.getZSet(key, false)
	if !ok {
		return 0
	}
	members := sortedZSetMembers(zset)
	start, stop = normalizeRange(start, stop, len(members))
	if start > stop {
		return 0
	}
	removed := 0
	for i := start; i <= stop; i++ {
		delete(zset, members[i].member)
		removed++
	}
	if len(zset) == 0 {
		delete(e.data, key)
	} else {
		e.saveZSet(key, zset)
	}
	return removed
}

// ---------- Expiration ----------

// isExpired checks if an item is expired and deletes it if so.
// Must be called without holding the write lock (acquires lock internally).
func (e *Engine) isExpired(key string, item Item) bool {
	if item.ExpiresAt == nil || time.Now().Before(*item.ExpiresAt) {
		return false
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	// Double-check under write lock.
	if current, ok := e.data[key]; ok && current.ExpiresAt != nil && time.Now().After(*current.ExpiresAt) {
		delete(e.data, key)
	}
	return true
}

// isExpiredLocked checks if an item is expired. Must be called with lock held.
func (e *Engine) isExpiredLocked(key string, item Item) bool {
	if item.ExpiresAt == nil || time.Now().Before(*item.ExpiresAt) {
		return false
	}
	delete(e.data, key)
	return true
}

// activeExpiration runs a background goroutine that periodically samples
// and deletes expired keys using Redis's probabilistic algorithm.
func (e *Engine) activeExpiration() {
	for {
		select {
		case <-e.ticker.C:
			e.expireSample()
		case <-e.stop:
			return
		}
	}
}

// expireSample implements Redis's active expiration algorithm:
// Sample 20 random keys, delete expired ones. If >25% expired, repeat.
func (e *Engine) expireSample() {
	e.mu.Lock()
	defer e.mu.Unlock()

	if len(e.data) == 0 {
		return
	}

	// Run up to 20 cycles per tick (Redis does this too).
	for cycle := 0; cycle < 20; cycle++ {
		expired := 0
		sampled := 0

		// Collect keys to sample.
		keys := make([]string, 0, len(e.data))
		for k := range e.data {
			keys = append(keys, k)
		}

		// Sample up to 20 random keys.
		for i := 0; i < 20 && i < len(keys); i++ {
			idx := rand.Intn(len(keys))
			key := keys[idx]
			item := e.data[key]
			sampled++
			if item.ExpiresAt != nil && time.Now().After(*item.ExpiresAt) {
				delete(e.data, key)
				expired++
			}
			// Remove sampled key to avoid re-sampling.
			keys[idx] = keys[len(keys)-1]
			keys = keys[:len(keys)-1]
		}

		if sampled == 0 || float64(expired)/float64(sampled) < 0.25 {
			break
		}
	}
}

// ---------- Helpers ----------

// normalizeRange converts Redis-style start/stop indices to Go slice indices.
// Negative indices count from the end. Out-of-bounds values are clamped.
func normalizeRange(start, stop, length int) (int, int) {
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}
	if start > stop {
		return start, stop
	}
	return start, stop
}

// normalizeIndex converts a Redis-style index to a Go slice index.
func normalizeIndex(index, length int) int {
	if index < 0 {
		index = length + index
	}
	return index
}

// parseInt parses a string as int64.
func parseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// matchGlob checks if a string matches a simple glob pattern.
// Supports * (any sequence) and ? (single char).
func matchGlob(s, pattern string) bool {
	// Fast path: exact match or all-match.
	if pattern == "*" {
		return true
	}
	if !strings.ContainsAny(pattern, "*?") {
		return s == pattern
	}
	// Simple recursive glob matcher.
	return globMatch([]rune(s), []rune(pattern))
}

func globMatch(s, p []rune) bool {
	for len(p) > 0 {
		switch p[0] {
		case '*':
			// Try all suffixes.
			for i := 0; i <= len(s); i++ {
				if globMatch(s[i:], p[1:]) {
					return true
				}
			}
			return false
		case '?':
			if len(s) == 0 {
				return false
			}
			s, p = s[1:], p[1:]
		default:
			if len(s) == 0 || s[0] != p[0] {
				return false
			}
			s, p = s[1:], p[1:]
		}
	}
	return len(s) == 0
}
