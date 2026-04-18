package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func initDB() *pgxpool.Pool {

	connstring := os.Getenv("DATABASE_URL")

	if connstring == "" {
		connstring = "postgres://postgres:pass@localhost:5432/postgres"

	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connstring)

	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Auto create table 👇
	_, err = pool.Exec(ctx,
		` CREATE TABLE IF NOT EXISTS urls ( 
		code TEXT PRIMARY KEY, original_url TEXT, clicks INT DEFAULT 0 
		)
	 `)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create urls table: %v\n", err)
		os.Exit(1)
	}

	_, err = pool.Exec(ctx,
		` CREATE TABLE IF NOT EXISTS users (
		  id SERIAL PRIMARY KEY, email TEXT UNIQUE NOT NULL, password TEXT NOT NULL ) 
	`)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create users table: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("connected to database")
	return pool
}
