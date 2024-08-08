package routerStrategy

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// weightedRoundRobin 平滑加权轮询路由策略

func weightedRoundRobin(myServicesStorage *ServicesStorage, svrName string) string {
	var addr string = ""
	myServicesStorage.RLock()
	defer myServicesStorage.RUnlock()
	if instances, ok := myServicesStorage.ServicesStorage[svrName]; ok && len(instances) > 0 {
		totalWeight := 0
		if myServicesStorage.CurrentWeight[svrName] == nil {
			myServicesStorage.CurrentWeight[svrName] = make(map[string]int)
		}
		for key, instanceAddr := range instances {
			parts := strings.Split(instanceAddr, ":")
			// 获取实例的权重
			part, _ := strconv.Atoi(parts[4])
			totalWeight += part
			myServicesStorage.CurrentWeight[svrName][key] += part
		}
		minn := math.MinInt64
		// 选择当前权重最大的实例
		for key := range instances {
			if myServicesStorage.CurrentWeight[svrName][key] >= minn {
				minn = myServicesStorage.CurrentWeight[svrName][key]
				addr = key
			}
		}
		// 减去总权重
		myServicesStorage.CurrentWeight[svrName][addr] -= totalWeight
		addr = myServicesStorage.ServicesStorage[svrName][addr]
		parts := strings.Split(addr, ":")
		addr = fmt.Sprintf("%s:%s", parts[1], parts[2])
	}
	return addr
}