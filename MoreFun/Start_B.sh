#!/bin/bash

# 可执行文件的路径
executable="go run ./server/server.go"

# 初始参数
initial_port=35000

# 执行次数
num_executions=100

# 服务器的其他参数
name="B"
ip=$(hostname -I | awk '{print $1}')

for ((i=1; i<=num_executions; i++))
do
    port=$((initial_port + (i - 1) * 2))
    weight=$(shuf -i 1-30 -n 1) # 生成随机的权重
    echo "Starting execution $i with port $port and weight $weight"
    $executable --name $name --ip $ip --port $port --weight $weight &
done

echo "All instances started."