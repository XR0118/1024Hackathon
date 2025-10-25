#!/bin/bash

# Boreas Operator PM Agent 配置测试脚本

set -e

echo "=== Boreas Operator PM Agent 配置测试 ==="

# 检查 Go 环境
if ! command -v go &> /dev/null; then
    echo "错误: 未找到 Go 环境"
    exit 1
fi

# 进入项目目录
cd "$(dirname "$0")/.."

echo "1. 编译 Agent..."
go build -o /tmp/boreas-agent ./cmd/operator-pm-agent

echo "2. 测试默认配置..."
echo "   运行: /tmp/boreas-agent --version"
/tmp/boreas-agent --version

echo "3. 测试配置文件加载..."
echo "   运行: /tmp/boreas-agent --config configs/agent.yaml --version"
/tmp/boreas-agent --config configs/agent.yaml --version

echo "4. 测试命令行参数覆盖..."
echo "   运行: /tmp/boreas-agent --port 9090 --agent-id test-agent --version"
/tmp/boreas-agent --port 9090 --agent-id test-agent --version

echo "5. 测试环境变量覆盖..."
echo "   设置环境变量: AGENT_ID=env-agent AGENT_WORK_DIR=/tmp/test-agent"
export AGENT_ID=env-agent
export AGENT_WORK_DIR=/tmp/test-agent
echo "   运行: /tmp/boreas-agent --version"
/tmp/boreas-agent --version

echo "6. 测试配置验证..."
echo "   测试无效端口..."
if /tmp/boreas-agent --port 99999 --version 2>/dev/null; then
    echo "   警告: 无效端口未被检测到"
else
    echo "   正常: 无效端口被正确检测"
fi

echo "7. 清理..."
rm -f /tmp/boreas-agent
unset AGENT_ID
unset AGENT_WORK_DIR

echo "=== 配置测试完成 ==="
