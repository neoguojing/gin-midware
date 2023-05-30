package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	lru "github.com/hashicorp/golang-lru/v2"
	"golang.org/x/time/rate"
)

type KeyFuc func(*gin.Context) string

var (
	cache, _ = lru.New[string, *rate.Limiter](128)
)

// RateLimiter is a middleware that limits the request rate based on a given key
func GinRateLimiter(key KeyFuc, maxToken int, interval time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		k := key(c)
		v, ok := cache.Get(k)
		if !ok {
			v = rate.NewLimiter(rate.Every(interval), maxToken)
			cache.Add(k, v)
		} else {
			if v.Allow() {
				c.Next()
			} else {
				c.AbortWithStatus(http.StatusTooManyRequests)
			}
		}
	}
}

func RateLimiter(first, every int, interval time.Duration, f func()) {
	s := rate.Sometimes{
		First:    first,
		Every:    every,
		Interval: interval,
	}
	s.Do(f)
}
