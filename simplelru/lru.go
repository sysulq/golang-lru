package simplelru

import (
	"container/list"
	"errors"
	"time"

	"github.com/hnlq715/golang-lru/shardmap"
)

// EvictCallback is used to get a callback when a cache entry is evicted
type EvictCallback func(key string, value interface{})

// LRU implements a non-thread safe fixed size LRU cache
type LRU struct {
	size      int
	evictList *list.List
	items     *shardmap.Map
	expire    time.Duration
	onEvict   EvictCallback
}

// entry is used to hold a value in the evictList
type entry struct {
	key    string
	value  interface{}
	expire *time.Time
}

func (e *entry) IsExpired() bool {
	if e.expire == nil {
		return false
	}
	return time.Now().After(*e.expire)
}

// NewLRU constructs an LRU of the given size
func NewLRU(size int, onEvict EvictCallback) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	c := &LRU{
		size:      size,
		evictList: list.New(),
		items:     shardmap.New(size),
		expire:    0,
		onEvict:   onEvict,
	}
	return c, nil
}

// NewLRUWithExpire contrusts an LRU of the given size and expire time
func NewLRUWithExpire(size int, expire time.Duration, onEvict EvictCallback) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	c := &LRU{
		size:      size,
		evictList: list.New(),
		items:     shardmap.New(size),
		expire:    expire,
		onEvict:   onEvict,
	}
	return c, nil
}

// Purge is used to completely clear the cache
func (c *LRU) Purge() {
	c.items.Range(func(k string, v interface{}) bool {
		if c.onEvict != nil {
			c.onEvict(k, v.(*list.Element).Value.(*entry).value)
		}
		return true
	})
	c.items.Clear()
	c.evictList.Init()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *LRU) Add(key string, value interface{}) bool {
	return c.AddEx(key, value, 0)
}

// AddEx adds a value to the cache with expire.  Returns true if an eviction occurred.
func (c *LRU) AddEx(key string, value interface{}, expire time.Duration) bool {
	var ex *time.Time = nil
	if expire > 0 {
		expire := time.Now().Add(expire)
		ex = &expire
	} else if c.expire > 0 {
		expire := time.Now().Add(c.expire)
		ex = &expire
	}
	// Check for existing item
	if ent, ok := c.get(key); ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		ent.Value.(*entry).expire = ex
		return false
	}

	// Add new item
	ent := &entry{key: key, value: value, expire: ex}
	entry := c.evictList.PushFront(ent)
	c.set(key, entry)

	evict := c.evictList.Len() > c.size
	// Verify size not exceeded
	if evict {
		c.removeOldest()
	}
	return evict
}

// Get looks up a key's value from the cache.
func (c *LRU) Get(key string) (value interface{}, ok bool) {
	if ent, ok := c.get(key); ok {
		if ent.Value.(*entry).IsExpired() {
			return nil, false
		}
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (c *LRU) Contains(key string) (ok bool) {
	if ent, ok := c.get(key); ok {
		if ent.Value.(*entry).IsExpired() {
			return false
		}
		return ok
	}
	return
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (c *LRU) Peek(key string) (value interface{}, ok bool) {
	if ent, ok := c.get(key); ok {
		if ent.Value.(*entry).IsExpired() {
			return nil, false
		}
		return ent.Value.(*entry).value, true
	}
	return nil, ok
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRU) Remove(key string) bool {
	if ent, ok := c.get(key); ok {
		c.removeElement(ent)
		return true
	}
	return false
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU) RemoveOldest() (string, interface{}, bool) {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return "", nil, false
}

// GetOldest returns the oldest entry
func (c *LRU) GetOldest() (string, interface{}, bool) {
	ent := c.evictList.Back()
	if ent != nil {
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return "", nil, false
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (c *LRU) Keys() []string {
	keys := make([]string, c.evictList.Len())
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

// Len returns the number of items in the cache.
func (c *LRU) Len() int {
	return c.evictList.Len()
}

// Resize changes the cache size.
func (c *LRU) Resize(size int) (evicted int) {
	diff := c.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		c.removeOldest()
	}
	c.size = size
	return diff
}

// removeOldest removes the oldest item from the cache.
func (c *LRU) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

// removeElement is used to remove a given list element from the cache
func (c *LRU) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	c.items.Delete(kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}

func (c *LRU) get(key string) (*list.Element, bool) {
	item, ok := c.items.Get(key)
	if !ok {
		return nil, false
	}

	return item.(*list.Element), true
}

func (c *LRU) set(key string, value *list.Element) {
	c.items.Set(key, value)
}
