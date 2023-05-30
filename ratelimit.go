package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type KeyFuc func(*gin.Context) string

var (
	keyMap = sync.Map{}
)

// RateLimiter is a middleware that limits the request rate based on a given key
func GinRateLimiter(key KeyFuc, maxToken int, interval time.Duration) gin.HandlerFunc {

	return func(c *gin.Context) {
		k := key(c)
		v, ok := keyMap.Load(k)
		if !ok {
			v = rate.NewLimiter(rate.Every(interval), maxToken)
			keyMap.Store(k, v)
		} else {
			limiter := v.(*rate.Limiter)
			if limiter.Allow() {
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
