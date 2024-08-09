package routerStrategy

// GetAddr 根据服务配置选择路由策略，返回服务地址
func GetAddr(myServicesStorage *ServicesStorage, svrConfig *ServiceConfig, fromServer string, status string) string {
	var addr string = ""
	switch svrConfig.Strategy {
	case ConsistentHash:
		addr = consistentHash(myServicesStorage, fromServer, status)
	case Random:
		addr = random(myServicesStorage, svrConfig.ServiceName)
	case WeightedRoundRobin:
		addr = weightedRoundRobin(myServicesStorage, svrConfig.ServiceName)
	case LeastConnections:
		addr = leastConnections(myServicesStorage, svrConfig.ServiceName)
	case WeightedLeastConnections:
		addr = weightedLeastConnections(myServicesStorage, svrConfig.ServiceName)
	case FastestResponse:
		addr = fastestResponse(myServicesStorage, svrConfig.ServiceName)
	}
	return addr
}
