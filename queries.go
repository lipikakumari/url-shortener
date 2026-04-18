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

func CreateUser(db *pgxpool.Pool, email string, hashPassword []byte) error {
	ctx := context.Background()

	_, err := db.Exec(ctx,
		"INSERT INTO users (email , password) values ($1 , $2)",
		email, hashPassword,
	)

	return err
}

func GetUserByEmail(db *pgxpool.Pool, email string) (int32, string, error) {
	ctx := context.Background()

	var id int32
	var pwd string

	err := db.QueryRow(ctx,
		"SELECT id , password FROM users WHERE email = $1",
		email,
	).Scan(&id, &pwd)

	if err != nil {
		return 0, "", err
	}

	return id, pwd, nil
}
