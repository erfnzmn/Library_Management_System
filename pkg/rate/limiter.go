package rate

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	rdb    *redis.Client
	limit  int
	window time.Duration
}

func New(rdb *redis.Client, limit int, window time.Duration) *Limiter {
	return &Limiter{
		rdb:   rdb,
		limit: limit,
		window: window,
	}
}

func (l *Limiter) TooMany(ctx context.Context, key string) (blocked bool, retryAfterSec int64, err error) {
	if l.rdb == nil {
		return false, 0, nil 
	}

	pipe := l.rdb.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, l.window) 
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	n := incr.Val()
	if int(n) > l.limit {
		ttl, _ := l.rdb.TTL(ctx, key).Result()
		if ttl < 0 {
			ttl = l.window
		}
		return true, int64(ttl.Seconds()), nil
	}
	return false, 0, nil
}


