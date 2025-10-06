package redisclient 

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Enabled bool
	Addr    string
	Password string
	DB      int
}

func New(cfg Config) (*redis.Client, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		DialTimeout: 3 * time.Second,
		ReadTimeout: 2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}