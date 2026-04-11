package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ClickEvent struct {
	Code string `json:"code"`
}

func Publish(code string) {

	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	fmt.Println("api connected to RabbitMQ!")

	ch, err := conn.Channel()

	if err != nil {
		log.Println("Failed to open channel", err)
		return
	}

	defer conn.Close()

	qu, err := ch.QueueDeclare(
		"clicks",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("Failed to declare queue", err)
		return
	}

	event := ClickEvent{Code: code}
	body, err := json.Marshal(event)

	if err != nil {
		log.Println("Failed to marshall message", err)
		return
	}

	err = ch.PublishWithContext(
		context.Background(),
		"",
		qu.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println("Failed to publish message", err)
		return
	}

	fmt.Println("published click event for code:", code)

}
