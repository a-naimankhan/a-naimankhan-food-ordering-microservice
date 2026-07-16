package main

import (
	"fmt"
	"log"
	"order-service/internal/delivery"
	"order-service/internal/infrastructure/rabbitmq"
	"order-service/internal/repository"
	"order-service/internal/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sqlx.Open("postgres", "postgres://user:password123@localhost:5432/orders_db?sslmode=disable")
	if err != nil {
		log.Fatalf("could not connect to postgres: %v", err)
	}

	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()

	publisher := rabbitmq.NewRabbitPublisher(*ch, "orders_events")
	repo := repository.NewOrderRepo(db)

	//paymentClient := domain.PaymentClient() // Initialize your payment client here
	svc := service.NewOrderService(repo /*paymentClient*/, publisher)
	handler := delivery.NewOrderHandler(svc)

	startServer(handler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("server is shutting down")
}

func startServer(handler *delivery.OrderHandler) {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/ping", handler.Ping)
		api.POST("/orders", handler.CreateOrder)
		api.GET("/orders/:id", handler.GetOrder)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
