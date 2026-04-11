package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func GetCachedURL(rds *redis.Client, code string) (string, error) {
	ctx := context.Background()

	var originalURL string

	originalURL, err := rds.Get(ctx, code).Result()

	return originalURL, err
}

func CacheURL(rds *redis.Client, code string, url string) error {
	ctx := context.Background()

	err := rds.Set(ctx, code, url, 24*time.Hour).Err()

	return err
}

func InsertURL(db *pgxpool.Pool, code string, url string) error {
	ctx := context.Background()

	_, err := db.Exec(ctx,
		"INSERT INTO urls (code , original_url) values ($1 , $2)",
		code, url,
	)

	return err
}

func GetURL(db *pgxpool.Pool, code string) (string, error) {
	ctx := context.Background()

	var originalURL string

	err := db.QueryRow(ctx,
		"SELECT original_url FROM urls WHERE code = $1",
		code,
	).Scan(&originalURL)

	if err != nil {
		return "", err
	}

	return originalURL, nil
}
