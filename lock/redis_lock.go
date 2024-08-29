package lock

import (
	"code/log"
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// lua脚本用于保证Redis加锁&设置过期时间的原子性操作
//EX second ：设置键的过期时间为 second 秒。 SET key value EX second 效果等同于 SETEX key second value 。
//PX millisecond ：设置键的过期时间为 millisecond 毫秒。 SET key value PX millisecond 效果等同于 PSETEX key millisecond value 。
//NX ：只在键不存在时，才对键进行设置操作。 SET key value NX 效果等同于 SETNX key value 。
//XX ：只在键已经存在时，才对键进行设置操作。

var (
	//go:embed script/lua/redis_unlock.lua
	delCmd string
	//go:embed script/lua/redis_refresh.lua
	renewCmd string
	//go:embed script/lua/redis_lock.lua
	lockCmd           string
	defaultExpireTime = 5 //单位：s
)

type RedisLock struct {
	key string
	// 锁的过期时间，单位: s
	expire uint32
	// 锁的标识
	Id string
	// Redis客户端
	redisCli *redis.Client
}

func NewRedisLock(cli *redis.Client, key, id string) *RedisLock {
	return &RedisLock{
		key:      key,
		expire:   uint32(defaultExpireTime),
		Id:       id,
		redisCli: cli,
	}
}

func (r *RedisLock) TryLock() bool {
	result, err := r.redisCli.Eval(context.TODO(), lockCmd, []string{r.key}, r.Id, r.expire).Result()
	if err != nil {
		log.Errorf("tryLock %s %v", r.key, err)
		return false
	}
	i := result.(int64)
	if i == 1 {
		//获取锁成功&自动续期
		go r.reNewExpire()
		return true
	}
	return false
}

func (r *RedisLock) Lock() {
	for {
		if r.TryLock() {
			fmt.Println("redis lock success")
			break
		}
		fmt.Println("redis lock fail")
		time.Sleep(time.Millisecond * 500)
	}
}

func (r *RedisLock) SetExpire(t uint32) {
	r.expire = t
}

func (r *RedisLock) Unlock() {
	//通过lua脚本删除锁
	//1. 查看锁是否存在，如果不存在，直接返回
	//2. 如果存在，对锁进行hincrby -1操作,当减到0时，表明已经unlock完成，可以删除key
	resp, err := r.redisCli.Eval(context.TODO(), delCmd, []string{r.key}, r.Id).Result()
	if err != nil && err != redis.Nil {
		log.Errorf("unlock %s %v", r.key, err)
	}
	if resp == nil {
		fmt.Println("delKey=", resp)
		return
	}
}

// 自动续期
func (r *RedisLock) reNewExpire() {
	ticker := time.NewTicker(time.Duration(r.expire/3) * time.Second)
	for {
		select {
		case <-ticker.C:
			//查看锁是否存在，如果存在进行续期
			resp, err := r.redisCli.Eval(context.TODO(), renewCmd, []string{r.key}, r.Id, r.expire).Result()
			if err != nil && err != redis.Nil {
				log.Errorf("renew key %s err %v", r.key, err)
			}
			if resp.(int64) == 0 {
				return
			}
			log.Infof("renew.....ing...")
		}
	}
}
