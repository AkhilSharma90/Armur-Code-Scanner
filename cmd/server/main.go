package main

import (
	"armur-codescanner/internal/api"
	"armur-codescanner/internal/redis"
	"armur-codescanner/internal/worker"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"log"
	"os"
)

func main() {
	router := gin.Default()

	go func() {
		if err := startAsynqWorker(); err != nil {
			log.Fatalf("Failed to start Asynq worker: %v", err)
		}
	}()

	router.POST("/api/v1/scan/repo", api.ScanHandler)
	router.GET("/api/v1/status/:task_id", api.TaskStatus)
	router.POST("/api/v1/advanced-scan/repo", api.AdvancedScanResult)
	router.POST("/api/v1/scan/file", api.ScanFile)
	router.GET("/api/v1/reports/owasp/:task_id", api.TaskOwasp)
	router.GET("/api/v1/reports/sans/:task_id", api.TaskSans)
	port := os.Getenv("APP_PORT")
	fmt.Println(port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

func startAsynqWorker() error {
	server := asynq.NewServer(
		redis.RedisClientOptions(),
		asynq.Config{
			Concurrency: 10,
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle("scan:repo", &worker.ScanTaskHandler{})

	// Start the Asynq server and process tasks
	return server.Start(mux)
}
