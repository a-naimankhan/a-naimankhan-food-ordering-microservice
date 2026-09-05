package main

import (
	"context"
	"log"
	"net/http"
	"notification-service/internal/delivery"
	"notification-service/internal/logger"
	"notification-service/internal/infrastructure/rabbitmq"
	"notification-service/internal/domain"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
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

	// INIT RABBITMQ PUBLISHER
	// connect to RabbitMQ and initialize a publisher, close it on shutdown via defer
	logger.Debug("Connecting to RabbitMQ broker", "host", "localhost", "port", 5672)
	amqpConn, aerr := amqp.Dial("******localhost:5672/")
	logger.Must(aerr, "Failed to connect to RabbitMQ")

	publisher, perr := rabbitmq.NewPublisher(amqpConn, domain.ExchangeOrders)
	logger.Must(perr, "Failed to initialize RabbitMQ publisher")
	// ensure publisher closed on shutdown
	defer func() {
		if err := publisher.Close(); err != nil {
			logger.Error("Failed to close publisher during defer", "error", err.Error())
		}
	}()

	//INIT HTTP SERVER
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
