package memdb

import (
	"sync"

	"github.com/hsn/tiny-redis/util"
)

const MaxConSize = 1<<31 - 1

// shard represents a shard within the ConcurrentMap.
// It stores key-value pairs in a map and uses a read-write mutex for concurrent access control.
type shard struct {
	mp   map[string]any // Map to store key-value pairs.
	rwMu *sync.RWMutex  // Mutex for concurrent access control.
}

// ConcurrentMap manages a table of shards with multiple hashmap shards.
// It provides thread-safe operations using read-write locks.
// The maximum table size is limited to MaxConSize.
type ConcurrentMap struct {
	table []*shard // Array of hashmap shards.
	size  int      // Maximum table size.
	count int      // Total number of keys across all shards.
}

// NewConcurrentMap creates a new ConcurrentMap
func NewConcurrentMap(size int) *ConcurrentMap {
	if size > MaxConSize || size <= 0 {
		size = MaxConSize
	}
	m := &ConcurrentMap{
		table: make([]*shard, size),
		size:  size,
		count: 0,
	}
	for i := 0; i < size; i++ {
		m.table[i] = &shard{
			mp:   make(map[string]any),
			rwMu: &sync.RWMutex{},
		}
	}
	return m
}

// getKeyPos calculates the shard index based on the hash of the key.
func (m *ConcurrentMap) getKeyPos(key string) int {
	return util.HashKey(key) % m.size
}

// Set stores a key-value pair in the ConcurrentMap.
// if successful return 1
// else return 0
func (m *ConcurrentMap) Set(key string, value any) int {
	added := 0
	pos := m.getKeyPos(key)
	shard := m.table[pos]
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()
	if _, ok := shard.mp[key]; !ok {
		m.count++
		added = 1
	}
	shard.mp[key] = value
	return added
}

// SetIfExist stores a key-value pair in the ConcurrentMap
// if it exists replace it with the new value return 1
// else return 0
func (m *ConcurrentMap) SetIfExist(key string, value any) int {
	pos := m.getKeyPos(key)
	shard := m.table[pos]
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()
	if _, ok := shard.mp[key]; ok {
		shard.mp[key] = value
		return 1
	}
	return 0
}
func (m *ConcurrentMap) SetIfNotExist(key string, value any) int {
	pos := m.getKeyPos(key)
	shard := m.table[pos]
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()
	if _, ok := shard.mp[key]; !ok {
		m.count++
		shard.mp[key] = value
		return 1
	}
	return 0
}
func (m *ConcurrentMap) Get(key string) (any, bool) {
	pos := m.getKeyPos(key)
	shard := m.table[pos]
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()
	value, ok := shard.mp[key]
	return value, ok
}
func (m *ConcurrentMap) Delete(key string) int {
	pos := m.getKeyPos(key)
	shard := m.table[pos]
	shard.rwMu.Lock()
	defer shard.rwMu.Unlock()
	if _, ok := shard.mp[key]; ok {
		delete(shard.mp, key)
		m.count--
		return 1
	}
	return 0
}
func (m *ConcurrentMap) Len() int {
	return m.count
}
func (m *ConcurrentMap) Clear() {
	*m = *NewConcurrentMap(m.size)
}
func (m *ConcurrentMap) Keys() []string {
	keys := make([]string, m.count)
	i := 0
	for _, shard := range m.table {
		shard.rwMu.RLock()
		for key := range shard.mp {
			keys[i] = key
			i++
		}
		shard.rwMu.RUnlock()
	}
	return keys
}
