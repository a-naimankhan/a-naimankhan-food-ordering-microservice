package main

import (
	"notification-service/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//TODO INIT CONFIGS

	//TODO INIT LOGGER
	logger := logger.Init(logger.DEBUG)
	defer logger.Close()
	logger.Info("Logger initialized successfully")
	logger.Info("Starting notification service...")

	//TODO INIT SERVICE LAYER

	//TODO INIT RABBITMQ CONSUMER

	//TODO GRACEFULL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Notification service is running. Press Ctrl+C to exit.")
	<-quit

	//TODO CLEANUP RESOURCES
	logger.Info("Shutdown signal received. Cleaning up resources...")
}
