package main

import (
	"log"

	"github.com/aasourav/go-rabbitmq/internal"
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

	messageBus, err := client.Consume("customer_created", "consumer-name", false)
	if err != nil {
		panic(err)
	}

	var blocking chan struct{}

	go func() {
		for message := range messageBus {
			str := string(message.Body)
			log.Println("New Message: ", str)
			// message.Nack() // requeue the msg
			if !message.Redelivered {
				message.Nack(false, true)
				continue
			}
			if err := message.Ack(true); err != nil {
				log.Println(err.Error())
			}

		}
	}()

	log.Println("Consuming, to close the program press CTRL + C")
	<-blocking
}
