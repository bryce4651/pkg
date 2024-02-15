package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type HashCache struct {
	cli *redis.Client
	//expire   time.Duration
	hashName string
}

func NewHashCache(name string, cli *redis.Client) *HashCache {
	return &HashCache{
		cli:      cli,
		hashName: name,
	}
}

//func (c *HashCache) getRandExpire() time.Duration {
//	n := rand.Intn(100)
//	e := c.expire / 10
//	r := e * time.Duration(n/100)
//	return c.expire + r
//}

func (c *HashCache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.cli.HGet(ctx, c.hashName, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (c *HashCache) Set(ctx context.Context, key string, value interface{}) error {
	return c.cli.HSet(ctx, c.hashName, key, value).Err()
}
func (c *HashCache) Del(ctx context.Context, keys ...string) error {
	return c.cli.HDel(ctx, c.hashName, keys...).Err()
}
func (c *HashCache) Exist(ctx context.Context, key string) (bool, error) {
	return c.cli.HExists(ctx, c.hashName, key).Result()
}
