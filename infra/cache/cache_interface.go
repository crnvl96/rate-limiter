package cache

import "time"

type RateLimiterCache interface {
	Get(key string) (string, error)
	Increment(key string) error
	Decrement(key string) (int64, error)
	Delete(key string) error
	ControlExpirationTime(limiterKey string)
	Set(key string, value string, expiration time.Duration) error
}
