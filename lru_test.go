package lru

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkLRU_Rand(b *testing.B) {
	l, err := New(8192)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		trace[i] = rand.Int63() % 32768
	}

	b.ResetTimer()

	var hit, miss int
	for i := 0; i < 2*b.N; i++ {
		if i%2 == 0 {
			l.Add(trace[i], trace[i])
		} else {
			_, ok := l.Get(trace[i])
			if ok {
				hit++
			} else {
				miss++
			}
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func BenchmarkLRU_Freq(b *testing.B) {
	l, err := New(8192)
	if err != nil {
		b.Fatalf("err: %v", err)
	}

	trace := make([]int64, b.N*2)
	for i := 0; i < b.N*2; i++ {
		if i%2 == 0 {
			trace[i] = rand.Int63() % 16384
		} else {
			trace[i] = rand.Int63() % 32768
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l.Add(trace[i], trace[i])
	}
	var hit, miss int
	for i := 0; i < b.N; i++ {
		_, ok := l.Get(trace[i])
		if ok {
			hit++
		} else {
			miss++
		}
	}
	b.Logf("hit: %d miss: %d ratio: %f", hit, miss, float64(hit)/float64(miss))
}

func TestLRU(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k interface{}, v interface{}) {
		if k != v {
			t.Fatalf("Evict values not equal (%v!=%v)", k, v)
		}
		evictCounter += 1
	}
	l, err := NewWithEvict(128, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for i := 0; i < 256; i++ {
		l.Add(i, i)
	}
	if l.Len() != 128 {
		t.Fatalf("bad len: %v", l.Len())
	}

	if evictCounter != 128 {
		t.Fatalf("bad evict count: %v", evictCounter)
	}

	for i, k := range l.Keys() {
		if v, ok := l.Get(k); !ok || v != k || v != i+128 {
			t.Fatalf("bad key: %v", k)
		}
	}
	for i := 0; i < 128; i++ {
		_, ok := l.Get(i)
		if ok {
			t.Fatalf("should be evicted")
		}
	}
	for i := 128; i < 256; i++ {
		_, ok := l.Get(i)
		if !ok {
			t.Fatalf("should not be evicted")
		}
	}
	for i := 128; i < 192; i++ {
		l.Remove(i)
		_, ok := l.Get(i)
		if ok {
			t.Fatalf("should be deleted")
		}
	}

	l.Get(192) // expect 192 to be last key in l.Keys()

	for i, k := range l.Keys() {
		if (i < 63 && k != i+193) || (i == 63 && k != 192) {
			t.Fatalf("out of order key: %v", k)
		}
	}

	l.Purge()
	if l.Len() != 0 {
		t.Fatalf("bad len: %v", l.Len())
	}
	if _, ok := l.Get(200); ok {
		t.Fatalf("should contain nothing")
	}
}

// test that Add returns true/false if an eviction occurred
func TestLRUAdd(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k interface{}, v interface{}) {
		evictCounter += 1
	}

	l, err := NewWithEvict(1, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if l.Add(1, 1) == true || evictCounter != 0 {
		t.Errorf("should not have an eviction")
	}
	if l.Add(2, 2) == false || evictCounter != 1 {
		t.Errorf("should have an eviction")
	}
}

// test that Contains doesn't update recent-ness
func TestLRUContains(t *testing.T) {
	l, err := New(2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add(1, 1)
	l.Add(2, 2)
	if !l.Contains(1) {
		t.Errorf("1 should be contained")
	}

	l.Add(3, 3)
	if l.Contains(1) {
		t.Errorf("Contains should not have updated recent-ness of 1")
	}
}

// test that Contains doesn't update recent-ness
func TestLRUContainsOrAdd(t *testing.T) {
	l, err := New(2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add(1, 1)
	l.Add(2, 2)
	contains, evict := l.ContainsOrAdd(1, 1)
	if !contains {
		t.Errorf("1 should be contained")
	}
	if evict {
		t.Errorf("nothing should be evicted here")
	}

	l.Add(3, 3)
	contains, evict = l.ContainsOrAdd(1, 1)
	if contains {
		t.Errorf("1 should not have been contained")
	}
	if !evict {
		t.Errorf("an eviction should have occurred")
	}
	if !l.Contains(1) {
		t.Errorf("now 1 should be contained")
	}
}

// test that Peek doesn't update recent-ness
func TestLRUPeek(t *testing.T) {
	l, err := New(2)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	l.Add(1, 1)
	l.Add(2, 2)
	if v, ok := l.Peek(1); !ok || v != 1 {
		t.Errorf("1 should be set to 1: %v, %v", v, ok)
	}

	l.Add(3, 3)
	if l.Contains(1) {
		t.Errorf("should not have updated recent-ness of 1")
	}
}

// test that NewWithExpireAndEvict works as expected
func TestNewWithExpireAndEvict(t *testing.T) {
	evictCounter := 0
	onEvicted := func(k interface{}, v interface{}) {
		evictCounter += 1
	}
	

	expire := 200 * time.Millisecond
	l, err := NewWithExpireAndEvict(3, expire, onEvicted)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	//test evict by expiration
	l.Add(1, 1)

	time.Sleep(201 * time.Millisecond) //eviction counter 1

	l.Add(2, 2)
	l.Add(3, 3)

	//test length
	if l.Len() != 3 {
		t.Fatalf("Wrong size!")
	}
	
	//simple get
	var observed interface{}
	var observedInt int
	var ok bool
	if observed, ok = l.Get(2); !ok {
		t.Fatalf("expected value does not exist in cache")
	}
	observedInt, ok = observed.(int)
	if !ok {
		t.Fatalf("type assertion failed expected value")
	}

	if observedInt != 2 {
		t.Fatalf("value retrieved from cache does equal expected value")
	}


	//test evict by manual removal
	l.Remove(2) // increment eviction counter = 2


	//test evict by surpassing capacity
	l.Add(4, 4)
	l.Add(5, 5) 
	l.Add(6, 6)
	//increment eviction counter = 3

	if evictCounter != 3 {
		t.Fatalf("evict handler did not increment correctly: count %v", evictCounter)
	}

}
