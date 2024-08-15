package main

import (
	"MoreFun/etcd"
	pb "MoreFun/proto"
	"context"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"time"
)

// MiniGameRouterServer 实现了 pb.UnimplementedMiniGameRouterServer 接口
type MiniGameRouterServer struct {
	pb.UnimplementedMiniGameRouterServer
}

// 定义命令行参数，用于指定 sidecar 的端口
var (
	ip   = flag.String("ip", "localhost", "The ns ip")
	port = flag.Int("port", 50051, "The ns port")
	kvs  map[string]string
	kvb  map[string]bool
	kvl  map[string]*clientv3.LeaseID
)

func (s *MiniGameRouterServer) CommitService(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	ops := make([]clientv3.Op, 0, len(kvs))
	for k, v := range kvs {
		var op clientv3.Op
		if leaseID, ok := kvl[k]; ok && kvb[k] {
			op = clientv3.OpPut(k, v, clientv3.WithLease(*leaseID))
		} else {
			op = clientv3.OpPut(k, v, clientv3.WithIgnoreLease())
		}
		ops = append(ops, op)
		if len(ops) == 128 {
			if _, err := cli.Txn(ctx).Then(ops...).Commit(); err != nil {
				log.Printf("批量提交失败: %v", err)
				return nil, err
			}
			ops = ops[:0] // 清空 ops 列表
		}
	}
	if len(ops) > 0 {
		if _, err := cli.Txn(context.Background()).Then(ops...).Commit(); err != nil {
			log.Printf("批量失败提交")
		}
	}
	// 清空全局变量
	kvs = make(map[string]string)
	kvb = make(map[string]bool)
	kvl = make(map[string]*clientv3.LeaseID)

	return &pb.CommitResponse{
		Msg: "批量提交成功",
	}, nil
}

func main() {
	kvs = make(map[string]string)
	kvb = make(map[string]bool)
	kvl = make(map[string]*clientv3.LeaseID)
	flag.Parse()

	// 启动Kafka消费者的协程
	go startKafkaConsumer()

	// 启动gRPC服务器
	var lis net.Listener
	var err error
	for i := 0; i < 20; i++ {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Printf("Failed to listen: %v", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}
	if lis != nil {
		defer lis.Close()
	}
	if err != nil {
		return
	}
	s := grpc.NewServer()
	pb.RegisterMiniGameRouterServer(s, &MiniGameRouterServer{})
	log.Printf("Sidecar is listening on port %d", *port)
	if err := s.Serve(lis); err != nil {
		log.Printf("Failed to serve: %v", err)
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
			log.Printf("rece succss")
			leaseID := clientv3.LeaseID(rRes.LeaseID)
			kvs[rRes.SvrKey] = rRes.SvrValue
			kvb[rRes.SvrKey] = rRes.IsLease
			kvl[rRes.SvrKey] = &leaseID
		} else {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
