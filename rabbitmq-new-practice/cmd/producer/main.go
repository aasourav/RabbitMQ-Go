package main

import (
	"context"
	"log"
	"time"

	"github.com/aasourav/go-rabbitmq/internal"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := internal.ConnectRabbitMQ("user", "7tFfMZcnRNA0H7yR", "localhost:5672", "customers")

	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)

	if err != nil {
		panic(err)
	}
	defer client.Close()

	if err := client.CreateQueue("customer_created", true, false); err != nil {
		panic(err)
	}

	if err := client.CreateQueue("customers_test", false, true); err != nil {
		panic(err)
	}

	if err := client.CreateExchange("my_exchange", "direct", true); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customer_created", "customer.created.create", "my_exchange"); err != nil {
		panic(err)
	}

	if err := client.CreateBinding("customers_test", "customer.created.test", "my_exchange"); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Send(ctx, "my_exchange", "customer.created.test", amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte("An cool message between services"),
	}); err != nil {
		panic(err)
	}

	if err := client.Send(ctx, "my_exchange", "customer.created.create", amqp091.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp091.Persistent,
		Body:         []byte("An cool message between services"),
	}); err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	log.Println(client)
}
