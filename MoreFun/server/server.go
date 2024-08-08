package main

import (
	"MoreFun/SDK"
	pb "MoreFun/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

// 定义命令行参数
var (
	name     = flag.String("name", "A", "The server name")
	ip       = flag.String("ip", "localhost", "The server ip")
	port     = flag.Int("port", 50050, "The server port")
	protocol = flag.String("protocol", "grpc", "The server protocol")
	weight   = flag.String("weight", "1", "The server weight")
	status   = flag.String("status", "0", "The server status")
)

// MiniGameRouterServer 实现了 MiniGameRouter gRPC 服务
type MiniGameRouterServer struct {
	pb.UnimplementedMiniGameRouterServer
}

var times int = 0

// SayHello 实现了 gRPC 的 SayHello 方法
func (s *MiniGameRouterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	times += 1
	fmt.Printf("Recv msg: %v times: %d\n", req.Msg, times)
	return &pb.HelloResponse{
		Msg: fmt.Sprintf("Hello, I am %s_%d", *name, *port),
	}, nil
}

// startGRPCServer 启动 gRPC 服务器
func startGRPCServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMiniGameRouterServer(s, &MiniGameRouterServer{})
	log.Printf("Server is listening on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func main() {
	flag.Parse()
	SDK.InitOne(name, ip, port, protocol, weight, status)
	// 启动 gRPC 服务器
	go startGRPCServer(*port)
	conn, client := SDK.InitTwo()
	defer conn.Close()
	SDK.Input(client)

}
