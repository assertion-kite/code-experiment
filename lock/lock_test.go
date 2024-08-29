package lock

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	redisCli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	c := NewClient(redisCli)
	l, err := c.TryLock(context.TODO(), "test", time.Second*10)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(l.Unlock(context.TODO()))

}
