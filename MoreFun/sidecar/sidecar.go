package main

import (
	"MoreFun/etcd"
	pb "MoreFun/proto"
	"MoreFun/routerStrategy"
	"context"
	"flag"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "The sidecar port")
)

type MiniGameRouterServer struct {
	pb.UnimplementedMiniGameRouterServer
}

var myIsDiscovered = &routerStrategy.IsDiscovered{
	IsDiscovered: map[string]bool{},
}

var myServicesStorage = &routerStrategy.ServicesStorage{
	ServicesStorage: map[string]map[string]string{},
	CurrentWeight:   make(map[string]map[string]int),
}

func WatchServiceName(serviceName string) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	keepAliveCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watchChannel := cli.Watch(keepAliveCtx, serviceName, clientv3.WithPrefix())
	for watchRes := range watchChannel {
		//log.Println("working?????????")
		for _, ev := range watchRes.Events {
			//log.Printf("%s__%s__%s__%s", ev.Type, serviceName, string(ev.Kv.Key), string(ev.Kv.Value))
			myServicesStorage.Lock()
			switch ev.Type {
			//删除myServicesStorage中的键值对
			case clientv3.EventTypeDelete:
				delete(myServicesStorage.ServicesStorage[serviceName], string(ev.Kv.Key))
			//添加myServicesStorage中的键值对
			case clientv3.EventTypePut:
				myServicesStorage.ServicesStorage[serviceName][string(ev.Kv.Key)] = string(ev.Kv.Value)
			}
			myServicesStorage.Unlock()
		}

	}
}
func (s *MiniGameRouterServer) RegisterService(ctx context.Context, req *pb.RegisterServiceRequest) (*pb.RegisterServiceResponse, error) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	var grantLease bool
	var leaseID clientv3.LeaseID

	serviceKey := fmt.Sprintf("%s_%s", req.Service.Name, req.Service.Port)
	serviceValue := fmt.Sprintf("%s:%s:%s:%s:%s", req.Service.Name, req.Service.Ip, req.Service.Port, req.Service.Protocol, req.Service.Weight)

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

func (s *MiniGameRouterServer) DiscoverService(ctx context.Context, req *pb.DiscoverServiceRequest) (*pb.DiscoverServiceResponse, error) {
	//服务发现
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	//路由选择 needtodo
	if !myIsDiscovered.IsDiscovered[req.ToMsg] {
		log.Println("FIRST TIME")
		//routerStrategy.MyIsDiscovered.Lock()
		myIsDiscovered.IsDiscovered[req.ToMsg] = true
		//routerStrategy.MyIsDiscovered.Unlock()
		myServicesStorage.ServicesStorage[req.ToMsg] = make(map[string]string)
		getRes, err := cli.Get(ctx, req.ToMsg, clientv3.WithPrefix())
		if err != nil {
			log.Fatalln(err)
		}

		if getRes.Count > 0 {
			myServicesStorage.Lock()
			for _, kv := range getRes.Kvs {
				//log.Printf("%s_%s_%s", req.ToMsg, string(kv.Key), string(kv.Value))
				myServicesStorage.ServicesStorage[req.ToMsg][string(kv.Key)] = string(kv.Value)
			}
			myServicesStorage.Unlock()
		}
		go WatchServiceName(req.ToMsg)
	}
	//打印req.ToMsg类的服务
	var index int = 0
	for key, value := range myServicesStorage.ServicesStorage[req.ToMsg] {
		index += 1
		log.Printf("%d_%s_%s", index, key, value)
	}
	/*
		getRes, err := cli.Get(ctx, req.ToMsg, clientv3.WithPrefix())
		if err != nil {
			log.Fatalln(err)
		}
		if len(getRes.Kvs) == 0 {
			return nil, fmt.Errorf("no key found for %s", req.ToMsg)
		}
		value := string(getRes.Kvs[0].Value)
		parts := strings.Split(value, ":")
		addr := fmt.Sprintf("%s:%s", parts[1], parts[2])
	*/

	config, ok := routerStrategy.ServiceConfigs[req.ToMsg]
	if !ok {
		return nil, fmt.Errorf("no configuration found for service %s", req.ToMsg)
	}
	addr := routerStrategy.GetAddr(myServicesStorage, config)
	//log.Printf("tsy__addr:%s", addr)
	//连接另外一个sidecar
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connnect another sidecar: %s", err)
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
func (s *MiniGameRouterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	addr := fmt.Sprintf("localhost:%d", *port-1)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Sidecar could not connect its server %s", err)
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
