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
	"github.com/boreas/internal/config"
	"github.com/boreas/internal/database"
	"github.com/boreas/internal/handler"
	"github.com/boreas/internal/logger"
	"github.com/boreas/internal/middleware"
	"github.com/boreas/internal/repository/postgres"
	"github.com/boreas/internal/service"
	"github.com/boreas/internal/utils"
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

	// 初始化验证器
	utils.InitValidator()

	// 创建仓库
	versionRepo := postgres.NewVersionRepository(database.GetDB())
	appRepo := postgres.NewApplicationRepository(database.GetDB())
	envRepo := postgres.NewEnvironmentRepository(database.GetDB())
	deploymentRepo := postgres.NewDeploymentRepository(database.GetDB())
	taskRepo := postgres.NewTaskRepository(database.GetDB())
	workflowRepo := postgres.NewWorkflowRepository(database.GetDB())

	// 创建服务
	versionService := service.NewVersionService(versionRepo)
	appService := service.NewApplicationService(appRepo)
	envService := service.NewEnvironmentService(envRepo)
	// TODO: 实现 WorkflowManager
	// workflowMgr := workflow.NewWorkflowManager(workflowRepo, taskRepo)
	// deploymentService := service.NewDeploymentService(deploymentRepo, versionRepo, appRepo, envRepo, workflowMgr)
	taskService := service.NewTaskService(taskRepo)

	// 暂时注释掉未使用的变量
	_ = deploymentRepo
	_ = workflowRepo

	// 创建处理器
	versionHandler := handler.NewVersionHandler(versionService)
	appHandler := handler.NewApplicationHandler(appService)
	envHandler := handler.NewEnvironmentHandler(envService)
	// deploymentHandler := handler.NewDeploymentHandler(deploymentService)
	taskHandler := handler.NewTaskHandler(taskService)

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
			"service":    "management-service",
			"version":    Version,
			"build_time": BuildTime,
			"go_version": GoVersion,
		})
	})

	router.GET("/ready", func(c *gin.Context) {
		// 检查数据库连接
		sqlDB, err := database.GetDB().DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  err.Error(),
			})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// API 路由组
	api := router.Group("/api/v1")
	api.Use(middleware.Auth())

	// 版本管理路由
	versions := api.Group("/versions")
	{
		versions.POST("", versionHandler.CreateVersion)
		versions.GET("", versionHandler.GetVersionList)
		versions.GET("/:id", versionHandler.GetVersion)
		versions.DELETE("/:id", versionHandler.DeleteVersion)
	}

	// 应用管理路由
	applications := api.Group("/applications")
	{
		applications.POST("", appHandler.CreateApplication)
		applications.GET("", appHandler.GetApplicationList)
		applications.GET("/:id", appHandler.GetApplication)
		applications.PUT("/:id", appHandler.UpdateApplication)
		applications.DELETE("/:id", appHandler.DeleteApplication)
	}

	// 环境管理路由
	environments := api.Group("/environments")
	{
		environments.POST("", envHandler.CreateEnvironment)
		environments.GET("", envHandler.GetEnvironmentList)
		environments.GET("/:id", envHandler.GetEnvironment)
		environments.PUT("/:id", envHandler.UpdateEnvironment)
		environments.DELETE("/:id", envHandler.DeleteEnvironment)
	}

	// 部署管理路由
	// deployments := api.Group("/deployments")
	// {
	// 	deployments.POST("", deploymentHandler.CreateDeployment)
	// 	deployments.GET("", deploymentHandler.GetDeploymentList)
	// 	deployments.GET("/:id", deploymentHandler.GetDeployment)
	// 	deployments.POST("/:id/cancel", deploymentHandler.CancelDeployment)
	// 	deployments.POST("/:id/rollback", deploymentHandler.RollbackDeployment)
	// }

	// 任务管理路由
	tasks := api.Group("/tasks")
	{
		tasks.GET("", taskHandler.GetTaskList)
		tasks.GET("/:id", taskHandler.GetTask)
		tasks.POST("/:id/retry", taskHandler.RetryTask)
	}

	// 启动服务器
	server := &http.Server{
		Addr:    cfg.GetServerAddr(),
		Handler: router,
	}

	// 启动服务器
	go func() {
		logger.GetLogger().Info("Starting management service",
			zap.String("version", Version),
			zap.String("build_time", BuildTime),
			zap.String("go_version", GoVersion),
			zap.String("addr", cfg.GetServerAddr()),
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
