package cache

import (
	"context"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

type StrCache struct {
	cli    *redis.Client
	expire time.Duration
}

func NewStrCache(cli *redis.Client, expire time.Duration) *StrCache {
	return &StrCache{
		cli:    cli,
		expire: expire,
	}
}

func (c *StrCache) getRandExpire() time.Duration {
	n := rand.Intn(100)
	e := c.expire / 10
	r := e * time.Duration(n/100)
	return c.expire + r
}

func (c *StrCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.cli.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (c *StrCache) Set(ctx context.Context, key string, value interface{}) error {
	return c.cli.Set(ctx, key, value, c.getRandExpire()).Err()
}
func (c *StrCache) SetEx(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return c.cli.Set(ctx, key, value, expire).Err()
}
func (c *StrCache) Del(ctx context.Context, keys ...string) error {
	return c.cli.Del(ctx, keys...).Err()
}
func (c *StrCache) Exist(ctx context.Context, key string) (bool, error) {
	res, err := c.cli.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if res > 0 {
		return true, nil
	}
	return false, nil
}
