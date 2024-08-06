package main

import (
	"MoreFun/etcd"
	pb "MoreFun/proto"
	"bufio"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 定义命令行参数
var (
	name     = flag.String("name", "A", "The server name")
	ip       = flag.String("ip", "localhost", "The server ip")
	port     = flag.Int("port", 50050, "The server port")
	protocol = flag.String("protocol", "grpc", "The server protocol")
	weight   = flag.String("weight", "1", "The server weight")
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
		Msg: fmt.Sprintf("Hello, I am %s:%d", *name, *port),
	}, nil
}

// DiscoverService 发现服务
func DiscoverService(client pb.MiniGameRouterClient, serviceName string, fixedRouter string) (*pb.DiscoverServiceResponse, error) {
	req := &pb.DiscoverServiceRequest{
		FromMsg:         fmt.Sprintf("%s:%d", *name, *port),
		ToMsg:           serviceName,
		FixedRouterAddr: fixedRouter,
	}
	return client.DiscoverService(context.Background(), req)
}

// RegisterService 注册服务
func RegisterService(client pb.MiniGameRouterClient) (*pb.RegisterServiceResponse, error) {
	req := &pb.RegisterServiceRequest{
		Service: &pb.Service{
			Name:     *name,
			Ip:       *ip,
			Port:     strconv.Itoa(*port + 1),
			Protocol: *protocol,
			Weight:   *weight,
		},
	}
	return client.RegisterService(context.Background(), req)
}

// startSidecar 启动 sidecar
func startSidecar(port int) error {
	sidecarPort := port + 1
	cmd := exec.Command("go", "run", "./sidecar/sidecar.go", "--port", strconv.Itoa(sidecarPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
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

// connectToSidecar 连接到 sidecar
func connectToSidecar(port int) (*grpc.ClientConn, pb.MiniGameRouterClient, error) {
	addr := fmt.Sprintf("localhost:%d", port+1)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to sidecar: %v", err)
	}
	client := pb.NewMiniGameRouterClient(conn)
	return conn, client, nil
}

func main() {
	flag.Parse()

	// 启动 sidecar
	if err := startSidecar(*port); err != nil {
		log.Fatalf("Failed to start sidecar: %v", err)
	}

	// 等待 sidecar 启动完成
	time.Sleep(etcd.DialTimeout)

	// 启动 gRPC 服务器
	go startGRPCServer(*port)

	// 连接 sidecar
	conn, client, err := connectToSidecar(*port)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer conn.Close()

	// 注册自己的服务
	rRes, err := RegisterService(client)
	if err != nil {
		log.Fatalf("Could not register service: %v", err)
	}
	log.Printf("Response: %s", rRes.Msg)

	// 请求别人的服务
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter the name of the service you need and the 'ip:port' for the service if your know it (or press Ctrl+C to exit):")
	for scanner.Scan() {
		var serviceName string
		var fixedRouterAddr string
		if scanner.Text() == "" {
			continue
		}
		parts := strings.Fields(scanner.Text())
		if len(parts) == 1 {
			serviceName = parts[0]
			fixedRouterAddr = ""
		} else {
			serviceName = parts[0]
			partsFixedRouterAddr := strings.Split(parts[1], ":")
			if len(partsFixedRouterAddr) != 2 {
				log.Printf("Wrong fomat for 'ip:port'")
				continue
			}
			portInt, err := strconv.Atoi(partsFixedRouterAddr[1])
			if err != nil {
				log.Printf("Wrong port")
				continue
			}
			fixedRouterAddr = fmt.Sprintf("%s:%d", partsFixedRouterAddr[0], portInt+1)
		}
		helloRes, err := DiscoverService(client, serviceName, fixedRouterAddr)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		log.Printf("Recv msg: %s", helloRes.Msg)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
}
