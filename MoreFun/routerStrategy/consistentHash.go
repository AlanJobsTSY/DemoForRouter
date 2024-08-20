package routerStrategy

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
	"sync"
)

const VirtualNodeNoPerNode = 200

type Node struct {
	ip          string
	serverSlice []string //管理的服务
	sync.RWMutex
}

func newNode(ip string) *Node {
	return &Node{
		ip: ip,
	}
}

// 计算哈希值
func hashCode(virtualNodeKey string) int {
	h := fnv.New32a()
	h.Write([]byte(fmt.Sprintf("%v", virtualNodeKey)))
	return int(h.Sum32())
}

// 插入节点
func (ss *ServicesStorage) AddNode(svrName string, ip string) {

	parts := strings.Split(ip, ":")
	ip = fmt.Sprintf("%s:%s", parts[1], parts[2])
	//log.Printf("????  %s", ip)
	node := newNode(ip)
	for i := 0; i < VirtualNodeNoPerNode; i++ {
		virtualKey := hashCode(svrName + strconv.Itoa(i))
		ss.HashRing.Put(virtualKey, node)
	}
	//有状态迁移
}

// 删除节点
func (ss *ServicesStorage) DeleteNode(svrName string) {
	for i := 0; i < VirtualNodeNoPerNode; i++ {
		virtualKey := hashCode(svrName + strconv.Itoa(i))
		ss.HashRing.Remove(virtualKey)
	}
}

// 查询节点
func (ss *ServicesStorage) getNode(svrName string) string {
	virtualKey := hashCode(svrName)
	_, node := ss.HashRing.Ceiling(virtualKey)
	if node == nil {
		_, node = ss.HashRing.Ceiling(0)
		if node == nil {
			return ""
		}
	}
	return node.(*Node).ip
	//有状态插入
}

// consistentHash 一致性哈希路由策略
func consistentHash(myServicesStorage *ServicesStorage, svrName string, status string) string {
	var addr string = ""
	addr = myServicesStorage.getNode(svrName)
	//log.Printf("why %s", addr)
	return addr
}
