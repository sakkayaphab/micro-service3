package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"strconv"
	"time"
)

type Message struct {
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	MsgId int64 `json:"Msg_id"`
	Sender string `json:"Sender"`
	Msg string `json:"Msg"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:rootpassword@mongodb_container:27017/?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&ssl=false"))
	if err!=nil {
		log.Fatal(err)
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "message",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	err = r.SetOffset(-2)
	if err!=nil{
		log.Fatal(err)
	}

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		var mm Message
		err = json.Unmarshal(m.Value, &mm)
		if err!=nil {
			log.Println(err)
		} else {
			fmt.Println(strconv.FormatInt(mm.MsgId, 10)+", "+mm.Sender+": "+mm.Msg+", "+time.Now().Format("2006-01-02T15:04:05.000Z"))
			collection := client.Database("mydatabase").Collection("messages")
			mm.ID = primitive.NewObjectID()
			_, err = collection.InsertOne(ctx, mm)
			if err!=nil {
				log.Println(err)
			}
		}
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	return

}
