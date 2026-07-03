package main

import (
	"fmt"
	"log"
	"order-service/internal/delivery"
	"order-service/internal/repository"
	"order-service/internal/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "postgres://user:password123@localhost:5432/orders_db?sslmode=disable")
	if err != nil {
		log.Fatalf("could not connect to postgres: %v", err)
	}

	repo := repository.NewOrderRepo(db)
	svc := service.NewOrderService(repo)
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
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
