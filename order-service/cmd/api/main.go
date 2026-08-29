package main

import (
	"context"
	"net/http"
	"order-service/internal/delivery"
	"order-service/internal/infrastructure/rabbitmq"
	"order-service/internal/logger"
	"order-service/internal/repository"
	"order-service/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	////TODO init configs
	//config, err := Must()
	//if err != nil {
	//	panic(err)
	//}

	log := logger.Init(logger.DEBUG)
	defer log.Close()

	log.Info("Order Service starting")

	log.Debug("Connecting to PostgreSQL database", "host", "localhost", "port", 5432)
	db, err := sqlx.Open("postgres", "postgres://user:password123@localhost:5432/orders_db?sslmode=disable")
	log.Must(err, "Failed to connect to PostgreSQL")

	log.Info("PostgreSQL connected successfully")
	log.Debug("Connecting to RabbitMQ broker", "host", "localhost", "port", 5672)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	log.Must(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	log.Must(err, "Failed to open RabbitMQ channel")

	log.Info("RabbitMQ connected successfully")
	log.Debug("Initializing repository, service, and handler")

	publisher := rabbitmq.NewRabbitPublisher(ch, "orders_events")
	repo := repository.NewOrderRepo(db)
	svc := service.NewOrderService(repo, publisher)
	handler := delivery.NewOrderHandler(svc)

	log.Info("All components initialized successfully")
	log.Info("Starting HTTP server on :8080")

	srv := startServer(handler, log)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Shutting down HTTP server
	log.Info("Shutting down HTTP server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown:", err)
	} else {
		log.Info("HTTP server stopped gracefully")
	}

	//shutting down RABBITMQ
	log.Info("Closing RabbitMQ channel and connection...")
	if err := ch.Close(); err != nil {
		log.Error("Failed to close RabbitMQ channel: ", err)
	}
	if err := conn.Close(); err != nil {
		log.Error("Failed to close RabbitMQ connection: ", err)
	}
	log.Info("RabbitMQ connection closed")

	//shutting down Database
	log.Info("Shutting down database connection...")
	if err := db.Close(); err != nil {
		log.Error("Failed to close database: ", err)
	} else {
		log.Info("Database connection closed")
	}

	log.Info("Server shutdown complete")
}

func startServer(handler *delivery.OrderHandler, log *logger.Logger) *http.Server {
	r := gin.Default()

	log.Debug("Registering HTTP routes")
	api := r.Group("/api/v1")
	{
		api.GET("/ping", handler.Ping)
		api.POST("/orders", handler.CreateOrder)
		api.GET("/orders/:id", handler.GetOrder)
	}

	log.Info("HTTP routes registered")
	log.Info("Listening on :8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("error while starting http server: %s", err)
		}
	}()
	log.Info("Server started")
	return srv

}
