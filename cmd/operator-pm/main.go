package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boreas/internal/pkg/logger"
	"github.com/boreas/internal/services/operator-pm/config"
	"github.com/boreas/internal/services/operator-pm/handler"
	"github.com/boreas/internal/services/operator-pm/service"
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
		showVersion = flag.Bool("version", false, "显示版本信息")
	)
	flag.Parse()

	// 显示版本信息
	if *showVersion {
		fmt.Printf("Boreas Operator PM\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
		os.Exit(0)
	}

	// 加载配置
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化日志
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	// 初始化服务
	operatorService := service.NewOperatorPMService(cfg)

	// 初始化处理器
	operatorHandler := handler.NewOperatorPMHandler(operatorService)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 注册路由
	operatorHandler.RegisterRoutes(r)

	// 启动服务
	log.Printf("Operator-PM service starting on %s", cfg.GetServerAddr())
	if err := r.Run(cfg.GetServerAddr()); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
