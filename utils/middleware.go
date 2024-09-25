package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	rateLimit    = 200
	requestCount = make(map[string]int)
	mutex        sync.Mutex
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mutex.Lock()
		defer mutex.Unlock()

		requestCount[ip]++

		if requestCount[ip] > rateLimit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status": "error",
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		go func() {
			time.Sleep(time.Minute)
			mutex.Lock()
			defer mutex.Unlock()
			requestCount[ip] = 0
		}()

		c.Next()
	}
}

func LowercaseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.URL.Path = strings.ToLower(c.Request.URL.Path)
		c.Next()
	}
}