package shardmap

import (
	"runtime"
	"sync"

	"github.com/cespare/xxhash"
)

// Map is a hashmap. Like map[string]interface{}, but sharded and thread-safe.
type Map struct {
	init   sync.Once
	cap    int
	shards int
	seed   uint32
	mus    []sync.RWMutex
	maps   []map[string]interface{}
}

// New returns a new hashmap with the specified capacity. This function is only
// needed when you must define a minimum capacity, otherwise just use:
//    var m shardmap.Map
func New(cap int) *Map {
	return &Map{cap: cap}
}

// Clear out all values from map
func (m *Map) Clear() {
	m.initDo()
	for i := 0; i < m.shards; i++ {
		m.mus[i].Lock()
		m.maps[i] = make(map[string]interface{}, m.cap/m.shards)
		m.mus[i].Unlock()
	}
}

// Set assigns a value to a key.
// Returns the previous value, or false when no value was assigned.
func (m *Map) Set(key string, value interface{}) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	m.maps[shard][key] = value
	m.mus[shard].Unlock()
	return
}

// Get returns a value for a key.
// Returns false when no value has been assign for key.
func (m *Map) Get(key string) (value interface{}, ok bool) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].RLock()
	value, ok = m.maps[shard][key]
	m.mus[shard].RUnlock()
	return value, ok
}

// Delete deletes a value for a key.
// Returns the deleted value, or false when no value was assigned.
func (m *Map) Delete(key string) {
	m.initDo()
	shard := m.choose(key)
	m.mus[shard].Lock()
	delete(m.maps[shard], key)
	m.mus[shard].Unlock()
}

// Len returns the number of values in map.
func (m *Map) Len() int {
	m.initDo()
	var lens int
	for i := 0; i < m.shards; i++ {
		m.mus[i].Lock()
		lens += len(m.maps[i])
		m.mus[i].Unlock()
	}
	return lens
}

// Range iterates overall all key/values.
// It's not safe to call or Set or Delete while ranging.
func (m *Map) Range(iter func(key string, value interface{}) bool) {
	m.initDo()
	var done bool
	for i := 0; i < m.shards; i++ {
		func() {
			m.mus[i].RLock()
			defer m.mus[i].RUnlock()
			for key, value := range m.maps[i] {
				if !iter(key, value) {
					done = true
					return
				}
			}
		}()
		if done {
			break
		}
	}
}

func (m *Map) choose(key string) int {
	return int(xxhash.Sum64String(key) & uint64(m.shards-1))
}

func (m *Map) initDo() {
	m.init.Do(func() {
		m.shards = 1
		for m.shards < runtime.NumCPU()*16 {
			m.shards *= 2
		}
		scap := m.cap / m.shards
		m.mus = make([]sync.RWMutex, m.shards)
		m.maps = make([]map[string]interface{}, m.shards)
		for i := 0; i < len(m.maps); i++ {
			m.maps[i] = make(map[string]interface{}, scap)
		}
	})
}
