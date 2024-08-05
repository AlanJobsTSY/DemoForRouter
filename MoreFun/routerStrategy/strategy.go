package routerStrategy

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

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

func consistentHash(myServicesStorage *ServicesStorage, svrName string) string {
	var addr string = ""
	return addr
}

func random(myServicesStorage *ServicesStorage, svrName string) string {
	var addr string = ""
	myServicesStorage.RLock()
	defer myServicesStorage.RUnlock()
	if instances, ok := myServicesStorage.ServicesStorage[svrName]; ok && len(instances) > 0 {
		instanceKeys := make([]string, 0, len(instances))
		for _, instanceAddr := range instances {
			instanceKeys = append(instanceKeys, instanceAddr)
		}
		parts := strings.Split(instanceKeys[rand.Intn(len(instanceKeys))], ":")
		addr = fmt.Sprintf("%s:%s", parts[1], parts[2])
	}
	return addr
}

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
			//log.Printf("tsytsy%s", parts[4])
			part, _ := strconv.Atoi(parts[4])
			totalWeight += part
			myServicesStorage.CurrentWeight[svrName][key] += part
		}
		minn := math.MinInt64
		for key, _ := range instances {
			if myServicesStorage.CurrentWeight[svrName][key] >= minn {
				minn = myServicesStorage.CurrentWeight[svrName][key]
				addr = key
			}

		}
		myServicesStorage.CurrentWeight[svrName][addr] -= totalWeight
		addr = myServicesStorage.ServicesStorage[svrName][addr]
		parts := strings.Split(addr, ":")
		addr = fmt.Sprintf("%s:%s", parts[1], parts[2])
	}
	return addr
}

func fixedRoute(myServicesStorage *ServicesStorage, svrName string, fixedRoute string) string {
	return fixedRoute
}
