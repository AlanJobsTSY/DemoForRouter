package SDK

import (
	"MoreFun/endPoint"
	pb "MoreFun/proto"
	"bufio"
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
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

// HandleUserInput 处理用户输入并选择测试场景
func HandleUserInput(epSlice []*endPoint.EndPoint, clientSlice []*pb.MiniGameRouterClient) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Choose a test scenario:")
		fmt.Println("1. All services send to a specific type of service, each sends multiple times")
		fmt.Println("2. A specific service sends to a specific type of service, sends multiple times")
		fmt.Println("3. All services send to a specific target route, each sends multiple times")
		fmt.Println("4. A specific service sends to a specific target route, sends multiple times")
		fmt.Println("5. A specific service registers multiple dynamic values")
		fmt.Print("Enter your choice: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "0":
			oneByOne(epSlice, clientSlice)
		case "1":
			testAllServicesToType(epSlice, clientSlice)
		case "2":
			testSpecificServiceToType(epSlice, clientSlice)
		case "3":
			testAllServicesToRoute(epSlice, clientSlice)
		case "4":
			testSpecificServiceToRoute(epSlice, clientSlice)
		case "5":
			//testRegisterMultipleValues(epSlice, clientSlice)
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func oneByOne(epSlice []*endPoint.EndPoint, clientSlice []*pb.MiniGameRouterClient) {
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
		Input(endpoint, client)
	}
}

// 测试场景 1: 所有服务发给某类服务，各发几次
func testAllServicesToType(epSlice []*endPoint.EndPoint, clientSlice []*pb.MiniGameRouterClient) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the type of service: ")
	scanner.Scan()
	serviceType := scanner.Text()
	fmt.Print("Enter the number of times to send: ")
	scanner.Scan()
	times, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid number")
		return
	}
	for index, _ := range epSlice {
		for i := 0; i < times; i++ {
			helloRes, err := discoverService(epSlice[index], *clientSlice[index], serviceType, "")
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Printf("Recv msg: %s", helloRes.Msg)
		}
	}
}

// 测试场景 2: 指定一个服务发给某类服务，发几次
// testSpecificServiceToType 并发测试指定服务发给某类服务
func testSpecificServiceToType(epSlice []*endPoint.EndPoint, clientSlice []*pb.MiniGameRouterClient) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the index of the service: ")
	scanner.Scan()
	index, err := strconv.Atoi(scanner.Text())
	if err != nil || index < 0 || index >= len(clientSlice) {
		fmt.Println("Invalid index")
		return
	}
	fmt.Print("Enter the type of service: ")
	scanner.Scan()
	serviceType := scanner.Text()
	fmt.Print("Enter the number of times to send: ")
	scanner.Scan()
	times, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid number")
		return
	}

	var wg sync.WaitGroup
	for i := 0; i < times; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			helloRes, err := discoverService(epSlice[index], *clientSlice[index], serviceType, "")
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Printf("Recv msg: %s", helloRes.Msg)
		}()
	}
	wg.Wait()
}

// 测试场景 3: 所有服务发给指定目标路由，各发几次
func testAllServicesToRoute(epSlice []*endPoint.EndPoint, clientSlice []*pb.MiniGameRouterClient) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the type of service: ")
	scanner.Scan()
	serviceType := scanner.Text()
	fmt.Print("Enter the target route: ")
	scanner.Scan()
	targetRoute := scanner.Text()
	fmt.Print("Enter the number of times to send: ")
	scanner.Scan()
	times, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid number")
		return
	}

	for index, _ := range epSlice {
		for i := 0; i < times; i++ {
			helloRes, err := discoverService(epSlice[index], *clientSlice[index], serviceType, targetRoute)
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Printf("Recv msg: %s", helloRes.Msg)
		}
	}
}

// 测试场景 4: 指定一个服务发给指定目标路由，发几次
func testSpecificServiceToRoute(epSlice []*endPoint.EndPoint, clientSlice []*pb.MiniGameRouterClient) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter the index of the service: ")
	scanner.Scan()
	index, err := strconv.Atoi(scanner.Text())
	if err != nil || index < 0 || index >= len(clientSlice) {
		fmt.Println("Invalid index")
		return
	}
	fmt.Print("Enter the type of service: ")
	scanner.Scan()
	serviceType := scanner.Text()
	fmt.Print("Enter the target route: ")
	scanner.Scan()
	targetRoute := scanner.Text()
	fmt.Print("Enter the number of times to send: ")
	scanner.Scan()
	times, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid number")
		return
	}
	limiter := rate.NewLimiter(5000, 5000)
	var wg sync.WaitGroup
	// 记录开始时间
	startTime := time.Now()
	for i := 0; i < times; i++ {
		limiter.Wait(context.Background())
		wg.Add(1)
		go func() {
			defer wg.Done()

			helloRes, err := discoverService(epSlice[index], *clientSlice[index], serviceType, targetRoute)
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Printf("Recv msg: %s", helloRes.Msg)
		}()
	}
	wg.Wait()
	elapsedTime := time.Since(startTime)
	fmt.Printf("Total time taken: %v ms\n", elapsedTime.Milliseconds())
}
