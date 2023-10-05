package main

import (
	"context"
	"log"
	"logger/data"
	"time"
)

type RPCserver struct{}
type RPCpayload struct {
	Name string
	Data string
}

func (r *RPCserver) LogInfo(payload RPCpayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error printing from mongo", err)
		return err
	}
	*resp = "process payload by RPC"
	return nil
}
