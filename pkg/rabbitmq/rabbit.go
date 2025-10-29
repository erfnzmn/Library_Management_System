package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
)

type Client struct{
	Conn *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil{
		conn.Close()
		return nil, err
	}
	_, err = ch.QueueDeclare(
		"reserve_requests",
		true, // durable
		false, // auto-delete
		true, // exclusive
		false, // no-wait
		nil, //args
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	log.Println("🐇 RabbitMQ connected and queue ready")
	return &Client{Conn: conn, Channel: ch}, nil
}

func (r *Client) Close() {
	if r.Channel != nil {
		_ = r.Channel.Close()
	}
	if r.Conn != nil {
		_ = r.Conn.Close()
	}
}
