package main

import (
	"MoreFun/SDK"
	"MoreFun/endPoint"
	pb "MoreFun/proto"
	"context"
	"flag"
	"fmt"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
	"sync"
)

// 定义命令行参数
var (
	name     = flag.String("name", "A", "The server name")
	ip       = flag.String("ip", "localhost", "The server ip")
	port     = flag.Int("port", 50050, "The server port")
	protocol = flag.String("protocol", "grpc", "The server protocol")
	weight   = flag.Int("weight", 1, "The server weight")
	status   = flag.String("status", "0", "The server status")
	num      = flag.Int("num", 3, "The server num")
)

// MiniGameRouterServer 实现了 MiniGameRouter gRPC 服务
type MiniGameRouterServer struct {
	pb.UnimplementedMiniGameRouterServer
	port  int
	times int
	mu    sync.Mutex
}

// SayHello 实现了 gRPC 的 SayHello 方法
func (s *MiniGameRouterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("Recv msg: %v times: %d\n", req.Msg, s.times)
	s.times += 1
	return &pb.HelloResponse{
		Msg: fmt.Sprintf("Hello, I am %s_%s:%d", *name, *ip, s.port),
	}, nil
}

// startGRPCServer 启动 gRPC 服务器
func startGRPCServer(port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	pb.RegisterMiniGameRouterServer(s, &MiniGameRouterServer{port: port, times: 1})
	log.Printf("Server is listening on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

var (
	epSlice     []*endPoint.EndPoint
	clientSlice []*pb.MiniGameRouterClient
	connSlice   []*grpc.ClientConn
	numServers  int
	wg          sync.WaitGroup
	mu          sync.Mutex
)

func initGRPCClients() {
	limiter := rate.NewLimiter(5, 5)
	for i := 0; i < numServers; i++ {
		limiter.Wait(context.Background())
		wg.Add(1)
		go func(i int, port int) {
			defer wg.Done()
			currPort := port + 2*i
			go startGRPCServer(currPort)
			currWeight := rand.Intn(10) + 1
			if *num == 1 {
				currWeight = *weight
			}
			endpoint := endPoint.EndPoint{
				Name:     name,
				Ip:       ip,
				Port:     &currPort,
				Protocol: protocol,
				Weight:   &currWeight,
				Status:   status,
			}
			conn, client := SDK.Init(&endpoint)
			//defer conn.Close()
			mu.Lock()
			epSlice = append(epSlice, &endpoint)
			clientSlice = append(clientSlice, &client)
			connSlice = append(connSlice, conn)
			mu.Unlock()
		}(i, *port)
	}
	wg.Wait()
}

func closeConnections() {
	for _, conn := range connSlice {
		conn.Close()
	}
}
func main() {
	flag.Parse()
	numServers = *num
	epSlice = make([]*endPoint.EndPoint, 0, numServers)
	clientSlice = make([]*pb.MiniGameRouterClient, 0, numServers)
	initGRPCClients()
	defer closeConnections()
	SDK.HandleUserInput(epSlice, clientSlice)
}
