package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestNewExpiring(t *testing.T) {
	cache := NewExpiring()

	// 模拟用户登录，并设置会话过期时间为10分钟
	sessionInfo, err := handleUserLogin("john.doe", "password123", cache, 10*time.Second)
	if err != nil {
		fmt.Println("登录失败:", err)
		return
	}
	fmt.Println("登录成功，会话信息:", sessionInfo)

	// 模拟用户请求，尝试从缓存中获取会话信息
	session, found := getUserSession(sessionInfo.SessionID, cache)
	if found {
		fmt.Println("获取会话信息成功:", session)
	} else {
		fmt.Println("未找到会话信息或会话已过期")
	}

	// 等待一段时间以模拟会话过期，然后再次尝试获取会话信息
	time.Sleep(15 * time.Second) // 假设现在过去了15秒
	session, found = getUserSession(sessionInfo.SessionID, cache)
	if found {
		fmt.Println("获取会话信息成功（但理论上应该已过期）:", session)
	} else {
		fmt.Println("未找到会话信息或会话已过期（符合预期）")
	}
}

// SessionInfo 表示用户的会话信息
type SessionInfo struct {
	SessionID string
	UserInfo  string // 示例，实际中可以包含更多用户信息
}

// handleUserLogin 模拟用户登录过程，并在登录成功后将会话信息存储到缓存中
func handleUserLogin(username, password string, cache *Expiring, expiration time.Duration) (*SessionInfo, error) {
	// 这里应该包含实际的用户验证逻辑，但为了示例简单，我们假设验证总是成功
	// 创建一个SessionInfo实例（实际中可能包含更多逻辑）
	sessionInfo := &SessionInfo{
		SessionID: "unique-session-id-" + time.Now().Format("20060102150405"),
		UserInfo:  username, // 假设用户名就是用户信息（实际中会更复杂）
	}

	// 将会话信息存储到缓存中，并设置过期时间
	cache.Set(sessionInfo.SessionID, sessionInfo, expiration)
	return sessionInfo, nil
}

// getUserSession 尝试从缓存中获取指定会话ID的会话信息
func getUserSession(sessionID string, cache *Expiring) (*SessionInfo, bool) {
	value, found := cache.Get(sessionID)
	if !found {
		return nil, false
	}
	sessionInfo, ok := value.(*SessionInfo)
	if !ok {
		// 缓存中的数据类型不匹配，可能是数据已损坏或类型被错误地存储
		return nil, false
	}
	return sessionInfo, true
}

func BenchmarkNewExpiring(b *testing.B) {
	b.Run("Run.N", func(b *testing.B) {
		b.ReportAllocs()
		c := NewExpiring()
		for i := 0; i < b.N; i++ {
			bs, ok := c.Get(i)
			b.Log(bs)
			if !ok {
				c.Set(i, i, time.Second*2)
			}
			c.Delete(i)
		}
	})
	b.Run("Run.Parallel", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			c := NewExpiring()
			for pb.Next() {
				i := 1024
				bs, ok := c.Get(i)
				b.Log(bs)
				if !ok {
					c.Set(i, i, time.Second*2)
				}
				c.Delete(i)
			}
		})
	})
}
