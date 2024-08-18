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

func weightedLeastConnections(myServicesStorage *ServicesStorage, svrName string) string {
	cli := etcd.NewEtcdCli()
	defer cli.Close()
	myServicesStorage.Lock()
	defer myServicesStorage.Unlock()
	var addr string
	if instances, ok := myServicesStorage.ServicesStorage[svrName]; ok && len(instances) > 0 {
		minnConn := math.MaxInt
		minnWeight := 0
		var key string
		for k, v := range instances {
			parts := strings.Split(v, ":")
			partConn, _ := strconv.Atoi(parts[6])
			partWeight, _ := strconv.Atoi(parts[4])
			if minnConn*partWeight >= partConn*minnWeight {
				minnConn = partConn
				minnWeight = partWeight
				addr = fmt.Sprintf("%s:%s", parts[1], parts[2])
				key = k
			}
		}

		// 替换 parts 切片中的第 6 个元素
		parts := strings.Split(instances[key], ":")
		parts[6] = strconv.Itoa(minnConn + 1)
		// 将 parts 切片重新拼接成字符串
		newValue := strings.Join(parts, ":")
		//log.Printf(newValue)
		_, err := cli.Put(context.Background(), key, newValue, clientv3.WithIgnoreLease())
		if err != nil {
			fmt.Println("Failed to update etcd:", err)
			return ""
		}
	}

	return addr
}
