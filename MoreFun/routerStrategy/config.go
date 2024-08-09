package routerStrategy

import (
	"github.com/emirpasic/gods/maps/treemap"
	"sync"
)

// RoutingStrategy 定义了路由策略的类型
type RoutingStrategy string

// 定义各种路由策略的常量
const (
	ConsistentHash           RoutingStrategy = "consistent_hash"
	Random                   RoutingStrategy = "random"
	WeightedRoundRobin       RoutingStrategy = "weighted_round_robin"
	LeastConnections         RoutingStrategy = "least_connections"
	WeightedLeastConnections RoutingStrategy = "weighted_least_connections"
	FastestResponse          RoutingStrategy = "fastest_response"
	//FixedRoute         RoutingStrategy = "fixed_route"
)

// ServiceConfig 定义了服务的配置，包括服务名称、路由策略和固定路由（仅在 FixedRoute 策略下使用）
type ServiceConfig struct {
	ServiceName string
	Strategy    RoutingStrategy
	//Status      bool // 是否是有状态服务
}

// ServiceConfigs 是一个全局变量，存储了所有服务的配置
var ServiceConfigs = map[string]*ServiceConfig{
	"A": {ServiceName: "A", Strategy: ConsistentHash},
	"B": {ServiceName: "B", Strategy: Random},
	"C": {ServiceName: "C", Strategy: WeightedRoundRobin},
	"D": {ServiceName: "D", Strategy: LeastConnections},
	"E": {ServiceName: "E", Strategy: WeightedLeastConnections},
	"F": {ServiceName: "F", Strategy: FastestResponse},
	//"D": {ServiceName: "D", Strategy: FixedRoute},
}

// IsDiscovered 用于存储服务是否已被发现的状态
type IsDiscovered struct {
	IsDiscovered map[string]bool
	sync.RWMutex
}

// ServicesStorage 用于存储服务实例信息和当前权重
type ServicesStorage struct {
	ServicesStorage map[string]map[string]string // 存储服务实例信息
	CurrentWeight   map[string]map[string]int    // 存储当前权重
	DynamicRouter   map[string]string
	HashRing        *treemap.Map
	sync.RWMutex
}
