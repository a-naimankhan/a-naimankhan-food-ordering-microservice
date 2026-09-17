package main

import (
	"context"
	"net/http"
	"notification-service/internal/domain"
	"notification-service/internal/infrastructure/rabbitmq"
	"notification-service/internal/logger"
	"notification-service/internal/probes"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	//TODO INIT CONFIGS
	rabbitURL := getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
	healthAddr := getEnv("HEALTH_ADDR", ":8081")
	//TODO INIT LOGGER
	log := logger.Init(logger.DEBUG)
	defer log.Close()
	log.Info("Logger initialized successfully")
	log.Info("Starting notification service...")

	//TODO INIT SERVICE LAYER

	// INIT RABBITMQ PUBLISHER
	// connect to RabbitMQ and initialize a publisher, close it on shutdown via defer
	log.Debug("Connecting to RabbitMQ broker", "host", "localhost", "port", 5672)
	amqpConn, aerr := amqp.Dial(rabbitURL)
	log.Must(aerr, "Failed to connect to RabbitMQ")

	publisher, perr := rabbitmq.NewPublisher(amqpConn, domain.ExchangeOrders)
	log.Must(perr, "Failed to initialize RabbitMQ publisher")

	//INIT HTTP SERVER
	serverErrChan := make(chan error, 1)
	srv := probes.NewHealthServer(healthAddr) //ну пока что так
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server closed while opening", "error", err.Error())
			//здесь если че добавим что приложение падало но пока одного логера хватает ведь я к8 не скоро добавлю
		}
	}()

	//TODO GRACEFULL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Notification service is running...")

	select {
	case sig := <-quit:
		log.Info("Shutdown signal received", "signal", sig.String())
	case err := <-serverErrChan:
		log.Fatal("Probe server failed to start", "error", err.Error())
	}
	//TODO CLEANUP RESOURCES
	//cleanup http

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Info("Shutting down HTTP server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown:", "error", err.Error())
	} else {
		log.Info("HTTP server stopped gracefully")
	}

	//RabbitMQ publisher closing
	log.Info("Closing RabbitMQ publisher...")
	if err := publisher.Close(); err != nil {
		log.Error("Failed to close publisher", "error", err)
	} else {
		log.Info("Publisher closed successfully")
	}

	//RabbiqMQ connection closing
	log.Info("Closing RabbitMQ connection...")
	if err := amqpConn.Close(); err != nil {
		log.Error("Failed to close connection", "error", err.Error())
	} else {
		log.Info("RabbitMQ connection closed succesfully")
	}

}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
