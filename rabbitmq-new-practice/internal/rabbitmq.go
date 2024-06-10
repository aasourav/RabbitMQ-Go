package internal

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	// the connection used by the client
	// we will reuse the connection to our whole application
	connection *amqp.Connection

	// channel is used to process / Send messages
	channel *amqp.Channel
	/**
		Channel is a multiplexed (allows multiple applications to send and receive data simultaneously) sub connection
	**/
}

func ConnectRabbitMQ(username, pass, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, pass, host, vhost))
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		connection: conn,
		channel:    ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.channel.Close()
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autoDelete bool) error {
	_, err := rc.channel.QueueDeclare(queueName, durable, autoDelete, false, false, amqp.Table{})
	return err
}

func (rc RabbitClient) CreateExchange(exchangeName, kind string, durable bool) error {
	return rc.channel.ExchangeDeclare(exchangeName, kind, durable, false, false, false, amqp.Table{})
}

func (rc RabbitClient) CreateBinding(name, binding, exhange string) error {
	return rc.channel.QueueBind(name, binding, exhange, false, amqp.Table{})
}

func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	// return rc.channel.Publish
	return rc.channel.PublishWithContext(ctx, exchange, routingKey,
		// mandatory is used to determine if an error should be returned upo failure
		true,
		// immidiate
		false, options)
}

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.channel.Consume(queue, consumer, autoAck,
		// exclusive true means this will be one and only consumer consume that queue
		// if false , server will distribute the message using loadbalancer technique
		autoAck,
		//publishing and consuming from the same domain
		false, false, amqp.Table{},
	)
}
