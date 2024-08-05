package routerStrategy

import "sync"

type RoutingStrategy string

const (
	ConsistentHash     RoutingStrategy = "consistent_hash"
	Random             RoutingStrategy = "random"
	WeightedRoundRobin RoutingStrategy = "weighted_round_robin"
	FixedRoute         RoutingStrategy = "fixed_route"
)

type ServiceConfig struct {
	ServiceName string
	Strategy    RoutingStrategy
	FixedRoute  string // 仅在 FixedRoute 策略下使用
}

var ServiceConfigs = map[string]*ServiceConfig{
	"A": {ServiceName: "A", Strategy: ConsistentHash},
	"B": {ServiceName: "B", Strategy: Random},
	"C": {ServiceName: "C", Strategy: WeightedRoundRobin},
	"D": {ServiceName: "D", Strategy: FixedRoute, FixedRoute: "localhost:55555"},
}

type IsDiscovered struct {
	IsDiscovered map[string]bool
	sync.RWMutex
}

type ServicesStorage struct {
	ServicesStorage map[string]map[string]string
	CurrentWeight   map[string]map[string]int
	sync.RWMutex
}
