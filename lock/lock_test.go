package lock

import (
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestName(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	NewClient(redisCli).TryLock()
}
