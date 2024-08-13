package routerStrategy

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
	"sync"
	"time"
)

func fastestResponse(myServicesStorage *ServicesStorage, svrName string) string {
	var addr string
	var minn = time.Duration(1<<63 - 1)
	var mu sync.Mutex
	var wg sync.WaitGroup

	myServicesStorage.RLock()
	defer myServicesStorage.RUnlock()

	if instances, ok := myServicesStorage.ServicesStorage[svrName]; ok && len(instances) > 0 {
		for _, instanceAddr := range instances {
			parts := strings.Split(instanceAddr, ":")
			addrInstance := fmt.Sprintf("%s:%s", parts[1], parts[2])

			wg.Add(1)
			// 用协程可以大大缩短大批量的连接时间计算
			go func(addrInstance string) {
				defer wg.Done()
				//记录开始时间
				startTime := time.Now()
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()

				conn, err := grpc.DialContext(ctx, addrInstance, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
				if err != nil {
					return
				}
				defer conn.Close()
				//记录结束时间
				endTime := time.Now()
				connectionTime := endTime.Sub(startTime)
				log.Printf("%s %v", addrInstance, connectionTime)
				// 保护minn和addr这两个共享变量
				mu.Lock()
				defer mu.Unlock()
				if connectionTime < minn {
					minn = connectionTime
					addr = addrInstance
				}
			}(addrInstance)
		}
		wg.Wait()
	}

	return addr
}
