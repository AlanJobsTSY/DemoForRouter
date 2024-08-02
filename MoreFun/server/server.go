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
	"time"
)

var (
	name     = flag.String("name", "A", "")
	ip       = flag.String("ip", "localhost", "")
	port     = flag.Int("port", 50050, "The server port")
	protocol = flag.String("protocol", "grpc", "The server protocol")
	weight   = flag.String("weight", "1", "")
)

type MiniGameRouterServer struct {
	pb.UnimplementedMiniGameRouterServer
}

func (s *MiniGameRouterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("Recv msg:%v \n", req.Msg)
	return &pb.HelloResponse{
		Msg: fmt.Sprintf("Hello, I am %s:%d", *name, *port),
	}, nil
}

func DiscoverService(client pb.MiniGameRouterClient, serviceName string) (*pb.DiscoverServiceResponse, error) {
	req := &pb.DiscoverServiceRequest{
		FromMsg: fmt.Sprintf("%s:%d", *name, *port),
		ToMsg:   serviceName,
	}
	return client.DiscoverService(context.Background(), req)
}

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

func startSidecar(port int) error {
	sidecarPort := port + 1
	cmd := exec.Command("go", "run", "./sidecar/sidecar.go", "--port", strconv.Itoa(sidecarPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

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
	fmt.Println("Enter the name of the service you need (or press Ctrl+C to exit):")
	for scanner.Scan() {
		serviceName := scanner.Text()
		if serviceName == "" {
			continue
		}
		helloRes, err := DiscoverService(client, serviceName)
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
