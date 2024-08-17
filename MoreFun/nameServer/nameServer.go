package main

import (
	"MoreFun/etcd"
	pb "MoreFun/proto"
	"context"
	"flag"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
	"log"
	"sync"
	"time"
)

var (
	kvSlice []kv
	mu      sync.Mutex
)

type kv struct {
	k string
	v string
	b bool
	l *clientv3.LeaseID
}

func main() {
	flag.Parse()
	// 启动Kafka消费者
	go startKafkaConsumer()
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	for {
		mu.Lock()
		for _, kvUnit := range kvSlice {
			var err error
			if kvUnit.b == true {
				//log.Printf("*leaseID %d", *leaseID)
				_, err = cli.Put(context.Background(), kvUnit.k, kvUnit.v, clientv3.WithLease(*kvUnit.l))
			} else {
				//log.Printf("*reallY????")
				_, err = cli.Put(context.Background(), kvUnit.k, kvUnit.v, clientv3.WithIgnoreLease())
			}
			if err != nil {
				log.Printf("提交失败: %v", err)
				continue
			}

			kvSlice = make([]kv, 0)
		}
		mu.Unlock()
		time.Sleep(5 * time.Second)
	}

}

func startKafkaConsumer() {
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
	index := 0
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			// 反序列化消息
			var rRes pb.RegisterServiceResponse
			err := proto.Unmarshal(msg.Value, &rRes)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				continue
			}
			index += 1
			log.Printf("[%d]: receive succss", index)
			mu.Lock()
			leaseID := clientv3.LeaseID(rRes.LeaseID)
			kvNew := kv{rRes.SvrKey, rRes.SvrValue, rRes.IsLease, &leaseID}
			kvSlice = append(kvSlice, kvNew)
			mu.Unlock()
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
