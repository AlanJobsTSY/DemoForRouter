package SDK

import (
	"MoreFun/endPoint"
	pb "MoreFun/proto"
	"bufio"
	"fmt"
	"log"
	"os"
)

func Input(endPoint *endPoint.EndPoint, client pb.MiniGameRouterClient) {
	// 请求别人的服务
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter the name of the service you need and the 'ip:port' for the service if your know it (or press Ctrl+C to exit):")
	for scanner.Scan() {
		if scanner.Text() == "exit" {
			return
		}
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
