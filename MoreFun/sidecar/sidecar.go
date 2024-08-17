package main

import (
	"MoreFun/etcd"
	pb "MoreFun/proto"
	"MoreFun/routerStrategy"
	"context"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/emirpasic/gods/maps/treemap"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"log"
	"net"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
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

// 缓存，用于存储服务发现状态
var myIsDiscovered = &routerStrategy.IsDiscovered{
	IsDiscovered: map[string]bool{},
}

// 缓存，用于存储服务信息
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

// 监听动态键值路由
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
	//var getRes *clientv3.GetResponse
	//var err error
	serviceKey := fmt.Sprintf("%s_%s:%s", req.Service.Name, req.Service.Ip, req.Service.Port)
	serviceValue := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", req.Service.Name, req.Service.Ip, req.Service.Port, req.Service.Protocol, req.Service.Weight, req.Service.Status, req.Service.ConnNum)

	// 查看当前的服务是否注册过
	/*
		for i := 0; i < 20; i++ {
			getRes, err = cli.Get(ctx, serviceKey, clientv3.WithCountOnly())
			if err != nil {
				log.Printf("Failed to get service key: %v", err)
				time.Sleep(time.Second)
				continue
			}
			break
		}*/
	// 没有注册过则创建一个租约，同时将租约的ID赋值给leaseID
	//if getRes.Count == 0 {
	grantLease = true
	leaseRes, err := cli.Grant(ctx, 10)
	if err != nil {
		log.Fatalf("Failed to grant lease: %v", err)
	}
	time.Sleep(5 * time.Second)
	leaseID = leaseRes.ID
	//}
	/*
		//log.Printf("%s", req.NsIp)
		conn, err := grpc.Dial(req.NsIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to another sidecar: %s", err)
		}
		defer conn.Close()
		c := pb.NewMiniGameRouterClient(conn)
		writeReq := &pb.WriteRequest{
			SvrKey:   serviceKey,
			SvrValue: serviceValue,
			IsLease:  grantLease,
			LeaseID:  int64(leaseID),
		}
		// 让另一个sidecar调用自己负责的服务
		_, err = c.WriteServers(ctx, writeReq)
		if err != nil {
			log.Printf("Write err: %s", err)
		}
		//time.Sleep(5 * time.Second)
	*/

	/*

		// 创建一个kv客户端实现数据插入etcd
		for i := 0; i < 20; i++ {
			if grantLease == true {
				_, err := cli.Put(context.Background(), serviceKey, serviceValue, clientv3.WithLease(leaseID))
				if err != nil {
					log.Printf("Failed to registe withlease: %v", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}
			} else {
				_, err := cli.Put(context.Background(), serviceKey, serviceValue, clientv3.WithIgnoreLease())
				if err != nil {
					log.Printf("Failed to registe withIgnoreLease: %v", err)
					time.Sleep(100 * time.Millisecond)
					continue
				}
			}

			break
		}

	*/
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

	// 创建一个新的Kafka生产者实例，配置Kafka服务器地址
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close() // 在函数结束时关闭生产者
	msgR := &pb.RegisterServiceResponse{
		SvrKey:   serviceKey,
		SvrValue: serviceValue,
		IsLease:  grantLease,
		LeaseID:  int64(leaseID),
	}
	topic := "myTopic" // 定义要发送消息的Kafka主题
	msgBytes, err := proto.Marshal(msgR)
	if err != nil {
		log.Fatalf("Failed to marshal message: %s", err)
	}
	// 创建Kafka消息
	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, // 设置消息的主题和分区
		Value:          msgBytes,                                                           // 设置消息的内容
	}
	// 发送消息到Kafka
	err = p.Produce(&msg, nil)
	if err != nil {
		log.Fatalf("Failed to produce message: %s", err)
	}

	// 等待消息传递完成，最多等待15秒
	p.Flush(10 * 1000)
	//fmt.Println("Message sent successfully")

	return msgR, nil
}

// 设置动态键值
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

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close() // 在函数结束时关闭生产者
	msgR := &pb.RegisterServiceResponse{
		SvrKey:   req.Key,
		SvrValue: req.Value,
		IsLease:  true,
		LeaseID:  int64(leaseID),
	}
	topic := "myTopic" // 定义要发送消息的Kafka主题
	msgBytes, err := proto.Marshal(msgR)
	if err != nil {
		log.Fatalf("Failed to marshal message: %s", err)
	}
	// 创建Kafka消息
	msg := kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, // 设置消息的主题和分区
		Value:          msgBytes,                                                           // 设置消息的内容
	}
	// 发送消息到Kafka
	err = p.Produce(&msg, nil)
	if err != nil {
		log.Fatalf("Failed to produce message: %s", err)
	}

	// 等待消息传递完成，最多等待15秒
	p.Flush(10 * 1000)

	return &pb.SetResponse{
		Msg: fmt.Sprintf("%s_%s sent successfully", req.Key, req.Value),
	}, nil
}

// 打印svrWant类服务的列表供用户选择
func printWant(svrWant string) {
	var index int = 0
	for _, value := range myServicesStorage.ServicesStorage[svrWant] {
		index += 1
		parts := strings.Split(value, ":")
		partPort, _ := strconv.Atoi(parts[2])
		parts[2] = strconv.Itoa(partPort - 1)
		newValue := strings.Join(parts, ":")
		log.Printf("[%d] %s", index, newValue)
	}
}

// 缓存初始化
func storageInit(ctx context.Context, cli *clientv3.Client, reqToMsg string) {
	myIsDiscovered.Lock()
	defer myIsDiscovered.Unlock()

	if !myIsDiscovered.IsDiscovered[reqToMsg] {
		//log.Println("FIRST TIME")
		myIsDiscovered.IsDiscovered[reqToMsg] = true
		myServicesStorage.ServicesStorage[reqToMsg] = make(map[string]string)
		getRes, err := cli.Get(ctx, reqToMsg, clientv3.WithPrefix())
		if err != nil {
			log.Fatalln(err)
		}

		if getRes.Count > 0 {
			myServicesStorage.Lock()
			for _, kv := range getRes.Kvs {
				myServicesStorage.ServicesStorage[reqToMsg][string(kv.Key)] = string(kv.Value)
				if routerStrategy.ConsistentHash == routerStrategy.ServiceConfigs[reqToMsg].Strategy {
					myServicesStorage.AddNode(string(kv.Key), string(kv.Value))
				}
				if routerStrategy.Random == routerStrategy.ServiceConfigs[reqToMsg].Strategy {
					myServicesStorage.AddRandom(string(kv.Key), string(kv.Value))
				}
			}
			myServicesStorage.Unlock()
		}
		// 开启监听
		go WatchServiceName(reqToMsg)
	}
}

// 服务发现
func (s *MiniGameRouterServer) DiscoverService(ctx context.Context, req *pb.DiscoverServiceRequest) (*pb.DiscoverServiceResponse, error) {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	var addr string
	var svrDiscoverTime int64
	var startTime time.Time
	var endTime time.Time
	if req.FixedRouterAddr == "" {
		// 拉取动态键值
		if _, ok := myServicesStorage.DynamicRouter[req.ToMsg]; !ok {
			getRes, err := cli.Get(ctx, req.ToMsg)
			if err != nil {
				log.Fatalln(err)
			}
			if getRes.Count > 0 {
				myServicesStorage.Lock()
				myServicesStorage.DynamicRouter[req.ToMsg] = string(getRes.Kvs[0].Value)
				myServicesStorage.Unlock()
				go WatchDynamicRouter(req.ToMsg)
			}
		}
		if endPoint, ok := myServicesStorage.DynamicRouter[req.ToMsg]; ok { //动态键值路由
			addr = endPoint
		} else if len(req.ToMsg) == 1 { //路由选择算法
			// 如果服务未被发现过，则进行服务发现
			storageInit(ctx, cli, req.ToMsg)
			// 打印 req.ToMsg 类的服务
			/*
				printWant(req.ToMsg)
			*/
			// 获取服务配置
			config, ok := routerStrategy.ServiceConfigs[req.ToMsg]
			if !ok {
				return nil, fmt.Errorf("no configuration found for service %s", req.ToMsg)
			}
			// 记录开始时间
			startTime = time.Now()
			//路由选择
			addr = routerStrategy.GetAddr(myServicesStorage, config, req.FromMsg, req.Status)
			// 记录结束时间
			endTime = time.Now()
			// 计算执行时间（以毫秒为单位）
			svrDiscoverTime = endTime.Sub(startTime).Milliseconds()
		}
	} else { // 指定目标路由
		// 如果服务未被发现过，则进行服务发现
		storageInit(ctx, cli, req.ToMsg)
		addr = req.FixedRouterAddr
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
		if err != nil {
			log.Printf("The specified target cannot be reached. Please select an IP address from the following %s services.", req.ToMsg)
			printWant(req.ToMsg)
			return nil, err
		}
		conn.Close()
	}
	// 连接另外一个 sidecar
	startTime = time.Now()
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to another sidecar: %s", err)
	}
	defer conn.Close()
	c := pb.NewMiniGameRouterClient(conn)
	helloReq := &pb.HelloRequest{
		Msg: req.FromMsg,
	}
	endTime = time.Now()
	dialTime := endTime.Sub(startTime).Milliseconds()

	// 让另一个sidecar调用自己负责的服务
	startTime = time.Now()
	helloRes, err := c.SayHello(ctx, helloReq)
	endTime = time.Now()
	returnTime := endTime.Sub(startTime).Milliseconds()
	if err != nil {
		return nil, fmt.Errorf("Failed to grpc another sidecar: %s", err)
	}
	return &pb.DiscoverServiceResponse{
		Msg:             helloRes.Msg,
		SvrDiscoverTime: svrDiscoverTime,
		DialTime:        dialTime,
		ReturnTime:      returnTime,
	}, nil
}

// sidecar调用服务端的问候服务
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
	//test
	// 创建 CPU 分析文件
	f, err := os.Create("sidecar.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	// 开始 CPU 分析
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	flag.Parse()
	var lis net.Listener
	for i := 0; i < 20; i++ {
		lis, err = net.Listen("tcp", fmt.Sprintf(":%d", *port))
		if err != nil {
			log.Printf("Failed to listen: %v", err)
			time.Sleep(time.Second)
			continue
		}
		break
	}
	if lis != nil {
		defer lis.Close()
	}
	if err != nil {
		return
	}
	s := grpc.NewServer()
	pb.RegisterMiniGameRouterServer(s, &MiniGameRouterServer{})
	log.Printf("Sidecar is listening on port %d", *port)
	if err := s.Serve(lis); err != nil {
		log.Printf("Failed to serve: %v", err)
	}
}
