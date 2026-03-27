package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	ctx    context.Context
}

func NewCache(addr string) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Cache{
		client: client,
		ctx:    context.Background(),
	}
}

func (c *Cache) Get(key string, dest any) error {
	val, err := c.client.Get(c.ctx, key).Result()
	if err == redis.Nil {
		return nil // key does not exist
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *Cache) Set(key string, value any, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(c.ctx, key, jsonData, expiration).Err()
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

func (c *Cache) HealthCheck() error {
	return c.client.Ping(c.ctx).Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}
