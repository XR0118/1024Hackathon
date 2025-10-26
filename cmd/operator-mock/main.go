package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boreas/internal/pkg/logger"
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
		showVersion = flag.Bool("version", false, "显示版本信息")
		port        = flag.String("port", "8082", "服务端口")
		logLevel    = flag.String("log-level", "info", "日志级别")
		logFormat   = flag.String("log-format", "json", "日志格式")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("Boreas Operator Mock\n")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build Time: %s\n", buildTime)
		os.Exit(0)
	}

	logger.Init(*logLevel, *logFormat)

	mockService := service.MockAgent

	mockHandler := handler.NewOperatorMockHandler(mockService)

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	mockHandler.RegisterRoutes(r)

	addr := ":" + *port
	log.Printf("Operator-Mock service starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
