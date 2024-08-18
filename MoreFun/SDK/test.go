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
	"strings"
	"sync"
	"time"
)

func HandleUserInput(endPoint *endPoint.EndPoint, client *pb.MiniGameRouterClient) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Choose a test scenario:")
		fmt.Println("1. Fixed target routing test")
		fmt.Println("2. Other types of routing tests")
		fmt.Println("3. Batch set dynamic key-value routing test")
		fmt.Println("4. Batch access dynamic key-value routing test")
		fmt.Println("5. Flexible experience")
		fmt.Print("Enter your choice: ")
		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			testFixedTypeRouting(endPoint, client)
		case "2":
			testOtherTypeRouting(endPoint, client)
		case "3":
			testRegisterDynamicKeyValueRouting(endPoint, client)
		case "4":
			testOtherTypeRouting(endPoint, client)
		case "5":
			Input(endPoint, *client)
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func testFixedTypeRouting(endPoint *endPoint.EndPoint, client *pb.MiniGameRouterClient) {
	var svrDiscoverTimeTotal int64
	var returnTimeTotal int64

	fmt.Print("Enter the type of service: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	serviceType := scanner.Text()
	fmt.Print("Enter the target route: ")
	scanner.Scan()
	targetRoute := scanner.Text()

	partsRouterAddr := strings.Split(targetRoute, ":")
	if len(partsRouterAddr) != 2 {
		log.Printf("Wrong fomat for 'ip:port'")
		return
	}
	portInt, err := strconv.Atoi(partsRouterAddr[1])
	if err != nil {
		log.Printf("Wrong port")
		return
	}
	targetRoute = fmt.Sprintf("%s:%d", partsRouterAddr[0], portInt+1)

	fmt.Print("Enter the number of times to send: ")
	scanner.Scan()
	times, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid number")
		return
	}
	limiter := rate.NewLimiter(1000, 1000)
	var wg sync.WaitGroup
	var mu sync.Mutex
	// 记录开始时间
	startTime := time.Now()
	for i := 0; i < times; i++ {
		limiter.Wait(context.Background())
		wg.Add(1)
		go func() {
			defer wg.Done()
			helloRes, err := discoverService(endPoint, *client, serviceType, targetRoute)
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Printf("Recv msg: %s", helloRes.Msg)
			mu.Lock()
			svrDiscoverTimeTotal += helloRes.SvrDiscoverTime
			returnTimeTotal += helloRes.ReturnTime
			mu.Unlock()
		}()
	}
	wg.Wait()
	endTime := time.Now()
	fmt.Printf("svrDiscoverTimeTotal taken: %v μs\n", svrDiscoverTimeTotal)
	fmt.Printf("returnTimeTotal taken: %v ms\n", returnTimeTotal)
	fmt.Printf("Total time taken: %v ms\n", endTime.Sub(startTime).Milliseconds())
}

func testOtherTypeRouting(endPoint *endPoint.EndPoint, client *pb.MiniGameRouterClient) {
	var svrDiscoverTimeTotal int64
	var returnTimeTotal int64

	fmt.Print("Enter the type of service or dynamicKey: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	serviceType := scanner.Text()

	fmt.Print("Enter the number of times to send: ")
	scanner.Scan()
	times, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid number")
		return
	}
	limiter := rate.NewLimiter(1000, 1000)
	var wg sync.WaitGroup
	var mu sync.Mutex
	// 记录开始时间
	startTime := time.Now()
	for i := 0; i < times; i++ {
		limiter.Wait(context.Background())
		wg.Add(1)
		go func() {
			defer wg.Done()
			helloRes, err := discoverService(endPoint, *client, serviceType, "")
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Printf("Recv msg: %s", helloRes.Msg)
			mu.Lock()
			svrDiscoverTimeTotal += helloRes.SvrDiscoverTime
			returnTimeTotal += helloRes.ReturnTime
			mu.Unlock()
		}()
	}
	wg.Wait()
	endTime := time.Now()
	fmt.Printf("svrDiscoverTimeTotal taken: %v μs\n", svrDiscoverTimeTotal)
	fmt.Printf("returnTimeTotal taken: %v ms\n", returnTimeTotal)
	fmt.Printf("Total time taken: %v ms\n", endTime.Sub(startTime).Milliseconds())
}

func testRegisterDynamicKeyValueRouting(endPoint *endPoint.EndPoint, client *pb.MiniGameRouterClient) {
	fmt.Print("Enter the dynamicKey: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	dynamicKey := scanner.Text()
	fmt.Print("Enter the dynamicValue: ")
	scanner.Scan()
	dynamicValue := scanner.Text()

	partsRouterAddr := strings.Split(dynamicValue, ":")
	if len(partsRouterAddr) != 2 {
		log.Printf("Wrong fomat for 'ip:port'")
		return
	}
	portInt, err := strconv.Atoi(partsRouterAddr[1])
	if err != nil {
		log.Printf("Wrong port")
		return
	}
	dynamicValue = fmt.Sprintf("%s:%d", partsRouterAddr[0], portInt+1)

	fmt.Print("Enter the timeout: ")
	scanner.Scan()
	timeout, err := strconv.Atoi(scanner.Text())
	fmt.Print("Enter the number of times to set: ")
	scanner.Scan()
	times, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid number")
		return
	}
	limiter := rate.NewLimiter(150, 150)
	var wg sync.WaitGroup
	// 记录开始时间
	startTime := time.Now()
	for i := 0; i < times; i++ {
		limiter.Wait(context.Background())
		wg.Add(1)
		go func(i string) {
			defer wg.Done()
			rsp, err := setCustomRoute(*client, dynamicKey+"_"+i, dynamicValue, strconv.Itoa(timeout))
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Printf("Recv msg: %s", rsp.Msg)
		}(strconv.Itoa(i))
	}
	wg.Wait()
	elapsedTime := time.Since(startTime)
	fmt.Printf("Total time taken: %v ms\n", elapsedTime.Milliseconds())
}

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
			setCustomRouteInput(client)
			continue
		}
		grpcDiscover(endPoint, client, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
}
