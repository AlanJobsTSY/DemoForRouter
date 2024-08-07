package routerStrategy

import (
	"fmt"
	"math/rand"
	"strings"
)

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
