package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/erfnzmn/Library_Management_System/internal/loans"
)

func ConsumeReservations(ch *amqp.Channel, loanService *loans.Service) error {
	msgs, err := ch.Consume(
		"reserve_requests", // queue name
		"loans-worker-1", // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil, // args
	)
	if err != nil { return err }

	go func() {
		for d := range msgs {
			var req ReserveMessage
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("invalid message: %v", err)
				_ = d.Nack(false, false) 
				continue
			}

			if err := loanService.ReserveBook(context.Background(), req.UserID, req.BookID); err != nil {
				log.Printf("reserve failed (user=%d book=%d): %v", req.UserID, req.BookID, err)
				_ = d.Nack(false, false)
			} else {
				_ = d.Ack(false)
			}
		}
	}()

	return nil
}
