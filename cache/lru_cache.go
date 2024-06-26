package cache

import (
	"container/list"
	"sync"
)

type Pair struct {
	key   int
	value int
}

type LRUCache struct {
	// mu 读写锁，用于保护并发访问。
	mu sync.RWMutex
	// cache 映射，用于快速查找缓存中的条目。
	cache map[int]*list.Element
	// list 双向链表，用于存储缓存条目的访问顺序。
	list *list.List
	// capacity 最大空间
	capacity int
}

// NewLRUCache 初始化
func NewLRUCache(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		cache:    make(map[int]*list.Element),
		list:     list.New(),
	}
}

// Get 查询
func (c *LRUCache) Get(key int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if elem, ok := c.cache[key]; ok {
		// 移动到前面
		c.list.MoveToFront(elem)
		return elem.Value.(*Pair).value
	}
	return -1
}

// Put 创建参数
func (c *LRUCache) Put(key int, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*Pair).value = value
		return
	}

	pair := &Pair{key: key, value: value}
	elem := c.list.PushFront(pair)
	c.cache[key] = elem

	for c.list.Len() > c.capacity {
		back := c.list.Back()
		if back != nil {
			// 删除最后一个
			c.list.Remove(back)
			delete(c.cache, back.Value.(*Pair).key)
		}
	}
}
