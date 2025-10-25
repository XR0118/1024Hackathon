package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/boreas/internal/pkg/logger"
	"github.com/boreas/internal/services/operator-pm-agent/config"
	"github.com/boreas/internal/services/operator-pm-agent/handler"
	"github.com/boreas/internal/services/operator-pm-agent/service"
	"github.com/gin-gonic/gin"
)

var (
	version   = "1.0.0"
	buildTime = "unknown"
)

func main() {
	// 命令行参数
	var (
		configPath  = flag.String("config", "", "配置文件路径")
		workDir     = flag.String("work-dir", "", "工作目录（覆盖配置文件）")
		port        = flag.Int("port", 0, "服务端口（覆盖配置文件）")
		host        = flag.String("host", "", "服务地址（覆盖配置文件）")
		agentID     = flag.String("agent-id", "", "Agent ID（覆盖配置文件）")
		showVersion = flag.Bool("version", false, "显示版本信息")
	)
	flag.Parse()

	// 显示版本信息
	if *showVersion {
		fmt.Printf("Boreas Operator PM Agent\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
		os.Exit(0)
	}

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 命令行参数覆盖配置
	if *host != "" {
		cfg.Server.Host = *host
	}
	if *port != 0 {
		cfg.Server.Port = *port
	}
	if *workDir != "" {
		cfg.Agent.WorkDir = *workDir
	}
	if *agentID != "" {
		cfg.Agent.ID = *agentID
	}

	// 初始化日志
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	// 确保工作目录存在
	agentWorkDir := cfg.GetAgentWorkDir()
	if err := os.MkdirAll(agentWorkDir, 0755); err != nil {
		log.Fatal("Failed to create work directory:", err)
	}

	// 初始化服务
	agentService, err := service.NewAgentServiceWithConfig(
		agentWorkDir,
		cfg.IsDockerEnabled(),
		cfg.GetDockerSocketPath(),
	)
	if err != nil {
		log.Fatal("Failed to initialize agent service:", err)
	}

	// 初始化处理器
	agentHandler := handler.NewAgentHandler(agentService)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 注册路由
	agentHandler.RegisterRoutes(r)

	// 启动服务
	serverAddr := cfg.GetServerAddr()
	log.Printf("Boreas Operator PM Agent starting on %s", serverAddr)
	log.Printf("Agent ID: %s", cfg.GetAgentID())
	log.Printf("Work directory: %s", agentWorkDir)
	log.Printf("Docker enabled: %v", cfg.IsDockerEnabled())
	log.Printf("Version: %s", version)

	// 优雅关闭
	go func() {
		if err := r.Run(serverAddr); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
