package rate

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Lua script برای پیاده‌سازی اتمیک Token Bucket در Redis
// State در Hash نگه‌داری می‌شود: fields => "tokens" (float) و "last_ms" (timestamp ms)
const tokenBucketLua = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local refill_rate = tonumber(ARGV[2]) -- tokens per second
local now_ms = tonumber(ARGV[3])
local cleanup_ttl_ms = tonumber(ARGV[4])

local data = redis.call('HMGET', key, 'tokens', 'last_ms')
local tokens = tonumber(data[1])
local last = tonumber(data[2])

if tokens == nil then
  tokens = capacity
  last = now_ms
else
  local delta = now_ms - last
  if delta > 0 then
    tokens = math.min(capacity, tokens + (delta / 1000.0) * refill_rate)
    last = now_ms
  end
end

local allowed = 0
local retry_after_ms = 0

if tokens >= 1.0 then
  tokens = tokens - 1.0
  allowed = 1
else
  allowed = 0
  retry_after_ms = math.ceil((1.0 - tokens) / refill_rate * 1000.0)
end

redis.call('HMSET', key, 'tokens', tokens, 'last_ms', last)
redis.call('PEXPIRE', key, cleanup_ttl_ms)

return {allowed, tostring(tokens), retry_after_ms}
`

type Limiter struct {
	rdb        *redis.Client
	capacity   float64       // ظرفیت سطل (حداکثر توکن‌ها)
	refillRate float64       // نرخ شارژ (tokens per second)
	cleanupTTL time.Duration // TTL برای پاک‌سازی خودکار bucketهای بلااستفاده
	script     *redis.Script
}

// NewTokenBucket سازندهٔ Limiter با مدل Token Bucket
// مثال استفاده: NewTokenBucket(rdb, 5, 1, 2*time.Minute, 20*time.Minute)
// یعنی: ظرفیت 5، هر 2 دقیقه 1 توکن شارژ، و TTL پاکسازی 20 دقیقه
func NewTokenBucket(rdb *redis.Client, capacity int, refillTokens int, refillInterval time.Duration, cleanupTTL time.Duration) *Limiter {
	rate := float64(refillTokens) / refillInterval.Seconds() // tokens per second
	return &Limiter{
		rdb:        rdb,
		capacity:   float64(capacity),
		refillRate: rate,
		cleanupTTL: cleanupTTL,
		script:     redis.NewScript(tokenBucketLua),
	}
}

// TooMany: سازگار با Handler فعلی شما (True یعنی بلاک شده‌ای)
// در Token Bucket: هر بار صدا خوردن، اگر توکن باشد یکی کم می‌شود؛ اگر نباشد => بلاک + زمان انتظار
func (l *Limiter) TooMany(ctx context.Context, key string) (blocked bool, retryAfterSec int64, err error) {
	args := []interface{}{
		strconv.FormatFloat(l.capacity, 'f', -1, 64),
		strconv.FormatFloat(l.refillRate, 'f', -1, 64),
		time.Now().UnixMilli(),
		l.cleanupTTL.Milliseconds(),
	}
	res, err := l.script.Run(ctx, l.rdb, []string{key}, args...).Result()
	if err != nil {
		return false, 0, err
	}

	vals := res.([]interface{})
	allowed := vals[0].(int64)       // 1 or 0
	retryMs := vals[2].(int64)       // زمان انتظار تا توکن بعدی (ms)
	if allowed == 1 {
		return false, 0, nil // بلاک نیست
	}
	// تبدیل به ثانیه با ceiling
	sec := (retryMs + 999) / 1000
	return true, sec, nil
}

// Reset: اگر جایی خواستی bucket را دستی پاک کنی
func (l *Limiter) Reset(ctx context.Context, key string) error {
	return l.rdb.Del(ctx, key).Err()
}
