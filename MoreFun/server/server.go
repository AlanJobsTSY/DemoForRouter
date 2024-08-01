package main

import (
	pb "MoreFun/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"strconv"
)

var (
	name       = flag.String("name", "A", "")
	ip         = flag.String("ip", "localhost", "")
	port       = flag.Int("port", 50050, "The server port")
	protocol   = flag.String("protocol", "grpc", "The server protocol")
	weight     = flag.String("weight", "1", "")
	needServer = flag.String("needServer", "", "")
)

type MiniGameRouterServer struct {
	pb.UnimplementedMiniGameRouterServer
}

func (s *MiniGameRouterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("Recv Client msg:%v \n", req.Msg)
	return &pb.HelloResponse{
		Msg: fmt.Sprintf("Hello, I am %s:%s", *name, strconv.Itoa(*port)),
	}, nil
}
func DiscoverService(client pb.MiniGameRouterClient, serviceName string) (*pb.DiscoverServiceResponse, error) {
	needServerReq := &pb.DiscoverServiceRequest{
		Msg: serviceName,
	}
	helloRes, err := client.DiscoverService(context.Background(), needServerReq)
	if err != nil {
		return nil, fmt.Errorf("could not discover service: %v", err)
	}
	return helloRes, nil
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
	rRes, err := client.RegisterService(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("could not register service: %v", err)
	}
	return rRes, nil
}
func main() {
	flag.Parse()

	//连接sidecar
	addr := fmt.Sprintf("localhost:%s", strconv.Itoa(*port+1))
	fmt.Println(addr)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	c := pb.NewMiniGameRouterClient(conn)

	//	注册自己的服务
	rRes, err := RegisterService(c)
	if err != nil {
		log.Fatalf("could not register service: %v", err)
	}

	log.Printf("Response: %s", rRes.Msg)

	// 请求别人的服务
	if *needServer != "" {
		helloRes, err := DiscoverService(c, *needServer)
		if err != nil {
			log.Fatalf("could not register service: %v", err)
		}
		log.Printf("Response: %s", helloRes.Msg)
	}

	//监听sidecar的请求
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
	log.Printf("tsytsytsy")
}
