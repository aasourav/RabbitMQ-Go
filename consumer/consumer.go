package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Consumer")
	conn, err := amqp.Dial("amqp://user:02C3dgdhmmIneepB@127.0.0.1:5672/")
	if err != nil {
		log.Println("ERR: %v", err.Error())
	} else {

		log.Println("Successfully Connected to our RabbitMQ")
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("ERR: %v", err.Error())
	}

	defer ch.Close()

	msgs, err := ch.Consume("TestQueue", "", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			val := string(d.Body)
			fmt.Println("Received message:", val)
		}
	}()

	fmt.Println("[*] - Waiting for messages")
	<-forever
}
