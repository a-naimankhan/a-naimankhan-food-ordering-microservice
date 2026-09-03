package main

import (
	"context"
	"log"
	"net/http"
	"notification-service/internal/delivery"
	"notification-service/internal/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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

	//TODO INIT HTTP SERVER
	handler := delivery.NewNotificationHandler()
	srv := startServer(handler, ":8081", logger)

	//TODO GRACEFULL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Notification service is running. Press Ctrl+C to exit.")
	<-quit

	//TODO CLEANUP RESOURCES
	//cleanup http

	logger.Info("Shutdown signal received. Cleaning up resources...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	logger.Info("Shutting down HTTP server...")

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown:", "error", err.Error())
	} else {
		logger.Info("HTTP server stopped gracefully")
	}
}

func startServer(handler *delivery.NotificationHandler, port string, logger *logger.Logger) *http.Server {
	r := gin.Default()
	logger.Debug("STARTED HTTP SERVER")
	api := r.Group("/api/v1/notification")
	{
		api.GET("ping", handler.Ping)
	}

	logger.Info("HTTP routes registered")
	logger.Info("Listening on :", "port", 8081)

	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("error while starting http server: %s", err)
		}
	}()
	logger.Info("Server started")
	return srv

}
