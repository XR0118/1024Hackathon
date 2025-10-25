package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/boreas/internal/pkg/models"
)

// AppStatus 应用状态
type AppStatus struct {
	Running bool
	PID     int
	Message string
}

// Runner 定义应用运行器接口
type Runner interface {
	// Restart 重启应用（如果未启动则直接启动）
	// currentDir: current 目录的绝对路径
	// version: 版本号
	// pkg: 部署包配置
	Restart(ctx context.Context, currentDir, version string, pkg *models.DeploymentPackage) error

	// Status 获取应用状态
	// appDir: 应用目录
	Status(ctx context.Context, appDir string) (*AppStatus, error)
}

// SimpleRunner 简单的启动脚本运行器（独立运行，不依赖 agent）
type SimpleRunner struct{}

// NewSimpleRunner 创建简单的启动脚本运行器
func NewSimpleRunner() *SimpleRunner {
	return &SimpleRunner{}
}

func (r *SimpleRunner) Restart(ctx context.Context, currentDir, version string, pkg *models.DeploymentPackage) error {
	appDir := filepath.Dir(currentDir)
	appName := filepath.Base(appDir)
	pidFile := filepath.Join(appDir, appName+".pid")

	// 1. 停止旧进程（如果存在）
	if pidData, err := os.ReadFile(pidFile); err == nil {
		pid := strings.TrimSpace(string(pidData))
		if pid != "" {
			// 发送 SIGTERM
			exec.Command("kill", "-TERM", pid).Run()
			// 等待 2 秒
			exec.Command("sleep", "2").Run()
			// 如果还在运行，发送 SIGKILL
			if exec.Command("kill", "-0", pid).Run() == nil {
				exec.Command("kill", "-9", pid).Run()
			}
			// 删除 PID 文件
			os.Remove(pidFile)
		}
	}

	// 2. 生成启动脚本
	var scriptContent string
	if pkg.Type == "binary" {
		// Binary 类型
		binaryPath := "app"
		if len(pkg.Command) > 0 {
			binaryPath = pkg.Command[0]
		}

		args := strings.Join(pkg.Args, " ")

		scriptContent = fmt.Sprintf(`#!/bin/bash
cd %s
nohup ./%s %s > /dev/null 2>&1 &
echo $! > %s
`, currentDir, binaryPath, args, pidFile)
	} else if pkg.Type == "script" {
		// Script 类型
		script := strings.Join(pkg.Command, "\n")
		scriptContent = fmt.Sprintf(`#!/bin/bash
cd %s
cat > temp_script.sh << 'SCRIPTEOF'
%s
SCRIPTEOF
nohup bash temp_script.sh > /dev/null 2>&1 &
echo $! > %s
rm -f temp_script.sh
`, currentDir, script, pidFile)
	} else {
		return fmt.Errorf("unsupported deployment type: %s", pkg.Type)
	}

	// 3. 写入环境变量文件（如果有）
	var envVars []string
	for k, v := range pkg.Environment {
		envVars = append(envVars, fmt.Sprintf("%s=%s", k, v))
	}

	if len(envVars) > 0 {
		envFile := filepath.Join(currentDir, ".env")
		envContent := strings.Join(envVars, "\n") + "\n"
		if err := os.WriteFile(envFile, []byte(envContent), 0644); err != nil {
			return fmt.Errorf("failed to create env file: %w", err)
		}
	}

	// 4. 写入并执行启动脚本
	startScript := filepath.Join(currentDir, "start.sh")
	if err := os.WriteFile(startScript, []byte(scriptContent), 0755); err != nil {
		return fmt.Errorf("failed to create start script: %w", err)
	}

	cmd := exec.CommandContext(ctx, "bash", startScript)
	cmd.Dir = currentDir

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start app: %w", err)
	}

	// 不等待命令完成，让它独立运行
	cmd.Process.Release()

	return nil
}

func (r *SimpleRunner) Status(ctx context.Context, appDir string) (*AppStatus, error) {
	appName := filepath.Base(appDir)
	pidFile := filepath.Join(appDir, appName+".pid")

	// 检查 PID 文件是否存在
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &AppStatus{
				Running: false,
				PID:     0,
				Message: "Not running",
			}, nil
		}
		return nil, fmt.Errorf("failed to read pid file: %w", err)
	}

	pid := strings.TrimSpace(string(pidData))
	if pid == "" {
		return &AppStatus{
			Running: false,
			PID:     0,
			Message: "Not running (empty PID file)",
		}, nil
	}

	// 解析 PID
	var pidInt int
	fmt.Sscanf(pid, "%d", &pidInt)

	// 通过 ps 命令检查进程是否存活
	cmd := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("ps -p %s > /dev/null 2>&1", pid))
	if err := cmd.Run(); err != nil {
		// 进程不存在，清理 PID 文件
		os.Remove(pidFile)
		return &AppStatus{
			Running: false,
			PID:     pidInt,
			Message: "Process not found",
		}, nil
	}

	return &AppStatus{
		Running: true,
		PID:     pidInt,
		Message: "Running",
	}, nil
}

// SupervisorRunner 使用 supervisor 的托管模式
type SupervisorRunner struct{}

// NewSupervisorRunner 创建 supervisor 运行器
func NewSupervisorRunner() *SupervisorRunner {
	return &SupervisorRunner{}
}

func (r *SupervisorRunner) Restart(ctx context.Context, currentDir, version string, pkg *models.DeploymentPackage) error {
	// TODO: 实现 supervisor 配置生成和重启
	return fmt.Errorf("supervisor runner not implemented yet")
}

func (r *SupervisorRunner) Status(ctx context.Context, appDir string) (*AppStatus, error) {
	// TODO: 实现 supervisor 状态查询
	return nil, fmt.Errorf("supervisor runner not implemented yet")
}

// SystemdRunner 使用 systemd 的托管模式
type SystemdRunner struct{}

// NewSystemdRunner 创建 systemd 运行器
func NewSystemdRunner() *SystemdRunner {
	return &SystemdRunner{}
}

func (r *SystemdRunner) Restart(ctx context.Context, currentDir, version string, pkg *models.DeploymentPackage) error {
	// TODO: 实现 systemd service 文件生成和重启
	return fmt.Errorf("systemd runner not implemented yet")
}

func (r *SystemdRunner) Status(ctx context.Context, appDir string) (*AppStatus, error) {
	// TODO: 实现 systemd 状态查询
	return nil, fmt.Errorf("systemd runner not implemented yet")
}
