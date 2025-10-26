package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boreas/internal/pkg/logger"
	"github.com/boreas/internal/services/operator-k8s/config"
	"github.com/boreas/internal/services/operator-k8s/handler"
	"github.com/boreas/internal/services/operator-k8s/service"
	"github.com/gin-gonic/gin"
)

var (
	version   = "1.0.0"
	buildTime = "unknown"
)

func main() {
	var (
		showVersion = flag.Bool("version", false, "显示版本信息")
		configPath  = flag.String("config", "", "配置文件路径")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("Boreas Operator K8s\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logger.Init(cfg.Log.Level, cfg.Log.Format)

	operatorService, err := service.NewOperatorK8sService(
		cfg.K8s.ConfigPath,
		cfg.K8s.Namespace,
		cfg.K8s.Timeout,
	)
	if err != nil {
		log.Fatal("Failed to initialize operator service:", err)
	}

	operatorHandler := handler.NewOperatorK8sHandler(operatorService)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	operatorHandler.RegisterRoutes(r)

	log.Printf("Operator-K8s service starting on %s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
