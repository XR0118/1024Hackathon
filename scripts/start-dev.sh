#!/bin/bash

# 启动开发环境脚本

set -e

echo "🚀 Starting Boreas Development Environment"

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# 启动数据库和 Redis
echo "📦 Starting PostgreSQL and Redis..."
docker-compose up -d postgres redis

# 等待数据库启动
echo "⏳ Waiting for database to be ready..."
sleep 10

# 检查数据库连接
echo "🔍 Checking database connection..."
until docker-compose exec postgres pg_isready -U boreas -d boreas; do
    echo "Waiting for database..."
    sleep 2
done

echo "✅ Database is ready!"

# 运行数据库迁移
echo "🗄️ Running database migrations..."
# 这里需要安装 migrate 工具
# make migrate-up

# 构建项目
echo "🔨 Building project..."
make build

# 启动服务
echo "🎯 Starting services..."
echo "Management Service: http://localhost:8080"
echo "Deploy Service: http://localhost:8081"
echo "Webhook Service: http://localhost:8082"
echo "Nginx: http://localhost:80"

# 在后台启动服务
make run-dev &

# 等待用户中断
echo "Press Ctrl+C to stop all services"
wait
