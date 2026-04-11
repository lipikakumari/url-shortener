package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimiter(rdb *redis.Client) gin.HandlerFunc {

	return func(c *gin.Context) {
		ip := c.ClientIP()
		ctx := context.Background()

		count, err := rdb.Incr(ctx, ip).Result()

		if err != nil {
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}

		if count == 1 {
			rdb.Expire(ctx, ip, time.Minute)
		}

		if count > 10 {
			c.JSON(429, gin.H{"error": "to many attemots"})
			c.Abort()
			return

		}

		c.Next()
	}

}
