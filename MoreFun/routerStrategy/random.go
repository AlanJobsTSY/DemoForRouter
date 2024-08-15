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
	// VERSION3
	if len(myServicesStorage.RandomStorage.RandomList) > 0 {
		//log.Printf("len   %d", len(myServicesStorage.RandomStorage.RandomList))
		indexRandom := rand.Intn(len(myServicesStorage.RandomStorage.RandomList))

		value := myServicesStorage.RandomStorage.RandomList[indexRandom]
		parts := strings.Split(value, ":")
		addr = fmt.Sprintf("%s:%s", parts[1], parts[2])
		//log.Printf("ADDR %s", addr)
	}
	// VERSION1
	// 随机选择一个实例地址
	/*
		if instances, ok := myServicesStorage.ServicesStorage[svrName]; ok && len(instances) > 0 {
			instanceKeys := make([]string, 0, len(instances))
			for _, instanceAddr := range instances {
				instanceKeys = append(instanceKeys, instanceAddr)
			}
			// 随机选择一个实例地址
			parts := strings.Split(instanceKeys[rand.Intn(len(instanceKeys))], ":")
			addr = fmt.Sprintf("%s:%s", parts[1], parts[2])
		}
	*/
	// VERSION2
	/*
		//由于go的map本身是无序的，本意是只选map的第一个达成随机的效果，但是测试下来分布却不太随机
		if instances, ok := myServicesStorage.ServicesStorage[svrName]; ok && len(instances) > 0 {
			for _, instanceAddr := range instances {
				parts := strings.Split(instanceAddr, ":")
				addr = fmt.Sprintf("%s:%s", parts[1], parts[2])
				break
			}
		}*/
	return addr
}

func (ss *ServicesStorage) DeleteRandom(svrName string) {
	if ss.RandomStorage.RandomMap[svrName] < len(ss.RandomStorage.RandomList) {
		ss.RandomStorage.RandomList[ss.RandomStorage.RandomMap[svrName]] = ss.RandomStorage.RandomList[len(ss.RandomStorage.RandomList)-1]
	}
	ss.RandomStorage.RandomList = ss.RandomStorage.RandomList[:len(ss.RandomStorage.RandomList)-1]
	delete(ss.RandomStorage.RandomMap, svrName)

}

func (ss *ServicesStorage) AddRandom(svrName string, ip string) {
	ss.RandomStorage.RandomList = append(ss.RandomStorage.RandomList, ip)
	ss.RandomStorage.RandomMap[svrName] = len(ss.RandomStorage.RandomList) - 1
}
