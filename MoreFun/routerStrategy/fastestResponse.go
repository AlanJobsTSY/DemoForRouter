package routerStrategy

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"strings"
	"time"
)

func fastestResponse(myServicesStorage *ServicesStorage, svrName string) string {
	var addr string
	var minn = time.Duration(1<<63 - 1)
	myServicesStorage.RLock()
	defer myServicesStorage.RUnlock()
	if instances, ok := myServicesStorage.ServicesStorage[svrName]; ok && len(instances) > 0 {
		for _, instanceAddr := range instances {
			parts := strings.Split(instanceAddr, ":")
			// 获取实例的权重
			addrInstance := fmt.Sprintf("%s:%s", parts[1], parts[2])
			startTime := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			// 创建 gRPC 连接，使用 grpc.WithBlock() 确保连接完全建立
			conn, err := grpc.DialContext(ctx, addrInstance, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
			if err != nil {
				continue
			}
			conn.Close()
			// 记录结束时间
			endTime := time.Now()
			connectionTime := endTime.Sub(startTime)
			if minn > connectionTime {
				minn = connectionTime
				addr = addrInstance
			}
		}
	}
	return addr
}
