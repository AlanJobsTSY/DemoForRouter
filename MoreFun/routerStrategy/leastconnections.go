package routerStrategy

import (
	"MoreFun/etcd"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"math"
	"strconv"
	"strings"
)

func leastConnections(myServicesStorage *ServicesStorage, svrName string) string {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	minn := math.MaxInt
	var addr string
	var key string
	for k, v := range myServicesStorage.ServicesStorage[svrName] {
		parts := strings.Split(v, ":")
		partInt, _ := strconv.Atoi(parts[6])
		if minn >= partInt {
			minn = partInt
			addr = fmt.Sprintf("%s:%s", parts[1], parts[2])
			key = k
		}
	}

	// 替换 parts 切片中的第 6 个元素
	parts := strings.Split(myServicesStorage.ServicesStorage[svrName][key], ":")
	parts[6] = strconv.Itoa(minn + 1)
	// 将 parts 切片重新拼接成字符串
	newValue := strings.Join(parts, ":")

	_, err := cli.Put(context.Background(), key, newValue, clientv3.WithIgnoreLease())
	if err != nil {
		fmt.Println("Failed to update etcd:", err)
		return ""
	}

	return addr
}
