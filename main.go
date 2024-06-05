package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Go Rabbitmq")
	conn, err := amqp.Dial("amqp://user:02C3dgdhmmIneepB@127.0.0.1:5672/")

	if err != nil {
		log.Println("ERR: %v", err.Error())
	} else {

		log.Println("Successfully Connected to our RabbitMQ")
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println("ERR: %v", err.Error())
	}

	queue, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Println("ERR: %v", err.Error())
	}

	defer ch.Close()
	defer conn.Close()

	fmt.Println(queue)

	err = ch.Publish("", "TestQueue", false, false, amqp.Publishing{ContentType: "text/plain", Body: []byte("Hello World")})
	if err != nil {
		log.Println("ERR: %v", err.Error())
	}

	fmt.Println("Publish Success")
}
