package main

import (
	"MoreFun/discovery"
	pb "MoreFun/discovery/proto"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

var (
	port = flag.Int("port", 50051, "")
)

var times int = 0

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	times++
	fmt.Printf("Recv Client msg:%v %d\n", in.Msg, times)
	return &pb.HelloReply{
		Msg: "hello client",
	}, nil
}

func ServiceRegister(s grpc.ServiceRegistrar, srv pb.GreeterServer) {
	pb.RegisterGreeterServer(s, srv)
	s1 := &discovery.Service{
		Name:     "hello.Greeter",
		Port:     strconv.Itoa(*port),
		IP:       "localhost",
		Protocol: "grpc",
	}
	go discovery.ServiceRegister(s1)
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}
	s := grpc.NewServer()
	ServiceRegister(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}
