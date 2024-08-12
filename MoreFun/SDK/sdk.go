package SDK

import (
	"MoreFun/endPoint"
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

// DiscoverService 发现服务
func DiscoverService(endPoint *endPoint.EndPoint, client pb.MiniGameRouterClient, serviceName string, fixedRouter string) (*pb.DiscoverServiceResponse, error) {
	req := &pb.DiscoverServiceRequest{
		FromMsg:         fmt.Sprintf("%s_%s:%d", *endPoint.Name, *endPoint.Ip, *endPoint.Port),
		ToMsg:           serviceName,
		FixedRouterAddr: fixedRouter,
		Status:          *endPoint.Status,
	}
	return client.DiscoverService(context.Background(), req)
}

// startSidecar 启动 sidecar
func startSidecar(endPoint *endPoint.EndPoint) error {
	sidecarPort := *endPoint.Port + 1
	sidecarIP := *endPoint.Ip
	cmd := exec.Command("go", "run", "./sidecar/sidecar.go", "--ip", sidecarIP, "--port", strconv.Itoa(sidecarPort))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func Init(endPoint *endPoint.EndPoint) (*grpc.ClientConn, pb.MiniGameRouterClient) {
	// 启动 sidecar
	if err := startSidecar(endPoint); err != nil {
		log.Fatalf("Failed to start sidecar: %v", err)
	}

	// 连接 sidecar
	conn, client, err := connectToSidecar(endPoint)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// 注册自己的服务
	rRes, err := RegisterService(endPoint, client)
	if err != nil {
		log.Fatalf("Could not register service: %v", err)
	}
	log.Printf("Response: %s", rRes.Msg)

	return conn, client
}

// connectToSidecar 连接到 sidecar
func connectToSidecar(endPoint *endPoint.EndPoint) (*grpc.ClientConn, pb.MiniGameRouterClient, error) {
	addr := fmt.Sprintf("%s:%d", *endPoint.Ip, *endPoint.Port+1)
	for i := 0; i < 200; i++ {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			client := pb.NewMiniGameRouterClient(conn)
			return conn, client, nil
		}
		// 连接失败，等待 0.1 秒后重试
		time.Sleep(100 * time.Millisecond)
	}
	return nil, nil, fmt.Errorf("failed to connect to sidecar")

}

// RegisterService 注册服务
func RegisterService(endPoint *endPoint.EndPoint, client pb.MiniGameRouterClient) (*pb.RegisterServiceResponse, error) {
	req := &pb.RegisterServiceRequest{
		Service: &pb.Service{
			Name:     *endPoint.Name,
			Ip:       *endPoint.Ip,
			Port:     strconv.Itoa(*endPoint.Port + 1),
			Protocol: *endPoint.Protocol,
			Weight:   *endPoint.Weight,
			Status:   *endPoint.Status,
			ConnNum:  "0",
		},
	}
	return client.RegisterService(context.Background(), req)
}

func grpcDiscover(endPoint *endPoint.EndPoint, client pb.MiniGameRouterClient, scan string) {
	var serviceName string
	var fixedRouterAddr string
	parts := strings.Fields(scan)
	if len(parts) == 1 { // 普通路由
		serviceName = parts[0]
		fixedRouterAddr = ""
	} else { // 指定目标路由
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
	helloRes, err := DiscoverService(endPoint, client, serviceName, fixedRouterAddr)
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
func Input(endPoint *endPoint.EndPoint, client pb.MiniGameRouterClient) {
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
		grpcDiscover(endPoint, client, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
}
