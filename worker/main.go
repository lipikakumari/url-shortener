package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ClickEvent struct {
	Code string `json:"code"`
}

func main() {
	// Step 1 - connect to database
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		connString = "postgres://postgres:pass@localhost:5432/postgres"
	}

	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()
	fmt.Println("Worker connected to database!")

	// Step 2 - connect to RabbitMQ

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitmqURL)
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()
	fmt.Println("Worker connected to RabbitMQ!")

	// Step 3 - open a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}
	defer ch.Close()

	// Step 4 - declare the queue
	q, err := ch.QueueDeclare(
		"clicks",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to declare queue:", err)
	}

	// Step 5 - start consuming
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal("Failed to consume:", err)
	}

	fmt.Println("Worker is listening for click events...")

	// Step 6 - process messages forever
	for msg := range msgs {
		var event ClickEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			log.Println("Failed to unmarshal message:", err)
			continue
		}

		fmt.Println("Processing click for code:", event.Code)

		_, err = db.Exec(context.Background(),
			"UPDATE urls SET clicks = clicks + 1 WHERE code = $1",
			event.Code,
		)
		if err != nil {
			log.Println("Failed to increment clicks:", err)
		}
	}
}
