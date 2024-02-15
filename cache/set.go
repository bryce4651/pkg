package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type SetCache struct {
	cli *redis.Client
	key string
}

func NewSetCache(name string, cli *redis.Client) *SetCache {
	return &SetCache{
		cli: cli,
		key: name,
	}
}

func (c *SetCache) Get(ctx context.Context, _ string) (string, error) {
	val, err := c.cli.SPop(ctx, c.key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func (c *SetCache) Set(ctx context.Context, member string, _ interface{}) error {
	return c.cli.SAdd(ctx, c.key, member).Err()
}

func (c *SetCache) Del(ctx context.Context, members ...string) error {
	var ms []interface{}
	for _, m := range members {
		ms = append(ms, m)
	}
	return c.cli.SRem(ctx, c.key, ms...).Err()
}
func (c *SetCache) Exist(ctx context.Context, member string) (bool, error) {
	return c.cli.SIsMember(ctx, c.key, member).Result()
}
