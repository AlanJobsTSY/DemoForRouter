package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"strconv"
	"strings"

	"MoreFun/etcd"
	pb "MoreFun/proto"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type MiniGameRouterServer struct {
	pb.UnimplementedMiniGameRouterServer
}

func (s *MiniGameRouterServer) RegisterService(ctx context.Context, req *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	cli := etcd.NewEtcdCli()

	var grantLease bool
	var leaseID clientv3.LeaseID

	serviceKey := fmt.Sprintf("%s_%s", req.Service.Name, req.Service.Port)
	serviceValue := fmt.Sprintf("%s:%s:%s:%s:%s", req.Service.Name, req.Service.Ip, req.Service.Port, req.Service.Protocol, req.Service.Weight)

	// 查看当前的服务是否注册过
	getRes, err := cli.Get(ctx, serviceKey, clientv3.WithCountOnly())
	if err != nil {
		log.Fatalf("Failed to get service key: %v", err)
	}

	// 没有注册过则创建一个租约，同时将租约的ID赋值给leaseID
	if getRes.Count == 0 {
		grantLease = true
		leaseRes, err := cli.Grant(ctx, 10)
		if err != nil {
			log.Fatalf("Failed to grant lease: %v", err)
		}
		leaseID = leaseRes.ID
	}

	// 创建一个kv客户端实现数据插入etcd
	kv := clientv3.NewKV(cli)

	// 开启事务
	txn := kv.Txn(ctx)
	// 判断数据库中是否存在 s *Service
	_, err = txn.If(clientv3.Compare(clientv3.CreateRevision(serviceKey), "=", 0)).
		Then(
			clientv3.OpPut(serviceKey, serviceValue, clientv3.WithLease(leaseID)),
		).
		Else(
			clientv3.OpPut(serviceKey, serviceValue, clientv3.WithIgnoreLease()),
		).
		Commit()

	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	// 租约保活机制
	if grantLease {
		go func() {
			defer cli.Close()

			keepAliveCtx, cancel := context.WithCancel(context.Background())
			defer cancel()

			leaseKeepAlive, err := cli.KeepAlive(keepAliveCtx, leaseID)
			if err != nil {
				log.Fatalf("Failed to keep lease alive: %v", err)
			}
			for lease := range leaseKeepAlive {
				fmt.Printf("leaseID:%x, ttl:%d\n", lease.ID, lease.TTL)
			}
		}()
	}

	return &pb.RegisterServiceResponse{
		Msg: serviceKey + "注册成功",
	}, nil
}

func (s *MiniGameRouterServer) DiscoverService(ctx context.Context, req *pb.DiscoverServiceRequest) (*pb.DiscoverServiceResponse, error) {
	//服务发现
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	//路由选择 needtodo
	getRes, err := cli.Get(context.Background(), req.Msg, clientv3.WithPrefix())
	if err != nil {
		log.Fatalln(err)
	}
	if len(getRes.Kvs) == 0 {
		return nil, fmt.Errorf("no key found for %s", req.Msg)
	}
	value := string(getRes.Kvs[0].Value)
	parts := strings.Split(value, ":")
	addr := fmt.Sprintf("%s:%s", parts[1], parts[2])
	//连接另外一个sidecar
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	c := pb.NewMiniGameRouterClient(conn)

	helloReq := &pb.HelloRequest{
		Msg: "hi",
	}
	rRes, err := c.SayHello(context.Background(), helloReq)
	if err != nil {
		log.Fatalf("could not SayHello: %v", err)
	}

	return &pb.DiscoverServiceResponse{
		Msg: rRes.Msg,
	}, nil

}
func (s *MiniGameRouterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	addr := fmt.Sprintf("localhost:%s", strconv.Itoa(*port-1))
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	c := pb.NewMiniGameRouterClient(conn)
	helloReq := &pb.HelloRequest{
		Msg: fmt.Sprintf("Hello, I am SideCar:%s", *port),
	}
	rRes, err := c.SayHello(context.Background(), helloReq)
	if err != nil {
		log.Fatalf("could not SayHello: %v", err)
	}
	return &pb.HelloResponse{
		Msg: rRes.Msg,
	}, nil
}
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMiniGameRouterServer(s, &MiniGameRouterServer{})
	log.Printf("Server is listening on port %d", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
