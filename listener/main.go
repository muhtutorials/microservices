package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener/event"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	rabbitMQConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitMQConn.Close()

	log.Println("Listening for and consuming RabbitMQ messages...")

	consumer, err := event.NewConsumer(rabbitMQConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.info", "log.warning", "log.error"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var count int64
	var wait = time.Second
	var conn *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not ready yet")
			count++
		} else {
			log.Println("Connected to RabbitMQ")
			conn = c
			break
		}

		if count > 5 {
			fmt.Println(err)
			return nil, err
		}

		wait = time.Duration(math.Pow(float64(count), 2)) * time.Second
		log.Println("Waiting...")
		time.Sleep(wait)
	}

	return conn, nil
}
