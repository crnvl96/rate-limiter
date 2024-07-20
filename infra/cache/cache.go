package cache

import (
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type Cache struct {
	client *redis.Client
}

func NewCache() *Cache {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	return &Cache{
		client: redisClient,
	}
}

func (c *Cache) Get(key string) (string, error) {
	return c.client.Get(key).Result()
}

func (c *Cache) Increment(key string) error {
	return c.client.Incr(key).Err()
}

func (c *Cache) Decrement(key string) (int64, error) {
	return c.client.Decr(key).Result()
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(key).Err()
}

func (c *Cache) ControlExpirationTime(limiterKey string) {
	time.Sleep(1 * time.Second)

	newValue, err := c.Decrement(limiterKey)
	if err != nil {
		log.Printf("Could not decrement counter: %v", err)
		return
	}

	if newValue == 0 {
		if err := c.Delete(limiterKey); err != nil {
			log.Printf("Could not delete key: %v", err)
		}
	}
}

func (c *Cache) Set(key string, value string, expiration time.Duration) error {
	return c.client.Set(key, value, expiration).Err()
}
