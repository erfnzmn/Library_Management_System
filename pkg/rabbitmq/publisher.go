package rabbitmq

import (
	"github.com/streadway/amqp"
	"encoding/json"
)
type ReserveMessage struct {
	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`
}

func PublishReservation(ch *amqp.Channel, msg ReserveMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",
		"reserve_requests",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			DeliveryMode: amqp.Persistent,
			Body:        body,
		},
	)
}