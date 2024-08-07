package main

import (
	"MoreFun/etcd"
	pb "MoreFun/proto"
	"MoreFun/routerStrategy"
	"context"
	"flag"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"strconv"
)

// 定义命令行参数，用于指定 sidecar 的端口
var (
	ip   = flag.String("ip", "localhost", "The sidecar ip")
	port = flag.Int("port", 50051, "The sidecar port")
)

// MiniGameRouterServer 实现了 pb.UnimplementedMiniGameRouterServer 接口
type MiniGameRouterServer struct {
	pb.UnimplementedMiniGameRouterServer
}

// 全局变量，用于存储服务发现状态
var myIsDiscovered = &routerStrategy.IsDiscovered{
	IsDiscovered: map[string]bool{},
}

// 全局变量，用于存储服务信息
var myServicesStorage = &routerStrategy.ServicesStorage{
	ServicesStorage: make(map[string]map[string]string),
	CurrentWeight:   make(map[string]map[string]int),
	DynamicRouter:   make(map[string]string),
	RandomStorage: &routerStrategy.RandomStruct{
		RandomList: make([]string, 0),
		RandomMap:  make(map[string]int),
	},
	HashRing: treemap.NewWithIntComparator(),
}

// 监听服务名称的变化
func WatchServiceName(serviceName string) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	keepAliveCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watchChannel := cli.Watch(keepAliveCtx, serviceName, clientv3.WithPrefix())
	for watchRes := range watchChannel {
		for _, ev := range watchRes.Events {
			myServicesStorage.Lock()
			switch ev.Type {
			// 删除 myServicesStorage 中的键值对
			case clientv3.EventTypeDelete:
				delete(myServicesStorage.ServicesStorage[serviceName], string(ev.Kv.Key))
				if routerStrategy.ConsistentHash == routerStrategy.ServiceConfigs[serviceName].Strategy {
					myServicesStorage.DeleteNode(string(ev.Kv.Key))
				}
				if routerStrategy.Random == routerStrategy.ServiceConfigs[serviceName].Strategy {
					myServicesStorage.DeleteRandom(string(ev.Kv.Key))
				}
			// 添加 myServicesStorage 中的键值对
			case clientv3.EventTypePut:
				myServicesStorage.ServicesStorage[serviceName][string(ev.Kv.Key)] = string(ev.Kv.Value)
				if routerStrategy.ConsistentHash == routerStrategy.ServiceConfigs[serviceName].Strategy {
					myServicesStorage.AddNode(string(ev.Kv.Key), string(ev.Kv.Value))
				}
				if routerStrategy.Random == routerStrategy.ServiceConfigs[serviceName].Strategy {
					myServicesStorage.AddRandom(string(ev.Kv.Key), string(ev.Kv.Value))
				}
			}
			myServicesStorage.Unlock()
		}
	}
}
func WatchDynamicRouter(key string) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	keepAliveCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	watchChannel := cli.Watch(keepAliveCtx, key, clientv3.WithPrefix())
	for watchRes := range watchChannel {
		for _, ev := range watchRes.Events {
			switch ev.Type {
			// 删除 myServicesStorage 中的键值对
			case clientv3.EventTypeDelete:
				delete(myServicesStorage.DynamicRouter, string(ev.Kv.Key))
				return
			// 添加 myServicesStorage 中的键值对
			case clientv3.EventTypePut:
				myServicesStorage.DynamicRouter[string(ev.Kv.Key)] = string(ev.Kv.Value)
			}
		}
	}
}

// 注册服务
func (s *MiniGameRouterServer) RegisterService(ctx context.Context, req *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	var grantLease bool
	var leaseID clientv3.LeaseID

	serviceKey := fmt.Sprintf("%s_%s", req.Service.Name, req.Service.Port)
	serviceValue := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", req.Service.Name, req.Service.Ip, req.Service.Port, req.Service.Protocol, req.Service.Weight, req.Service.Status, req.Service.ConnNum)

	// 查看当前的服务是否注册过
	getRes, err := cli.Get(ctx, serviceKey, clientv3.WithCountOnly())
	if err != nil {
		log.Fatalf("Failed to get service key: %v", err)
	}

	// 没有注册过则创建一个租约，同时将租约的ID赋值给leaseID
	if getRes.Count == 0 {
		grantLease = true
		leaseRes, err := cli.Grant(ctx, 10)
		if err != nil {
			log.Fatalf("Failed to grant lease: %v", err)
		}
		leaseID = leaseRes.ID
	}

	// 创建一个kv客户端实现数据插入etcd
	kv := clientv3.NewKV(cli)

	// 开启事务
	txn := kv.Txn(ctx)
	// 判断数据库中是否存在 s *Service
	_, err = txn.If(clientv3.Compare(clientv3.CreateRevision(serviceKey), "=", 0)).
		Then(
			clientv3.OpPut(serviceKey, serviceValue, clientv3.WithLease(leaseID)),
		).
		Else(
			clientv3.OpPut(serviceKey, serviceValue, clientv3.WithIgnoreLease()),
		).
		Commit()

	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	// 租约保活机制
	if grantLease {
		go func() {
			cli := etcd.NewEtcdCli()
			defer cli.Close()

			keepAliveCtx, cancel := context.WithCancel(context.Background())
			defer cancel()

			leaseKeepAlive, err := cli.KeepAlive(keepAliveCtx, leaseID)
			if err != nil {
				log.Fatalf("Failed to keep lease alive: %v", err)
			}
			for range leaseKeepAlive {
				//fmt.Printf("leaseID:%x, ttl:%d\n", lease.ID, lease.TTL)
			}
		}()
	}

	return &pb.RegisterServiceResponse{
		Msg: serviceKey + "注册成功",
	}, nil
}

// 服务发现
func (s *MiniGameRouterServer) DiscoverService(ctx context.Context, req *pb.DiscoverServiceRequest) (*pb.DiscoverServiceResponse, error) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	var addr string
	if req.FixedRouterAddr == "" {
		if _, ok := myServicesStorage.DynamicRouter[req.ToMsg]; !ok {
			getRes, err := cli.Get(ctx, req.ToMsg)
			if err != nil {
				log.Fatalln(err)
			}
			if getRes.Count > 0 {
				myServicesStorage.Lock()
				//log.Printf("tsytsy%s,%s", req.ToMsg, string(getRes.Kvs[0].Value))
				myServicesStorage.DynamicRouter[req.ToMsg] = string(getRes.Kvs[0].Value)
				myServicesStorage.Unlock()
				go WatchDynamicRouter(req.ToMsg)
			}
		}
		if endPoint, ok := myServicesStorage.DynamicRouter[req.ToMsg]; ok {
			addr = endPoint
		} else if len(req.ToMsg) == 1 {
			// 如果服务未被发现过，则进行服务发现
			if !myIsDiscovered.IsDiscovered[req.ToMsg] {
				//log.Println("FIRST TIME")
				myIsDiscovered.IsDiscovered[req.ToMsg] = true
				myServicesStorage.ServicesStorage[req.ToMsg] = make(map[string]string)
				getRes, err := cli.Get(ctx, req.ToMsg, clientv3.WithPrefix())
				if err != nil {
					log.Fatalln(err)
				}

				if getRes.Count > 0 {
					myServicesStorage.Lock()
					for _, kv := range getRes.Kvs {
						myServicesStorage.ServicesStorage[req.ToMsg][string(kv.Key)] = string(kv.Value)
						if routerStrategy.ConsistentHash == routerStrategy.ServiceConfigs[req.ToMsg].Strategy {
							myServicesStorage.AddNode(string(kv.Key), string(kv.Value))
						}
						if routerStrategy.Random == routerStrategy.ServiceConfigs[req.ToMsg].Strategy {
							myServicesStorage.AddRandom(string(kv.Key), string(kv.Value))
						}
					}
					myServicesStorage.Unlock()
				}
				go WatchServiceName(req.ToMsg)
			}

			// 打印 req.ToMsg 类的服务
			var index int = 0
			for key, value := range myServicesStorage.ServicesStorage[req.ToMsg] {
				index += 1
				log.Printf("%d_%s_%s", index, key, value)
			}

			// 获取服务配置
			config, ok := routerStrategy.ServiceConfigs[req.ToMsg]
			if !ok {
				return nil, fmt.Errorf("no configuration found for service %s", req.ToMsg)
			}
			addr = routerStrategy.GetAddr(myServicesStorage, config, req.FromMsg, req.Status)
		}
	} else {
		addr = req.FixedRouterAddr
	}
	//log.Printf("%s,  %s", req.FixedRouterAddr, addr)
	// 连接另外一个 sidecar
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to another sidecar: %s", err)
	}
	defer conn.Close()

	c := pb.NewMiniGameRouterClient(conn)

	helloReq := &pb.HelloRequest{
		Msg: req.FromMsg,
	}
	helloRes, err := c.SayHello(ctx, helloReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to grpc another sidecar: %s", err)
	}
	return &pb.DiscoverServiceResponse{
		Msg: helloRes.Msg,
	}, nil
}
func (s *MiniGameRouterServer) SetCustomRoute(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	var leaseID clientv3.LeaseID
	timeInt, _ := strconv.Atoi(req.Timeout)
	leaseRes, err := cli.Grant(ctx, int64(timeInt))
	if err != nil {
		log.Fatalf("Failed to grant lease: %v", err)
	}
	leaseID = leaseRes.ID
	_, err = cli.Put(ctx, req.Key, req.Value, clientv3.WithLease(leaseID))
	if err != nil {
		return nil, fmt.Errorf("Fail SetCustomRoute")
	}
	return &pb.SetResponse{
		Msg: fmt.Sprintf("%s_%s set success", req.Key, req.Value),
	}, nil
}

// 简单的问候服务
func (s *MiniGameRouterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	addr := fmt.Sprintf("%s:%d", *ip, *port-1)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Sidecar could not connect to its server %s", err)
	}
	defer conn.Close()

	c := pb.NewMiniGameRouterClient(conn)
	helloReq := &pb.HelloRequest{
		Msg: fmt.Sprintf("Hello, I am %s", req.Msg),
	}
	rRes, err := c.SayHello(ctx, helloReq)
	if err != nil {
		return nil, fmt.Errorf("Failed to grpc its server: %s", err)
	}
	return &pb.HelloResponse{
		Msg: rRes.Msg,
	}, nil
}

// 主函数，启动 gRPC 服务器
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMiniGameRouterServer(s, &MiniGameRouterServer{})
	log.Printf("Sidecar is listening on port %d", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
