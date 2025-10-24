package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/guaguasong/1024Hackathon/internal/config"
	"github.com/guaguasong/1024Hackathon/internal/database"
	"github.com/guaguasong/1024Hackathon/internal/logger"
	"github.com/guaguasong/1024Hackathon/internal/middleware"
	"go.uber.org/zap"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log.Level, cfg.Log.Format); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 初始化数据库
	if err := database.Init(cfg); err != nil {
		logger.GetLogger().Fatal("Failed to initialize database", zap.Error(err))
	}
	defer database.Close()

	// 设置 Gin 模式
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	router := gin.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":     "healthy",
			"service":    "deploy-service",
			"version":    Version,
			"build_time": BuildTime,
			"go_version": GoVersion,
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		// 检查数据库连接
		if err := database.GetDB().DB().Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// 内部 API 路由组
	internal := router.Group("/internal/v1")
	{
		deploy := internal.Group("/deploy")
		{
			deploy.GET("/info/:deployment_id", func(c *gin.Context) {
				// TODO: 实现获取部署信息
				c.JSON(http.StatusOK, gin.H{
					"deployment_id":  c.Param("deployment_id"),
					"status":         "running",
					"replicas":       3,
					"ready_replicas": 2,
					"updated_at":     time.Now(),
				})
			})

			deploy.GET("/health/:deployment_id", func(c *gin.Context) {
				// TODO: 实现健康检查
				c.JSON(http.StatusOK, gin.H{
					"healthy": true,
					"message": "All replicas are healthy",
					"checks": []gin.H{
						{
							"name":    "pod_status",
							"status":  "healthy",
							"message": "3/3 pods running",
						},
					},
				})
			})

			deploy.GET("/logs/:deployment_id", func(c *gin.Context) {
				// TODO: 实现获取日志
				c.JSON(http.StatusOK, gin.H{
					"logs": []gin.H{
						{
							"timestamp": time.Now(),
							"message":   "Deployment service is running",
						},
					},
				})
			})
		}
	}

	// 启动服务器
	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	// 启动服务器
	go func() {
		logger.GetLogger().Info("Starting deploy service",
			zap.String("version", Version),
			zap.String("build_time", BuildTime),
			zap.String("go_version", GoVersion),
			zap.String("addr", ":8081"),
		)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.GetLogger().Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.GetLogger().Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.GetLogger().Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.GetLogger().Info("Server exited")
}
