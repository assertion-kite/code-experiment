package cache

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNewLRUCache(t *testing.T) {
	// 创建一个LRU缓存实例，容量为10，
	lruCache := NewLRUCache(10)

	// 创建一个等待组来等待所有goroutines完成
	var wg sync.WaitGroup

	// 启动多个goroutines来并发地访问缓存
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				// 生成一个随机的键
				key := rand.Intn(100)
				// 模拟从某个数据源获取值
				value := rand.Intn(1000)

				// 将键值对放入缓存
				lruCache.Put(key, value)

				// 休眠一小段时间以增加并发性
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

				// 从缓存中获取值
				if cachedValue := lruCache.Get(key); cachedValue != -1 {
					fmt.Printf("Goroutine %d: Got value %d for key %d from cache\n", id, cachedValue, key)
				} else {
					fmt.Printf("Goroutine %d: Key %d not found in cache\n", id, key)
				}
			}
		}(i)
	}

	// 等待所有goroutines完成
	wg.Wait()

	// 打印最终缓存的容量（可能由于扩展而增加）
	fmt.Printf("Final cache capacity: %d\n", lruCache.capacity)
}
