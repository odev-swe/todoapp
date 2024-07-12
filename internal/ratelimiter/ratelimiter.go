package ratelimiter

import "time"

type RateLimiter interface {
	Allow(ip string) (bool, time.Duration)
}
