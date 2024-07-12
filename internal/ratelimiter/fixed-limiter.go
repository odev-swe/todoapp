package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	sync.RWMutex
	client map[string]int
	limit  int
	window time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		client: make(map[string]int),
		limit:  limit,
		window: window,
	}
}

func (rl *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	rl.RLock()
	count, exist := rl.client[ip]
	rl.RUnlock()

	if !exist || count < rl.limit {
		rl.Lock()
		if !exist {
			rl.client[ip] = 0
			go rl.Reset(ip)
		}
		rl.client[ip]++
		rl.Unlock()

		return true, rl.window
	}

	return false, rl.window
}

func (rl *FixedWindowLimiter) Reset(clientId string) {
	time.Sleep(rl.window * time.Second)
	rl.Lock()
	defer rl.Unlock()
	delete(rl.client, clientId)
}
