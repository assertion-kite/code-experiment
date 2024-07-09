package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
)

func TestNewRedisCache(t *testing.T) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	cache := NewRedisCache(rdb)
	var res struct {
		Id int64
	}
	cache.Push(ctx, "queue", struct {
		Id int64
	}{
		Id: 1,
	})
	err := cache.Pop(ctx, "queue", &res)
	t.Log(err)
	fmt.Println(res)
	cache.Push(ctx, "queue", 2)
	cache.Push(ctx, "queue", 3)

}
