package main

import (
	pb "MoreFun/proto"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
)

func main() {
	// 创建一个新的Kafka生产者实例，配置Kafka服务器地址
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close() // 在函数结束时关闭生产者

	topic := "myTopic" // 定义要发送消息的Kafka主题

	// 创建WriteRequest消息
	writeReq := &pb.WriteRequest{
		SvrKey:   "serviceKey",   // 设置服务键
		SvrValue: "serviceValue", // 设置服务值
		IsLease:  true,           // 设置是否租约
		LeaseID:  12345,          // 设置租约ID
	}

	// 序列化消息，将WriteRequest对象转换为字节数组
	msgBytes, err := proto.Marshal(writeReq)
	if err != nil {
		log.Fatalf("Failed to marshal message: %s", err)
	}

	// 创建Kafka消息
	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, // 设置消息的主题和分区
		Value:          msgBytes,                                                           // 设置消息的内容
	}

	// 发送消息到Kafka
	err = p.Produce(&msg, nil)
	if err != nil {
		log.Fatalf("Failed to produce message: %s", err)
	}

	// 等待消息传递完成，最多等待15秒
	p.Flush(30 * 1000)
	fmt.Println("Message sent successfully")
}
