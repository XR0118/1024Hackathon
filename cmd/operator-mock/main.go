package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boreas/internal/pkg/logger"
	"github.com/boreas/internal/services/operator-mock/config"
	"github.com/boreas/internal/services/operator-mock/handler"
	"github.com/boreas/internal/services/operator-mock/service"
	"github.com/gin-gonic/gin"
)

var (
	version   = "1.0.0"
	buildTime = "unknown"
)

func main() {
	var (
		configPath  = flag.String("config", "", "配置文件路径")
		showVersion = flag.Bool("version", false, "显示版本信息")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("Boreas Operator Mock\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
		os.Exit(0)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	logger.Init(cfg.Log.Level, cfg.Log.Format)

	mockService := service.MockAgent

	mockHandler := handler.NewOperatorMockHandler(mockService)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	mockHandler.RegisterRoutes(r)

	log.Printf("Operator-Mock service starting on %s", cfg.GetServerAddr())
	if err := r.Run(cfg.GetServerAddr()); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
