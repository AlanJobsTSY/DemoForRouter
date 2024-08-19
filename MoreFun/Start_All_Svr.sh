#!/bin/bash

# 可执行文件的路径
executable="go run ./server/server.go"

# 初始参数
initial_port=30000

# 每个名称的执行次数
num_executions=20

# 服务器的其他参数
names=("A" "B" "C" "D" "E" "F")
ip=$(hostname -I | awk '{print $1}')

for name in "${names[@]}"
do
    for ((i=1; i<=num_executions; i++))
    do
        port=$((initial_port + (i - 1) * 2))
        weight=$(shuf -i 1-30 -n 1) # 生成随机的权重
        echo "Starting execution $i for name $name with port $port and weight $weight"
        $executable --name $name --ip $ip --port $port --weight $weight &
    done
    initial_port=$((initial_port + num_executions * 2)) # 更新初始端口以避免端口冲突
done

echo "All instances started."