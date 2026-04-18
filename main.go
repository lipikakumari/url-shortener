package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func main() {
	db := initDB()
	defer db.Close()

	rdb := initRedis()
	defer rdb.Close()

	h := Handler{
		db:    db,
		redis: rdb,
	}

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Gin is working"})
	})

	r.Use(RateLimiter(h.redis))

	r.POST("/shorten", JWTMiddleware(), h.shortenURL)
	r.GET("/:code", h.redirectURL)
	r.POST("/register", h.register)
	r.POST("/login", h.login)

	r.Run(":8080")
}

func randomCode(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (h *Handler) shortenURL(c *gin.Context) {

	userId, _ := c.Get("userId")
	fmt.Println("Request from user:", userId)

	var req struct {
		URL string `json:"url"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	shortCode := randomCode(6)

	err := InsertURL(h.db, shortCode, req.URL)

	if err != nil {
		fmt.Println("insert error", err)
		c.JSON(500, gin.H{"error": "failed to save URL"})
		return
	}

	c.JSON(200, gin.H{"short_url": "http://localhost:8080/" + shortCode})

	CacheURL(h.redis, shortCode, req.URL)

}

func (h *Handler) redirectURL(c *gin.Context) {

	code := c.Param("code")

	fmt.Println(code)

	// check in Redis first

	cached, err := GetCachedURL(h.redis, code)

	if err == nil {
		fmt.Println("cache hit")
		Publish(code)
		c.Redirect(302, cached)
		return
	}

	originalURL, err := GetURL(h.db, code)

	if err != nil {
		c.JSON(404, gin.H{
			"error": "short URL not found",
		})
		return
	}

	CacheURL(h.redis, code, originalURL)

	Publish(code)

	c.Redirect(http.StatusFound, originalURL)

}

func (h *Handler) register(c *gin.Context) {

	var req struct {
		Email string `json:"email"`
		Pwd   string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	hashPassword, err := HashPassword(req.Pwd)

	if err != nil {
		fmt.Println("error while hashing password", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	err = CreateUser(h.db, req.Email, hashPassword)

	if err != nil {
		fmt.Println("insert error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created"})

}

func (h *Handler) login(c *gin.Context) {

	var req struct {
		Email string `json:"email"`
		Pwd   string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON",
		})
		return
	}

	id, pwd, err := GetUserByEmail(h.db, req.Email)

	if err != nil {
		fmt.Println("error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user"})
		return
	}

	err = CheckPassword(req.Pwd, []byte(pwd))

	if err != nil {
		fmt.Println("password mismatch", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password do not match"})
		return
	}

	tokenString, err := GenerateToken(id)

	if err != nil {
		fmt.Println(" token creation error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "token not created"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}
