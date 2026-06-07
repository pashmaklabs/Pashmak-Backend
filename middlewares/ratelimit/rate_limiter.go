package middlewares_ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type KeyFunc func(c *gin.Context) string

// RateLimiter holds config for a sliding-window rate limiter.
type RateLimiter struct {
	redis    *redis.Client
	limit    int           // max requests
	window   time.Duration // per window
	keyFunc  KeyFunc
	keyPrefix string
}

func NewRateLimiter(redisClient *redis.Client, limit int, window time.Duration, keyPrefix string, keyFunc KeyFunc) *RateLimiter {
	if keyFunc == nil {
		keyFunc = KeyByIP
	}
	return &RateLimiter{
		redis:     redisClient,
		limit:     limit,
		window:    window,
		keyFunc:   keyFunc,
		keyPrefix: keyPrefix,
	}
}

// Common key strategies
func KeyByIP(c *gin.Context) string {
	return c.ClientIP()
}

func KeyByEmail(c *gin.Context) string {
	// Reads email from validated body set by ValidationMiddleware
	if v, exists := c.Get("validated"); exists {
		type hasEmail interface{ GetEmail() string }
		if e, ok := v.(hasEmail); ok {
			return e.GetEmail()
		}
	}
	return c.ClientIP() // fallback
}

// Middleware returns a gin.HandlerFunc using a sliding window (INCR + EXPIRE).
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		key := fmt.Sprintf("rl:%s:%s", rl.keyPrefix, rl.keyFunc(c))

		// Sliding window via INCR + TTL set only on first request
		count, err := rl.redis.Incr(ctx, key).Result()
		if err != nil {
			// Fail open: don't block users if Redis is down
			c.Next()
			return
		}

		if count == 1 {
			rl.redis.Expire(ctx, key, rl.window)
		}

		// Set informational headers
		ttl, _ := rl.redis.TTL(ctx, key).Result()
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, rl.limit-int(count))))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(ttl).Unix()))

		if int(count) > rl.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "درخواست‌های شما بیش از حد مجاز است. لطفاً کمی صبر کنید.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}