package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type CacheCfg struct {
	ClientName string `json:"client_name" yaml:"client_name" env:"name=CACHE_CLIENT_NAME"`
	Address    string `json:"address" yaml:"address" env:"name=CACHE_ADDRESS"`
	Username   string `json:"username" yaml:"username" env:"name=CACHE_USERNAME"`
	Password   string `json:"password" yaml:"password" env:"name=CACHE_PASSWORD"`
	DB         int    `json:"db" yaml:"db" env:"name=CACHE_DB"`
}

func GetCacheCli(config *CacheCfg) (*redis.Client, error) {
	if config == nil {
		return nil, fmt.Errorf("cache config is nil")
	}
	return redis.NewClient(&redis.Options{
		Addr:       config.Address,
		DB:         config.DB,
		ClientName: config.ClientName,
		Username:   config.Username,
		Password:   config.Password,
	}), nil
}

type ICache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}) error
	Del(ctx context.Context, keys ...string) error
	Exist(ctx context.Context, key string) (bool, error)
}
