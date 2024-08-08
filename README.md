# MiniGameRouter  

MiniGameRouter是一个基于 gRPC、Etcd、sidecar架构的微服务路由器示例项目。该项目展示了如何使用 gRPC 实现服务注册、服务发现以及简单的服务调用。

## 背景

在微服务架构中，服务注册和发现是非常重要的功能。MiniGameRouter 项目展示了如何使用 gRPC 实现这些功能，并提供了一个简单的示例服务。

## 体验流程

### 1.在个人服务器使用docker部署一个Etcd集群

首先下载Etcd

``` bash
sudo docker pull quay.io/coreos/etcd:v3.5.5
```

然后修改成自己的部署路径和配置路径（以下为参考），配置文件可参考etcd0.yaml、etcd1.yaml、etcd2.yaml（主要修改ip地址）

``` bash
docker run -d -p 2379:2379 -p 2380:2380 -v /tmp/etcd0-data:/etcd-data -v /data/home/shenytong/workspace/etcdconf:/etcd-conf --name etcd0 quay.io/coreos/etcd:v3.5.5 /usr/local/bin/etcd --config-file=/etcd-conf/etcd0.yaml
docker run -d -p 12379:12379 -p 12380:12380 -v /tmp/etcd1-data:/etcd-data -v /data/home/shenytong/workspace/etcdconf:/etcd-conf --name etcd1 quay.io/coreos/etcd:v3.5.5 /usr/local/bin/etcd --config-file=/etcd-conf/etcd1.yaml
docker run -d -p 22379:22379 -p 22380:22380 -v /tmp/etcd2-data:/etcd-data -v /data/home/shenytong/workspace/etcdconf:/etcd-conf --name etcd2 quay.io/coreos/etcd:v3.5.5 /usr/local/bin/etcd --config-file=/etcd-conf/etcd2.yaml
```

进入docker bash中

``` bash
docker exec -it etcd0 bash
```

查看集群状态

``` bash
etcdctl endpoint status --cluster -w table
```

![etcd集群](.\pic\etcd集群.png)

### 2. 启动服务

可能需要安装最新的环境

``` bash
 go get google.golang.org/grpc@latest
```

只支持`A`,`B`,`C`,`D`作为服务名字，同`ip`的`port`需要间隔两个启动（如`--port 50050` 下一个需要`--port 50052`）

``` bash
go run .\server\server.go --name A  --ip "服务器地址,默认localhost" --port 50050 --weight 1 
```

### 3. 服务请求调用

键入`A`,`B`,`C`,`D`,`[key]`请求各类服务

``` bash
A "键入A则使用一致性哈希算法访问A类服务器"
B "键入B则使用随机算法访问B类服务器"
C "键入C则使用平滑权重轮询算法访问C类服务器"
D "键入D则使用最少连接数路由算法访问D类服务器"
[key] "键入[key]则实现动态键值路由访问[endpoint]服务器（timeout时间内）"
```

键入 `[任意字符串]  ip:port`实现固定键值路由

``` bash
[svrName] localhost:50050 "前者任意字符串，后者为点对点路由地址"
```

键入`set`然后键入`key endpoints timeout`设定动态键值路由

``` bash
set

[key] [endpoints] timeout "如 `tsy localhost:50050 60`,60秒为动态键值持续时间，可被覆盖"
```

## 项目细节

```bash
MoreFun
├── etcd                    # etcd相关代码
│   └── etcd.go             	# etcd 相关配置
├── go.mod                  
├── go.sum                 
├── proto                   # gRPC 协议文件和生成的代码
│   ├── minigame_router_grpc.pb.go  # gRPC 服务代码
│   ├── minigame_router.pb.go       # Protocol Buffers 代码
│   └── minigame_router.proto       # Protocol Buffers 定义文件
├── routerStrategy          # 路由策略相关代码
│   ├── config.go           	# 路由策略配置、存储数据结构
│   └── strategy.go         	# 路由策略实现
├── server                  # 服务器相关代码
│   └── server.go           	# 服务器实现
└── sidecar                 # Sidecar 相关代码
	└── sidecar.go          	# Sidecar 实现
```

### 1. minigame_router.proto

1. **文件头部**：
   - `syntax = "proto3";`：指定使用 Protobuf 的版本 3。
   - `option go_package = "MoreFun/proto";`：指定生成的 Go 代码的包路径。
   - `package proto;`：定义 proto 文件的包名。
2. **服务定义**：
   - `service MiniGameRouter`：定义一个名为 `MiniGameRouter` 的服务。
   - `rpc RegisterService(RegisterServiceRequest) returns (RegisterServiceResponse);`：定义一个 `RegisterService` RPC 方法，接收 `RegisterServiceRequest` 请求并返回 `RegisterServiceResponse` 响应。
   - `rpc DiscoverService(DiscoverServiceRequest) returns (DiscoverServiceResponse);`：定义一个 `DiscoverService` RPC 方法，接收 `DiscoverServiceRequest` 请求并返回 `DiscoverServiceResponse` 响应。
   - `rpc SayHello(HelloRequest) returns (HelloResponse);`：定义一个 `SayHello` RPC 方法，接收 `HelloRequest` 请求并返回 `HelloResponse` 响应。
3. **消息定义**：
   - `message Service`：定义一个 `Service` 消息，包含服务的基本信息。
   - `message RegisterServiceRequest`：定义一个 `RegisterServiceRequest` 消息，包含一个 `Service` 类型的字段。
   - `message RegisterServiceResponse`：定义一个 `RegisterServiceResponse` 消息，包含一个字符串类型的字段 `msg`。
   - `message DiscoverServiceRequest`：定义一个 `DiscoverServiceRequest` 消息，包含两个字符串类型的字段 `from_msg` 和 `to_msg`。
   - `message DiscoverServiceResponse`：定义一个 `DiscoverServiceResponse` 消息，包含一个字符串类型的字段 `msg`。
   - `message HelloRequest`：定义一个 `HelloRequest` 消息，包含一个字符串类型的字段 `msg`。
   - `message HelloResponse`：定义一个 `HelloResponse` 消息，包含一个字符串类型的字段 `msg`。

### 2. config.go

1. **RoutingStrategy 类型**：
   - 定义了路由策略的类型 `RoutingStrategy`，使用字符串表示。
2. **路由策略常量**：
   - `ConsistentHash`：一致性哈希路由策略。
   - `Random`：随机路由策略。
   - `WeightedRoundRobin`：加权轮询路由策略。
   - `FixedRoute`：固定路由策略。
3. **ServiceConfig 结构体**：
   - `ServiceName`：服务名称。
   - `Strategy`：路由策略。
   - `FixedRoute`：固定路由地址，仅在 `FixedRoute` 策略下使用。
4. **ServiceConfigs 变量**：
   - 存储了所有服务的配置，每个服务都有一个 `ServiceConfig` 实例。
5. **IsDiscovered 结构体**：
   - `IsDiscovered`：存储服务是否已被发现的状态。
   - `sync.RWMutex`：读写锁，用于并发控制。
6. **ServicesStorage 结构体**：
   - `ServicesStorage`：存储服务实例信息，使用嵌套的 map 结构。
   - `CurrentWeight`：存储当前权重，使用嵌套的 map 结构。
   - `sync.RWMutex`：读写锁，用于并发控制。

### 3. strategy.go

1. **GetAddr 函数**：
   - 根据服务配置选择路由策略，返回服务地址。
   - 支持的路由策略包括：一致性哈希、随机、加权轮询和固定路由。
2. **consistentHash 函数**：
   - 一致性哈希路由策略（未实现）。
3. **random 函数**：
   - 随机路由策略，从服务实例中随机选择一个实例地址。
4. **weightedRoundRobin 函数**：
   - 加权轮询路由策略，根据实例的权重进行轮询选择。
   - 计算每个实例的当前权重，选择当前权重最大的实例，并减去总权重。
5. **fixedRoute 函数**：
   - 固定路由策略，直接返回固定的路由地址。

### 4. server.go

1. **命令行参数**：
   - `name`：服务器名称。
   - `ip`：服务器 IP 地址。
   - `port`：服务器端口。
   - `protocol`：服务器协议。
   - `weight`：服务器权重。
2. **MiniGameRouterServer 结构体**：
   - 实现了 `MiniGameRouter` gRPC 服务。
3. **SayHello 方法**：
   - 实现了 gRPC 的 `SayHello` 方法，接收请求并返回响应。
4. **DiscoverService 函数**：
   - 发现服务，发送 `DiscoverServiceRequest` 请求并返回响应。
5. **RegisterService 函数**：
   - 注册服务，发送 `RegisterServiceRequest` 请求并返回响应。
6. **startSidecar 函数**：
   - 启动 sidecar，使用 `go run` 命令运行 sidecar 程序。
7. **startGRPCServer 函数**：
   - 启动 gRPC 服务器，监听指定端口并注册 `MiniGameRouterServer` 服务。
8. **connectToSidecar 函数**：
   - 连接到 sidecar，返回 gRPC 客户端连接和 `MiniGameRouterClient` 实例。
9. **main 函数**：
   - 解析命令行参数。
   - 启动 sidecar 并等待其启动完成。
   - 启动 gRPC 服务器。
   - 连接到 sidecar。
   - 注册自己的服务。
   - 从标准输入读取服务名称并请求服务。

### 5. sidecar.go

1. **命令行参数**：
   - 使用 `flag` 包定义命令行参数 `port`，用于指定 sidecar 的端口。
2. **MiniGameRouterServer 结构体**：
   - 实现了 `pb.UnimplementedMiniGameRouterServer` 接口。
3. **全局变量**：
   - `myIsDiscovered`：用于存储服务发现状态。
   - `myServicesStorage`：用于存储服务信息。
4. **WatchServiceName 函数**：
   - 监听服务名称的变化，并更新 `myServicesStorage`。
5. **RegisterService 方法**：
   - 注册服务，将服务信息存储到 etcd 中，并设置租约保活机制。
6. **DiscoverService 方法**：
   - 服务发现，获取服务信息并连接到另一个 sidecar。
7. **SayHello 方法**：
   - 简单的问候服务，连接到本地的另一个 sidecar 并发送问候消息。
8. **main 函数**：
   - 启动 gRPC 服务器，监听指定端口
