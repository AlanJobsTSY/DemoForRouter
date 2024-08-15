package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// DialTimeout 连接etcd服务器的超时时间
const DialTimeout = time.Second * 20

// GetEtcdEndpoints etcd集群所有节点的ip列表
func GetEtcdEndpoints() []string {
	return []string{"9.135.119.71:2379", "9.135.119.71:12379", "9.135.119.71:22379"}
}

// NewEtcdCli new一个etcd客户端
func NewEtcdCli() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   GetEtcdEndpoints(),
		DialTimeout: DialTimeout,
	})
	if err != nil {
		log.Fatalf("ETCD Error: %v\n", err)
	}
	return cli
}
