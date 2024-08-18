#!/bin/bash

# 可执行文件的路径
executable="go run ./server/server.go"

# 初始参数
initial_port=40000

# 执行次数
num_executions=100

# 服务器的其他参数
name="C"
ip=$(hostname -I | awk '{print $1}')
weight=5

for ((i=1; i<=num_executions; i++))
do
    port=$((initial_port + (i - 1) * 2))
    echo "Starting execution $i with port $port"
    $executable --name $name --ip $ip --port $port --weight $weight &
done

echo "All instances started."