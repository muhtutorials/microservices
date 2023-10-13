package main

import (
	"context"
	"log"
	"logger/data"
	"time"
)

type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

func (receiver *RPCServer) LogInfo(payload RPCPayload, res *string) error {
	col := client.Database("logger").Collection("logs")
	_, err := col.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to MongoDB:", err)
		return err
	}

	*res = "Processed payload via RPC " + payload.Name

	return nil
}
