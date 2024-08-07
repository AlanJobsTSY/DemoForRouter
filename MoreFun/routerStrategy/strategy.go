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
		/*
			case FixedRoute:
				addr = fixedRoute(myServicesStorage, svrConfig.ServiceName, svrConfig.FixedRoute)
		*/
	}
	return addr
}

/*
// fixedRoute 固定路由策略
func fixedRoute(myServicesStorage *ServicesStorage, svrName string, fixedRoute string) string {
	return fixedRoute
}
*/
