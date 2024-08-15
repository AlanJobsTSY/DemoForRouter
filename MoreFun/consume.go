package main

import (
	pb "MoreFun/proto"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
)

func main() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer c.Close()

	c.SubscribeTopics([]string{"myTopic"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			// 反序列化消息
			var writeReq pb.WriteRequest
			err := proto.Unmarshal(msg.Value, &writeReq)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				continue
			}

			// 打印消息内容
			fmt.Printf("Received message: %+v\n", writeReq)
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
