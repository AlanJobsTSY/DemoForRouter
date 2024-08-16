#!/bin/bash

# 可执行文件的路径
executable="go run ./map.go"

# 执行次数
num_executions=20

for ((i=1; i<=num_executions; i++))
do
    $executable  &
done

echo "All instances started."