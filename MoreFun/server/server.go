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
	"os/exec"
	"sync"
	"time"
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

func killProcessOnPort(port int) error {
	cmd := exec.Command("fuser", "-k", fmt.Sprintf("%d/tcp", port))
	return cmd.Run()
}

// startGRPCServer 启动 gRPC 服务器
func startGRPCServer(port int) {
	var lis net.Listener
	var err error
	for i := 0; i < 10; i++ {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			break
		}
		if killErr := killProcessOnPort(port); killErr != nil {
			log.Printf("Failed to free port %d: %v", port, killErr)
		} else {
			log.Printf("Successfully freed port %d", port)
		}
		time.Sleep(2 * time.Second)
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

func isListen(port int) (flag bool) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Server failed to listen on port %d: %v\n", port, err)
		return false
	}
	defer ln.Close() // 确保在函数返回前关闭监听器

	sidecarLn, err := net.Listen("tcp", fmt.Sprintf(":%d", port+1))
	if err != nil {
		log.Printf("Sidecar failed to listen on port %d: %v\n", port+1, err)
		return false
	}
	defer sidecarLn.Close() // 确保在函数返回前关闭监听器

	return true
}
func initGRPCClients() {
	limiter := rate.NewLimiter(20, 20)
	for i := 0; i < numServers; i++ {
		limiter.Wait(context.Background())
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			currPort := *port + 2*i
			if ok := isListen(currPort); ok == false {
				return
			}
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

		}(i)
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
