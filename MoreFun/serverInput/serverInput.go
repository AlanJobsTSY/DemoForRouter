package main

import (
	"MoreFun/SDK"
	"MoreFun/endPoint"
	pb "MoreFun/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"net"
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
	//portNS   = flag.Int("portNS", 50050, "The server port")
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
	var lis net.Listener
	var err error
	for i := 0; i < 10; i++ {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return
	}

	s := grpc.NewServer()
	pb.RegisterMiniGameRouterServer(s, &MiniGameRouterServer{port: port, times: 1})
	log.Printf("Server is listening on port %d", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

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

func main() {
	//test
	// 创建 CPU 分析文件
	/*
		f, err := os.Create("server.pprof")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()

		// 开始 CPU 分析
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	*/

	flag.Parse()

	currPort := *port
	for {
		if isListen(currPort) {
			break
		}
		currPort = *port + 2*(rand.Intn(1000))
	}
	go startGRPCServer(currPort)
	endpoint := endPoint.EndPoint{
		Name:     name,
		Ip:       ip,
		Port:     &currPort,
		Protocol: protocol,
		Weight:   weight,
		Status:   status,
	}
	rRes, conn, _ := SDK.Init(&endpoint)
	defer conn.Close()
	if rRes == nil {
		return
	}
	//SDK.Input(&endpoint, client)
	log.Printf("server exit")
}
