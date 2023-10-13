package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"logger/data"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const (
	webPort    = "4000"
	mongoDBURL = "mongodb://mongodb:27017"
	rpcPort    = "9000"
	gRPCPort   = "9001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToDB()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	err = rpc.Register(new(RPCServer))
	go func() {
		err = app.rpcListen()
		if err != nil {
			log.Println(err)
		}
	}()

	go app.gRPCListen()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Println("Starting logger service on port:", webPort)

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToDB() (*mongo.Client, error) {
	opts := options.Client().ApplyURI(mongoDBURL)
	opts.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	conn, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Println("Error connecting to MongoDB:", err)
		return nil, err
	}

	return conn, nil
}

func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port:", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
