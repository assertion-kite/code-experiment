package lock

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestNewRedisLock(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	hTable := NewRedisLock(redisCli, "order_table", "1")
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			hTable.Lock()
			defer hTable.Unlock()
			fmt.Println(i)
			time.Sleep(1 * time.Second)
		}()
	}

	time.Sleep(20 * time.Second)
}
