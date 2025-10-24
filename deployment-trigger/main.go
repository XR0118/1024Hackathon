package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/XR0118/1024Hackathon/deployment-trigger/internal/handler"
	"github.com/XR0118/1024Hackathon/deployment-trigger/internal/service"
)

func main() {
	config := &service.Config{
		WebhookSecret:  os.Getenv("GITHUB_WEBHOOK_SECRET"),
		WorkDir:        os.Getenv("WORK_DIR"),
		DockerRegistry: os.Getenv("DOCKER_REGISTRY"),
		ManagementAPI:  os.Getenv("MANAGEMENT_API"),
	}

	if config.WorkDir == "" {
		config.WorkDir = "/tmp/deployment-trigger"
	}
	if config.ManagementAPI == "" {
		config.ManagementAPI = "http://localhost:8080/api/v1"
	}

	versionService := service.NewVersionService(config)
	webhookHandler := handler.NewWebhookHandler(versionService, config.WebhookSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook/github", webhookHandler.HandleGitHubWebhook)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Starting deployment-trigger service on :8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
