package main

import (
	"MoreFun/etcd"
	pb "MoreFun/proto"
	"context"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
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

func (s *MiniGameRouterServer) WriteServers(ctx context.Context, req *pb.WriteRequest) (*pb.WriteResponse, error) {
	leaseID := clientv3.LeaseID(req.LeaseID)
	kvs[req.SvrKey] = req.SvrValue
	kvb[req.SvrKey] = req.IsLease
	kvl[req.SvrKey] = &leaseID
	return nil, nil
}

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
