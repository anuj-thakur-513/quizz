package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/anuj-thakur-513/quizz/pkg/core"
	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	requests    map[string]*requestInfo
	mutex       sync.Mutex
	maxRequests int
	duration    time.Duration
}

type requestInfo struct {
	count     int
	timestamp time.Time
}

func NewRateLimiter(maxRequests int, duration time.Duration) *rateLimiter {
	return &rateLimiter{
		requests:    make(map[string]*requestInfo),
		maxRequests: maxRequests,
		duration:    duration,
	}

}

func (rl *rateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mutex.Lock()
		defer rl.mutex.Unlock()

		if info, exists := rl.requests[ip]; exists {
			if time.Since(info.timestamp) > rl.duration {
				// Reset counter if the time window has passed
				info.count = 1
				info.timestamp = time.Now()
			} else if info.count >= rl.maxRequests {
				c.JSON(http.StatusTooManyRequests, core.NewAppError(400, "Too many requests"))
				c.Abort()
				return
			} else {
				info.count++
			}
		} else {
			// New IP entry
			rl.requests[ip] = &requestInfo{
				count:     1,
				timestamp: time.Now(),
			}
		}

		c.Next()
	}
}
