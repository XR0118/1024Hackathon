package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boreas/internal/pkg/client/operator"
	"github.com/boreas/internal/pkg/config"
	"github.com/boreas/internal/pkg/database"
	"github.com/boreas/internal/pkg/logger"
	"github.com/boreas/internal/pkg/middleware"
	"github.com/boreas/internal/pkg/utils"
	masterconfig "github.com/boreas/internal/services/master/config"
	"github.com/boreas/internal/services/master/handler"
	"github.com/boreas/internal/services/master/repository/postgres"
	"github.com/boreas/internal/services/master/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
	GoVersion = "unknown"
)

func main() {
	// 命令行参数
	var (
		showVersion = flag.Bool("version", false, "显示版本信息")
	)
	flag.Parse()

	// 显示版本信息
	if *showVersion {
		fmt.Printf("Boreas Master Service\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		fmt.Printf("Go Version: %s\n", GoVersion)
		os.Exit(0)
	}

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 加载 master 特定配置（包含 operator 配置）
	masterCfg, err := masterconfig.Load("")
	if err != nil {
		log.Fatalf("Failed to load master config: %v", err)
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

	// 初始化 Operator Manager
	logger.GetLogger().Info("Initializing operator manager")
	envList, _, err := envRepo.List(context.Background(), nil)
	if err != nil {
		logger.GetLogger().Fatal("Failed to list environments", zap.Error(err))
	}

	operatorConfig := &operator.Config{
		K8SOperatorURL: masterCfg.Operator.K8SOperatorURL,
		PMOperatorURL:  masterCfg.Operator.PMOperatorURL,
		UseMock:        masterCfg.Operator.UseMock,
	}

	operatorManager, err := operator.InitializeOperators(envList, operatorConfig)
	if err != nil {
		logger.GetLogger().Fatal("Failed to initialize operator manager", zap.Error(err))
	}
	logger.GetLogger().Info("Operator manager initialized",
		zap.Int("registered_operators", len(operatorManager.ListOperators())),
		zap.Bool("use_mock", masterCfg.Operator.UseMock),
	)

	// 创建服务
	versionService := service.NewVersionService(versionRepo)
	appService := service.NewApplicationService(appRepo, versionRepo, deploymentRepo, operatorManager)
	envService := service.NewEnvironmentService(envRepo)
	triggerService := service.NewTriggerService(&service.TriggerConfig{
		WebhookSecret:  cfg.Trigger.WebhookSecret,
		WorkDir:        cfg.Trigger.WorkDir,
		DockerRegistry: cfg.Trigger.DockerRegistry,
		Apps:           appService,
		Version:        versionService,
	})
	webhookHandler := handler.NewWebhookHandler(triggerService, cfg.Trigger.WebhookSecret)

	// 创建工作流控制器
	workflowController := service.NewWorkflowController(
		taskRepo,
		deploymentRepo,
		versionRepo,
		service.WorkflowConfig{},
		logger.GetLogger(),
	)
	workflowController.Start()
	defer workflowController.Stop()

	deploymentService := service.NewDeploymentService(deploymentRepo, versionRepo, appRepo, envRepo, workflowController)
	taskService := service.NewTaskService(taskRepo)

	// 创建处理器
	versionHandler := handler.NewVersionHandler(versionService)
	appHandler := handler.NewApplicationHandler(appService)
	envHandler := handler.NewEnvironmentHandler(envService)
	deploymentHandler := handler.NewDeploymentHandler(deploymentService)
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

	router.POST("/webhook/github", gin.WrapF(webhookHandler.HandleGitHubWebhook))

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
	// 开发环境：暂时禁用认证中间件，方便前后端联调
	// api.Use(middleware.Auth())

	// 版本管理路由
	versions := api.Group("/versions")
	{
		versions.POST("", versionHandler.CreateVersion)
		versions.GET("", versionHandler.GetVersionList)
		versions.GET("/:version", versionHandler.GetVersion)
		versions.DELETE("/:version", versionHandler.DeleteVersion)
		versions.POST("/:version/rollback", versionHandler.RollbackVersion)
	}

	// 应用管理路由
	applications := api.Group("/applications")
	{
		applications.POST("", appHandler.CreateApplication)
		applications.GET("", appHandler.GetApplicationList)
		applications.GET("/:name", appHandler.GetApplication) // 按应用名称查询
		applications.PUT("/:id", appHandler.UpdateApplication)
		applications.DELETE("/:id", appHandler.DeleteApplication)
	}

	// 应用版本信息路由（使用应用名称）
	api.GET("/applications/:name/versions", appHandler.GetApplicationVersions)                          // 版本详情（按环境组织）
	api.GET("/applications/:name/versions/summary", appHandler.GetApplicationVersionsSummary)           // 版本概要
	api.GET("/applications/:name/versions/:version/coverage", appHandler.GetApplicationVersionCoverage) // 版本覆盖率（累积）

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
	deployments := api.Group("/deployments")
	{
		deployments.POST("", deploymentHandler.CreateDeployment)
		deployments.GET("", deploymentHandler.GetDeploymentList)
		deployments.GET("/:id", deploymentHandler.GetDeployment)
		deployments.POST("/:id/start", deploymentHandler.StartDeployment)
		deployments.POST("/:id/pause", deploymentHandler.PauseDeployment)
		deployments.POST("/:id/resume", deploymentHandler.ResumeDeployment)
	}

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
