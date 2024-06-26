package cache

import (
	"container/list"
	"math"
	"sync"
)

// LFUPair 表示键值对
type LFUPair struct {
	key   interface{}
	value interface{}
	count int // 访问次数
}

// LFUCache 是LFU缓存的结构体
type LFUCache struct {
	mu       sync.RWMutex
	capacity int
	minFreq  int                           // 当前缓存中最小的访问频率
	freqMap  map[int]*list.List            // 频率到双向链表的映射，链表保存具有相同频率的键值对
	keyMap   map[interface{}]*list.Element // 键到链表元素的映射
}

// NewLFUCache 创建一个新的LFU缓存
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		minFreq:  1,
		freqMap:  make(map[int]*list.List),
		keyMap:   make(map[interface{}]*list.Element),
	}
}

// Get 从缓存中获取一个值，如果键不存在则返回-1
func (c *LFUCache) Get(key interface{}) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if elem, ok := c.keyMap[key]; ok {
		pair := elem.Value.(*LFUPair)
		pair.count++
		// 更新频率
		c.updateFreq(elem, pair.count)
		return pair.value, true
	}
	return -1, false
}

// Put 向缓存中添加一个键值对
func (c *LFUCache) Put(key, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.keyMap[key]; ok {
		// 如果键已存在，更新值并更新频率
		pair := elem.Value.(*LFUPair)
		pair.value = value
		pair.count++
		c.updateFreq(elem, pair.count)
		return
	}

	// 如果缓存已满，则淘汰一个最小频率的键值对
	if c.capacity <= len(c.keyMap) {
		c.evict()
	}

	// 插入新元素
	pair := &LFUPair{key: key, value: value, count: 1}
	elem := c.getOrCreateFreqList(1).PushFront(pair)
	c.keyMap[key] = elem
	c.minFreq = 1 // 插入新元素后，最小频率至少为1
}

// updateFreq 更新元素频率并维护freqMap和minFreq
func (c *LFUCache) updateFreq(elem *list.Element, newFreq int) {
	if elem == nil {
		// 如果elem是nil，直接返回或抛出错误，取决于你的需求
		return
	}

	oldPair := elem.Value.(*LFUPair) // 假设LFUPair是你的键值对结构体
	oldFreq := oldPair.count

	// 从旧频率的链表中移除元素
	oldList, exists := c.freqMap[oldFreq]
	if !exists || oldList == nil {
		// 如果oldFreq对应的链表不存在，可能是并发问题或数据不一致，应该处理这种情况
		// 这里简单返回或抛出错误
		return
	}
	oldList.Remove(elem)

	// 如果旧频率的链表为空，从freqMap中删除它，并可能需要更新minFreq
	if oldList.Len() == 0 {
		delete(c.freqMap, oldFreq)
		// 如果删除的是minFreq链表，更新minFreq
		if oldFreq == c.minFreq {
			// 寻找新的minFreq
			c.minFreq = c.findNewMinFreq()
		}
	}

	// 将元素添加到新频率的链表中
	newList := c.getOrCreateFreqList(newFreq)
	newList.PushFront(elem)
	oldPair.count = newFreq
}

// findNewMinFreq 辅助函数，用于找到新的最小频率
func (c *LFUCache) findNewMinFreq() int {
	minFreq := math.MaxInt32 // 初始化为一个很大的数
	for freq, list := range c.freqMap {
		if list.Len() > 0 && freq < minFreq {
			minFreq = freq
		}
	}
	if minFreq == math.MaxInt32 {
		// 如果没有找到有效的minFreq（即所有链表都为空），则返回1或根据需要处理
		return 1
	}
	return minFreq
}

func (c *LFUCache) getOrCreateFreqList(freq int) *list.List {
	if list1, ok := c.freqMap[freq]; ok {
		return list1
	}

	// 如果该频率的链表不存在，则创建它
	list1 := list.New()
	c.freqMap[freq] = list1
	return list1
}

// evict 淘汰一个最小频率的键值对
func (c *LFUCache) evict() {
	// 确保minFreq对应的链表不为空
	minList := c.freqMap[c.minFreq]
	if minList == nil || minList.Len() == 0 {
		return // 如果没有minFreq链表或链表为空，则直接返回
	}

	// 淘汰minFreq链表中的最后一个元素
	elem := minList.Back()
	if elem == nil {
		return // 这通常不会发生，但为了完整性还是检查一下
	}
	pair := elem.Value.(*LFUPair)
	delete(c.keyMap, pair.key)
	minList.Remove(elem)

	// 如果minFreq链表为空，则更新minFreq
	if minList.Len() == 0 {
		delete(c.freqMap, c.minFreq)
		if c.minFreq > 1 {
			c.minFreq--
		}
	}
}
