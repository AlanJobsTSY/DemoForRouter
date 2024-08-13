package main

import (
	"MoreFun/SDK"
	"MoreFun/endPoint"
	pb "MoreFun/proto"
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
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
		Msg: fmt.Sprintf("Hello, I am %s_%s:%d", *name, *ip, *port),
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

var (
	epSlice     []*endPoint.EndPoint
	clientSlice []*pb.MiniGameRouterClient
	connSlice   []*grpc.ClientConn
	numServers  = 2
	wg          sync.WaitGroup
	mu          sync.Mutex
)

func initGRPCClients() {
	for i := 0; i < numServers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			currPort := *port + 2*i
			go startGRPCServer(currPort)
			endpoint := endPoint.EndPoint{
				Name:     name,
				Ip:       ip,
				Port:     &currPort,
				Protocol: protocol,
				Weight:   weight,
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

func handleUserInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter index: ")
		scanner.Scan()
		input := scanner.Text()
		index, err := strconv.Atoi(input)
		if err != nil || index < 0 || index >= len(epSlice) {
			fmt.Println("Invalid index")
			continue
		}

		endpoint := epSlice[index]
		client := *clientSlice[index]

		SDK.Input(endpoint, client)
	}
}
func closeConnections() {
	for _, conn := range connSlice {
		conn.Close()
	}
}
func main() {
	flag.Parse()
	epSlice = make([]*endPoint.EndPoint, 0, numServers)
	clientSlice = make([]*pb.MiniGameRouterClient, 0, numServers)
	initGRPCClients()
	defer closeConnections()
	handleUserInput()
}
