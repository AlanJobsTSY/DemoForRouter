#!/bin/bash

# 安装 Docker Compose
echo "Installing Docker Compose..."
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 检查 Docker 是否已安装
if ! command -v docker &> /dev/null
then
    echo "Docker is not installed. Please install Docker first."
    exit 1
fi

# 检查 Docker Compose 是否已安装
if ! command -v docker-compose &> /dev/null
then
    echo "Docker Compose installation failed."
    exit 1
fi

# 切换到 etcd-config 目录
cd etcd-config

# 检查 docker-compose.yml 文件是否存在
if [ ! -f "docker-compose.yml" ]; then
    echo "docker-compose.yml file not found in the etcd-config directory."
    exit 1
fi

# 启动 Docker Compose 服务
echo "Starting Docker Compose services..."
docker-compose up -d

echo "All services started."