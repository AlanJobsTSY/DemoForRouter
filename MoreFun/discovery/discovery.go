package discovery

import (
	"MoreFun/etcd"
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
)

// 服务的结构体
type Service struct {
	Name     string
	IP       string
	Port     string
	Protocol string
}

// 服务注册
func ServiceRegister(s *Service) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()

	var grantLease bool
	var leaseID clientv3.LeaseID

	//serviceKey := fmt.Sprintf("%s:%s", s.Name, s.Port)

	//查看当前的服务是否注册过
	getRes, err := cli.Get(context.Background(), s.Name, clientv3.WithCountOnly())
	if err != nil {
		log.Fatalln(err)
	}

	//没有注册过则创建一个租约，同时将租约的ID赋值给leaseID
	if getRes.Count == 0 {
		grantLease = true
		leaseRes, err := cli.Grant(context.Background(), 10)
		if err != nil {
			log.Fatalln(err)
		}
		leaseID = leaseRes.ID
	}

	//创建一个kv客户端实现数据插入etcd
	kv := clientv3.NewKV(cli)

	//开启事务
	txn := kv.Txn(context.Background())
	//判断数据库中是否存在 s *Service
	_, err = txn.If(clientv3.Compare(clientv3.CreateRevision(s.Name), "=", 0)).
		Then(
			clientv3.OpPut(s.Name, s.Name, clientv3.WithLease(leaseID)),
			clientv3.OpPut(s.Name+".IP", s.IP, clientv3.WithLease(leaseID)),
			clientv3.OpPut(s.Name+".Port", s.Port, clientv3.WithLease(leaseID)),
			clientv3.OpPut(s.Name+".Protocol", s.Protocol, clientv3.WithLease(leaseID)),
		).
		Else(
			clientv3.OpPut(s.Name, s.Name, clientv3.WithIgnoreLease()),
			clientv3.OpPut(s.Name+".IP", s.IP, clientv3.WithIgnoreLease()),
			clientv3.OpPut(s.Name+".Port", s.Port, clientv3.WithIgnoreLease()),
			clientv3.OpPut(s.Name+".Protocol", s.Protocol, clientv3.WithIgnoreLease()),
		).
		Commit()

	if err != nil {
		log.Fatalln(err)
	}

	//租约保活机制
	if grantLease {
		leaseKeepAlive, err := cli.KeepAlive(context.Background(), leaseID)
		if err != nil {
			log.Fatalln(err)
		}
		for lease := range leaseKeepAlive {
			fmt.Printf("leaseID:%x,ttl:%d\n", lease.ID, lease.TTL)
		}
	}

}

// 存在etcd的信息在本地也会存储一份
type Services struct {
	services map[string]*Service
	sync.RWMutex
}

// 实例化Services
var myServices = &Services{
	services: map[string]*Service{},
}

// 服务发现
func ServiceDiscovery(svcName string) *Service {
	var s *Service = nil
	myServices.RLock()
	s, _ = myServices.services[svcName]
	myServices.RUnlock()
	return s
}

// 将服务加入本地
func WatchServiceName(svcName string) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()

	getRes, err := cli.Get(context.Background(), svcName, clientv3.WithPrefix())
	if err != nil {
		log.Fatalln(err)
	}
	if getRes.Count > 0 {

		s := kvsToService(svcName, getRes.Kvs)
		myServices.Lock()
		myServices.services[svcName] = s
		myServices.Unlock()
	}

	//监听以 svcName 为前缀的键值对变化
	watchChannel := cli.Watch(context.Background(), svcName, clientv3.WithPrefix())
	for watchRes := range watchChannel {
		for _, ev := range watchRes.Events {
			//删除myServices中的键值对
			if ev.Type == clientv3.EventTypeDelete {
				myServices.Lock()
				delete(myServices.services, svcName)
				myServices.Unlock()
			}
			//添加myServices中的键值对
			if ev.Type == clientv3.EventTypePut {
				myServices.Lock()
				//若无svcName的key，则要创建一个空的对象
				if _, ok := myServices.services[svcName]; !ok {
					myServices.services[svcName] = &Service{}
				}
				switch string(ev.Kv.Key) {
				case svcName:
					myServices.services[svcName].Name = string(ev.Kv.Value)
				case svcName + ".IP":
					myServices.services[svcName].IP = string(ev.Kv.Value)
				case svcName + ".Port":
					myServices.services[svcName].Port = string(ev.Kv.Value)
				case svcName + ".Protocol":
					myServices.services[svcName].Protocol = string(ev.Kv.Value)
				}
				myServices.Unlock()
			}
		}

	}
}

func kvsToService(svcName string, kvs []*mvccpb.KeyValue) *Service {
	s := &Service{}
	for _, kv := range kvs {
		switch string(kv.Key) {
		case svcName:
			s.Name = string(kv.Value)
		case svcName + ".IP":
			s.IP = string(kv.Value)
		case svcName + ".Port":
			s.Port = string(kv.Value)
		case svcName + ".Protocol":
			s.Protocol = string(kv.Value)
		}
	}
	return s
}
