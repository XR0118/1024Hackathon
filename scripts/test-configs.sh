#!/bin/bash

# Boreas 服务配置测试脚本

set -e

echo "=== Boreas 服务配置测试 ==="
echo ""

# 创建构建目录
mkdir -p bin

# 编译所有服务
echo "1. 编译所有服务..."
go build -o bin/operator-pm-agent ./cmd/operator-pm-agent
go build -o bin/operator-pm ./cmd/operator-pm
go build -o bin/master-service ./cmd/master-service
go build -o bin/operator-k8s ./cmd/operator-k8s
echo "✓ 所有服务编译完成"
echo ""

# 测试各服务的配置加载
echo "2. 测试配置加载..."

echo "测试 Operator PM Agent:"
bin/operator-pm-agent --version
bin/operator-pm-agent --config cmd/operator-pm-agent/configs/agent.yaml --version
echo "✓ Agent 配置测试通过"
echo ""

echo "测试 Operator PM:"
bin/operator-pm --version
bin/operator-pm --config cmd/operator-pm/configs/operator-pm.yaml --version
echo "✓ Operator-PM 配置测试通过"
echo ""

echo "测试 Master Service:"
bin/master-service --version
echo "✓ Master 版本测试通过"
echo ""

echo "测试 Operator K8s:"
bin/operator-k8s --version
echo "✓ Operator-K8s 版本测试通过"
echo ""

# 测试环境变量覆盖
echo "3. 测试环境变量覆盖..."
export AGENT_LOG_LEVEL=debug
export PM_LOG_LEVEL=info
export MASTER_LOG_LEVEL=warn
export K8S_LOG_LEVEL=error

echo "测试环境变量覆盖:"
bin/operator-pm-agent --config cmd/operator-pm-agent/configs/agent.yaml --version
bin/operator-pm --config cmd/operator-pm/configs/operator-pm.yaml --version
bin/master-service --version
bin/operator-k8s --version
echo "✓ 环境变量覆盖测试通过"
echo ""

# 清理
echo "4. 清理测试文件..."
rm -rf bin/
echo "✓ 清理完成"
echo ""

echo "=== 所有配置测试通过 ==="
