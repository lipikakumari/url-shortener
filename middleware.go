package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func tooManyAttempts(rdb *redis.Client, ip string, limit int64) bool {

	ctx := context.Background()

	count, err := rdb.Incr(ctx, ip).Result()

	if err != nil {
		return false
	}

	if count == 1 {
		rdb.Expire(ctx, ip, time.Minute)
	}

	return count > 10

}

func RateLimiter(rdb *redis.Client) gin.HandlerFunc {

	return func(c *gin.Context) {
		ip := c.ClientIP()

		tooManyAttempts := tooManyAttempts(rdb, ip, 10)

		if tooManyAttempts {
			c.JSON(429, gin.H{"error": "to many attemots"})
			c.Abort()
			return

		}

		c.Next()
	}

}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1 - get token from header
		authHeader := c.GetHeader("Authorization")
		fmt.Println("Auth header:", authHeader)

		if authHeader == "" {
			c.JSON(401, gin.H{"error": "no token provided"})
			c.Abort()
			return
		}

		// Step 2 - remove "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		fmt.Println("tokenString:", tokenString)

		// Step 3 - validate token
		userId, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Step 4 - save userId for use in handler
		c.Set("userId", userId)

		c.Next()
	}
}
