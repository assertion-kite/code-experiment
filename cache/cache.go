package cache

import (
	"code/clock"
	"container/heap"
	"sync"
	"time"
)

// NewExpiring 初始化
func NewExpiring() *Expiring {
	return NewExpiringWithClock(clock.RealClock{})
}

// NewExpiringWithClock 添加可控时间参数
func NewExpiringWithClock(clock clock.Clock) *Expiring {
	return &Expiring{
		clock: clock,
		cache: make(map[interface{}]entry),
	}
}

// Expiring 缓存实现
type Expiring struct {
	// clock 时钟实现
	clock clock.Clock
	// mu 锁
	mu sync.RWMutex
	// cache 缓存的map
	cache map[interface{}]entry
	// generation 当前缓存的版本号
	generation uint64
	// heap 一个最小堆
	heap expiringHeap
}

type entry struct {
	val        interface{}
	expiry     time.Time
	generation uint64
}

// Get 在缓存中查找条目
func (c *Expiring) Get(key interface{}) (val interface{}, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.cache[key]
	if !ok || !c.clock.Now().Before(e.expiry) {
		return nil, false
	}
	return e.val, true
}

// Set 设置映射中的键/值/到期条目，覆盖以前的任何条目
func (c *Expiring) Set(key interface{}, val interface{}, ttl time.Duration) {
	now := c.clock.Now()
	expiry := now.Add(ttl)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.generation++

	c.cache[key] = entry{
		val:        val,
		expiry:     expiry,
		generation: c.generation,
	}

	// 在推送新条目之前，以内联方式运行GC
	c.gc(now)

	heap.Push(&c.heap, &expiringHeapEntry{
		key:        key,
		expiry:     expiry,
		generation: c.generation,
	})
}

// Delete 删除缓存条目
func (c *Expiring) Delete(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.del(key, 0)
}

// del 删除给定密钥的条目
func (c *Expiring) del(key interface{}, generation uint64) {
	e, ok := c.cache[key]
	if !ok {
		return
	}
	if generation != 0 && generation != e.generation {
		return
	}
	delete(c.cache, key)
}

// Len 返回缓存中的项数
func (c *Expiring) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// gc 垃圾回收方法。它遍历堆中的元素，检查哪些元素已过期，并从缓存中删除它们。
// 这个过程会一直持续，直到堆中的第一个元素未过期为止。
// 管理过期的缓存条目。
// 会占用内存资源，直到被显式地删除。
// 导致可用内存减少，性能下降，甚至可能引发内存泄漏。
func (c *Expiring) gc(now time.Time) {
	for {
		if len(c.heap) == 0 || now.Before(c.heap[0].expiry) {
			return
		}
		cleanup := heap.Pop(&c.heap).(*expiringHeapEntry)
		c.del(cleanup.key, cleanup.generation)
	}
}

type expiringHeapEntry struct {
	key        interface{}
	expiry     time.Time
	generation uint64
}

// expiringHeap 一个堆
type expiringHeap []*expiringHeapEntry

var _ heap.Interface = &expiringHeap{}

// Len 返回堆中元素的数量。
func (cq expiringHeap) Len() int {
	return len(cq)
}

// Less 比较两个元素的过期时间，以确定它们在堆中的顺序。这里，它比较的是 expiry 字段，因此这是一个最小堆，堆顶元素总是具有最早的过期时间。
func (cq expiringHeap) Less(i, j int) bool {
	return cq[i].expiry.Before(cq[j].expiry)
}

// Swap 交换两个元素的位置。
func (cq expiringHeap) Swap(i, j int) {
	cq[i], cq[j] = cq[j], cq[i]
}

// Push 方法向堆中添加一个新元素。这里，它假设传入的 c 是一个指向 expiringHeapEntry 的指针，并将其添加到切片中。
func (cq *expiringHeap) Push(c interface{}) {
	*cq = append(*cq, c.(*expiringHeapEntry))
}

// Pop 方法从堆中移除并返回具有最早过期时间的元素（即堆顶元素）。它首先获取最后一个元素，然后移除并返回
func (cq *expiringHeap) Pop() interface{} {
	c := (*cq)[cq.Len()-1]
	*cq = (*cq)[:cq.Len()-1]
	return c
}
