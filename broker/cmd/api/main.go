package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"os"
	"time"
)

const webPort = "8080"

type Config struct {
	RabbitMQConn *amqp.Connection
}

func main() {
	rabbitMQConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitMQConn.Close()

	app := Config{
		RabbitMQConn: rabbitMQConn,
	}

	log.Println("Starting broker service on port:", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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
