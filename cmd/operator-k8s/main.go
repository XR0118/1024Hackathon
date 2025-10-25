package main

import (
	"fmt"
	"log"

	"github.com/XR0118/1024Hackathon/internal/pkg/config"
	"github.com/XR0118/1024Hackathon/internal/pkg/database"
	"github.com/XR0118/1024Hackathon/internal/pkg/logger"
	"github.com/XR0118/1024Hackathon/internal/services/operator-k8s/handler"
	"github.com/XR0118/1024Hackathon/internal/services/operator-k8s/service"
	"github.com/gin-gonic/gin"
)

func main() {
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
