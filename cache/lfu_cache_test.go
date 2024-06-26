package cache

import "testing"

func TestNewLFUCache(t *testing.T) {
	// 示例用法
	cache := NewLFUCache(2)
	cache.Put(1, 1)
	cache.Put(2, 2)
	val, found := cache.Get(1)
	if found {
		println("Got 1:", val)
	}
	cache.Put(3, 3) // 淘汰键2
	val, found = cache.Get(2)
	if !found {
		println("2 not found")
	}
	val, found = cache.Get(3)
	if found {
		println("Got 3:", val)
	}
	cache.Put(4, 4) // 淘汰键1
	val, found = cache.Get(1)
	if !found {
		println("1 not found")
	}
	val, found = cache.Get(3)
	if found {
		println("Got 3:", val)
	}
	val, found = cache.Get(4)
	if found {
		println("Got 4:", val)
	}
}
