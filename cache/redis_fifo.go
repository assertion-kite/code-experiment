package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

// Cache 定义了一个缓存接口
type Cache interface {
	Push(ctx context.Context, key string, value interface{}) error
	Pop(ctx context.Context, key string, result interface{}) error
	Close() error
}

// RedisCache 实现了Cache接口，使用Redis作为底层存储
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建一个新的RedisCache实例
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Push 将值推送到Redis列表的尾部（队尾）
func (c *RedisCache) Push(ctx context.Context, key string, value interface{}) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = c.client.LPush(ctx, key, val).Result()
	return err
}

// Pop 从Redis列表的头部（队头）移除并返回值
func (c *RedisCache) Pop(ctx context.Context, key string, result interface{}) error {
	val, err := c.client.RPop(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return errors.New("queue is empty")
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(val, result)
	if err != nil {
		return errors.New("json unmarshal fail")
	}
	return nil
}

// Close 关闭Redis连接（如果需要的话）
func (c *RedisCache) Close() error {
	// 通常，Redis客户端库不需要显式关闭连接
	// 但如果需要，可以在这里实现关闭逻辑
	return nil
}
