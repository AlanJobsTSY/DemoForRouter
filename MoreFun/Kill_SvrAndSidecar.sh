#!/bin/bash

# 检查是否提供了端口号
if [ -z "$1" ]; then
    echo "Usage: $0 <port>"
    exit 1
fi

# 要关闭的端口
target_port=$1
next_port=$((target_port + 1))

# 定义一个函数来终止指定端口上的进程
terminate_process_on_port() {
    local port=$1
    local pid=$(lsof -t -i :$port)

    if [ -z "$pid" ]; then
        echo "No process found running on port $port."
    else
        echo "Killing process $pid running on port $port."
        kill -9 $pid
        echo "Process $pid has been terminated."
    fi
}

# 终止 target_port 和 next_port 上的进程
terminate_process_on_port $target_port
terminate_process_on_port $next_port