package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boreas/internal/pkg/config"
	"github.com/boreas/internal/pkg/database"
	"github.com/boreas/internal/pkg/logger"
	"github.com/boreas/internal/services/operator-k8s/handler"
	"github.com/boreas/internal/services/operator-k8s/service"
	"github.com/gin-gonic/gin"
)

var (
	version   = "1.0.0"
	buildTime = "unknown"
)

func main() {
	// 命令行参数
	var (
		showVersion = flag.Bool("version", false, "显示版本信息")
	)
	flag.Parse()

	// 显示版本信息
	if *showVersion {
		fmt.Printf("Boreas Operator K8s\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
		os.Exit(0)
	}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化日志
	logger.Init(cfg.Log.Level, cfg.Log.Format)

	// 初始化数据库
	err = database.Init(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化服务
	operatorService := service.NewOperatorK8sService(database.GetDB())

	// 初始化处理器
	operatorHandler := handler.NewOperatorK8sHandler(operatorService)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 注册路由
	operatorHandler.RegisterRoutes(r)

	// 启动服务
	log.Printf("Operator-K8s service starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
