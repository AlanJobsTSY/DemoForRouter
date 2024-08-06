package routerStrategy

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

// GetAddr 根据服务配置选择路由策略，返回服务地址
func GetAddr(myServicesStorage *ServicesStorage, svrConfig *ServiceConfig) string {
	var addr string = ""
	switch svrConfig.Strategy {
	case ConsistentHash:
		addr = consistentHash(myServicesStorage, svrConfig.ServiceName)
	case Random:
		addr = random(myServicesStorage, svrConfig.ServiceName)
	case WeightedRoundRobin:
		addr = weightedRoundRobin(myServicesStorage, svrConfig.ServiceName)
	case FixedRoute:
		addr = fixedRoute(myServicesStorage, svrConfig.ServiceName, svrConfig.FixedRoute)
	}

	return addr
}

// consistentHash 一致性哈希路由策略（未实现）
func consistentHash(myServicesStorage *ServicesStorage, svrName string) string {
	var addr string = ""
	return addr
}

// random 随机路由策略
func random(myServicesStorage *ServicesStorage, svrName string) string {
	var addr string = ""
	myServicesStorage.RLock()
	defer myServicesStorage.RUnlock()
	if instances, ok := myServicesStorage.ServicesStorage[svrName]; ok && len(instances) > 0 {
		instanceKeys := make([]string, 0, len(instances))
		for _, instanceAddr := range instances {
			instanceKeys = append(instanceKeys, instanceAddr)
		}
		// 随机选择一个实例地址
		parts := strings.Split(instanceKeys[rand.Intn(len(instanceKeys))], ":")
		addr = fmt.Sprintf("%s:%s", parts[1], parts[2])
	}
	return addr
}

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

// fixedRoute 固定路由策略
func fixedRoute(myServicesStorage *ServicesStorage, svrName string, fixedRoute string) string {
	return fixedRoute
}
