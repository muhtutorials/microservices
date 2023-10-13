package event

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type Publisher struct {
	conn *amqp.Connection
}

func (p *Publisher) setup() error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Publisher) Push(event []byte, severity string) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	log.Println("Pushing to channel")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"logs_topic", // exchange
		severity,     // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        event,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func NewEventPublisher(conn *amqp.Connection) (Publisher, error) {
	publisher := Publisher{
		conn: conn,
	}

	err := publisher.setup()
	if err != nil {
		return Publisher{}, err
	}

	return publisher, nil
}
