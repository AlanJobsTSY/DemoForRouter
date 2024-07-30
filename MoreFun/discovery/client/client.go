package main

import (
	"MoreFun/discovery"
	pb "MoreFun/discovery/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func getServerAddress(svcName string) string {
	s := discovery.ServiceDiscovery(svcName)
	if s == nil {
		return ""
	}
	if s.IP == "" && s.Port == "" {
		return ""
	}
	return s.IP + ":" + s.Port
}
func sayHello() {
	addr := getServerAddress("hello.Greeter")
	if addr == "" {
		log.Println("未发现可用服务")
		return
	}
	fmt.Println(addr)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	in := &pb.HelloRequest{
		Msg: "hello server",
	}
	r, err := c.SayHello(context.Background(), in)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Recv server msg:", r.Msg)

}

func main() {
	//log.SetFlags(log.Llongfile)
	go discovery.WatchServiceName("hello.Greeter")
	for {
		sayHello()
		time.Sleep(time.Second * 2)
	}

}
