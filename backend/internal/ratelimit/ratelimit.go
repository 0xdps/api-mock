package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// Limiter is a sliding-window rate limiter backed by Redis.
// When Redis is unavailable it fails open (all requests pass).
type Limiter struct {
	rdb *redis.Client
}

// New creates a Limiter. rdb may be nil (fail-open).
func New(rdb *redis.Client) *Limiter {
	return &Limiter{rdb: rdb}
}

// Allow checks whether the given key is within the RPM limit.
// Returns (allowed, currentCount, error).
func (l *Limiter) Allow(ctx context.Context, key string, rpm int) (bool, int, error) {
	if l.rdb == nil {
		return true, 0, nil
	}

	window := time.Now().UTC().Truncate(time.Minute).Unix()
	redisKey := fmt.Sprintf("rl:%s:%d", key, window)

	pipe := l.rdb.Pipeline()
	incr := pipe.Incr(ctx, redisKey)
	pipe.Expire(ctx, redisKey, 2*time.Minute) // keep a little beyond the window
	if _, err := pipe.Exec(ctx); err != nil {
		return true, 0, nil // fail open on Redis error
	}

	count := int(incr.Val())
	return count <= rpm, count, nil
}

// Middleware returns an http.Handler middleware that rate-limits by a dynamic
// key resolved from the request. If keyFn returns "" the limiter is skipped.
func (l *Limiter) Middleware(rpm int, keyFn func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFn(r)
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			allowed, count, _ := l.Allow(r.Context(), key, rpm)
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rpm))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(max(0, rpm-count)))

			if !allowed {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Retry-After", "60")
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprintf(w, `{"error":"rate limit exceeded","limit":%d,"retry_after":60}`, rpm)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
