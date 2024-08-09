package SDK

import (
	"MoreFun/etcd"
	pb "MoreFun/proto"
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// MiniGameRouterServer 实现了 MiniGameRouter gRPC 服务
type MiniGameRouterServer struct {
	pb.UnimplementedMiniGameRouterServer
}

var (
	name     *string
	ip       *string
	port     *int
	protocol *string
	weight   *string
	status   *string
)

// DiscoverService 发现服务
func DiscoverService(client pb.MiniGameRouterClient, serviceName string, fixedRouter string) (*pb.DiscoverServiceResponse, error) {
	req := &pb.DiscoverServiceRequest{
		FromMsg:         fmt.Sprintf("%s_%d", *name, *port),
		ToMsg:           serviceName,
		FixedRouterAddr: fixedRouter,
		Status:          *status,
	}
	return client.DiscoverService(context.Background(), req)
}

// startSidecar 启动 sidecar
func startSidecar(ip string, port int) error {
	sidecarPort := port + 1
	sidecarIP := ip
	cmd := exec.Command("go", "run", "./sidecar/sidecar.go", "--ip", sidecarIP, "--port", strconv.Itoa(sidecarPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func InitOne(svrName *string, svrIP *string, svrPort *int, svrProtocol *string, svrWeight *string, svrStatus *string) {
	name = svrName
	ip = svrIP
	port = svrPort
	protocol = svrProtocol
	weight = svrWeight
	status = svrStatus
	// 启动 sidecar
	if err := startSidecar(*ip, *port); err != nil {
		log.Fatalf("Failed to start sidecar: %v", err)
	}
	// 等待 sidecar 启动完成
	time.Sleep(etcd.DialTimeout)

}

// connectToSidecar 连接到 sidecar
func connectToSidecar(port int) (*grpc.ClientConn, pb.MiniGameRouterClient, error) {
	addr := fmt.Sprintf("%s:%d", *ip, port+1)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to sidecar: %v", err)
	}
	client := pb.NewMiniGameRouterClient(conn)
	return conn, client, nil
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
			Status:   *status,
			ConnNum:  "0",
		},
	}
	return client.RegisterService(context.Background(), req)
}

func InitTwo() (*grpc.ClientConn, pb.MiniGameRouterClient) {
	// 连接 sidecar
	conn, client, err := connectToSidecar(*port)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// 注册自己的服务
	rRes, err := RegisterService(client)
	if err != nil {
		log.Fatalf("Could not register service: %v", err)
	}
	log.Printf("Response: %s", rRes.Msg)

	return conn, client
}
func grpcDiscover(client pb.MiniGameRouterClient, scan string) {
	var serviceName string
	var fixedRouterAddr string
	parts := strings.Fields(scan)
	if len(parts) == 1 {
		serviceName = parts[0]
		fixedRouterAddr = ""
	} else {
		serviceName = parts[0]
		partsFixedRouterAddr := strings.Split(parts[1], ":")
		if len(partsFixedRouterAddr) != 2 {
			log.Printf("Wrong fomat for 'ip:port'")
			return
		}
		portInt, err := strconv.Atoi(partsFixedRouterAddr[1])
		if err != nil {
			log.Printf("Wrong port")
			return
		}
		fixedRouterAddr = fmt.Sprintf("%s:%d", partsFixedRouterAddr[0], portInt+1)
	}
	helloRes, err := DiscoverService(client, serviceName, fixedRouterAddr)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	log.Printf("Recv msg: %s", helloRes.Msg)
}
func setCustomRoute(client pb.MiniGameRouterClient) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	parts := strings.Fields(scanner.Text())
	partsRouterAddr := strings.Split(parts[1], ":")
	if len(partsRouterAddr) != 2 {
		log.Printf("Wrong fomat for 'ip:port'")
		return
	}
	portInt, err := strconv.Atoi(partsRouterAddr[1])
	if err != nil {
		log.Printf("Wrong port")
		return
	}
	routerAddr := fmt.Sprintf("%s:%d", partsRouterAddr[0], portInt+1)
	req := &pb.SetRequest{
		Key:     parts[0],
		Value:   routerAddr,
		Timeout: parts[2],
	}
	rsp, err := client.SetCustomRoute(context.Background(), req)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	log.Printf(rsp.Msg)
}
func Input(client pb.MiniGameRouterClient) {
	// 请求别人的服务
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter the name of the service you need and the 'ip:port' for the service if your know it (or press Ctrl+C to exit):")
	for scanner.Scan() {

		if scanner.Text() == "" {
			continue
		}
		if scanner.Text() == "set" {
			setCustomRoute(client)
			continue
		}
		grpcDiscover(client, scanner.Text())

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
}
