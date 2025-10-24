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
			"service":    "webhook-service",
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

	// Webhook 路由组
	webhooks := router.Group("/api/v1/webhooks")
	{
		webhooks.POST("/github", func(c *gin.Context) {
			// TODO: 实现 GitHub Webhook 处理
			event := c.GetHeader("X-GitHub-Event")
			delivery := c.GetHeader("X-GitHub-Delivery")

			logger.GetLogger().Info("Received GitHub webhook",
				zap.String("event", event),
				zap.String("delivery", delivery),
			)

			c.JSON(http.StatusOK, gin.H{
				"message":  "Webhook received and processed",
				"event":    event,
				"delivery": delivery,
			})
		})
	}

	// 启动服务器
	server := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	// 启动服务器
	go func() {
		logger.GetLogger().Info("Starting webhook service",
			zap.String("version", Version),
			zap.String("build_time", BuildTime),
			zap.String("go_version", GoVersion),
			zap.String("addr", ":8082"),
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
