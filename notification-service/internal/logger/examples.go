package logger

// Пример использования логгера

/*
// В вашем main.go:

package main

import (
	"order-service/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Инициализация логгера с режимом DEBUG
	log := logger.Init(logger.DEBUG, "logs_test.md")
	defer log.Close()

	log.Info("Service started")

	// Обработка критических ошибок
	config, err := loadConfig()
	log.Must(err, "Failed to load configuration")

	log.Debug("Configuration loaded: %v", config)

	// Проверка nil значений
	db := initDatabase()
	log.MustNotNil(db, "Database connection is required")

	log.Info("Database connected successfully")

	// Запуск приложения
	log.Info("Starting HTTP server on :8080")
	server := startServer()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Info("Received shutdown signal")
	err = server.Shutdown()
	log.Must(err, "Failed to shutdown server gracefully")

	log.Info("Server shutdown completed successfully")
}

// Примеры использования в разных слоях:

// В repository層
func (r *userRepository) GetUser(ctx context.Context, id string) (*User, error) {
	log := logger.GetLogger()

	query := "SELECT * FROM users WHERE id = $1"
	var user User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		log.Error("Failed to fetch user %s: %v", id, err)
		return nil, err
	}

	log.Debug("User fetched successfully: %s", id)
	return &user, nil
}

// В service層
func (s *orderService) CreateOrder(ctx context.Context, order *Order) error {
	log := logger.GetLogger()

	log.Debug("Creating new order: %v", order.ID)

	// Проверяем что заказ не nil
	log.MustNotNil(order, "Order cannot be nil in CreateOrder")

	// Проверяем валидность
	if order.UserID == "" {
		log.Error("Order creation failed: missing user ID")
		return fmt.Errorf("user_id is required")
	}

	// Сохраняем в БД
	err := s.orderRepo.Create(ctx, order)
	if err != nil {
		log.Error("Failed to save order to database: %v", err)
		return err
	}

	log.Info("Order created successfully: %s", order.ID)
	return nil
}

// В delivery層 (HTTP handlers)
func (h *orderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	log := logger.GetLogger()

	var req CreateOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Error("Invalid request body: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Debug("Creating order from request: %v", req)

	order := &Order{
		ID:     uuid.New().String(),
		UserID: req.UserID,
		Amount: req.Amount,
	}

	err = h.orderService.CreateOrder(r.Context(), order)
	log.Must(err, "Failed to create order through service")  // Критическая ошибка - паника!

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
	log.Info("Order creation response sent: %s", order.ID)
}

// Примеры переключения режимов:
func enableDebugMode() {
	log := logger.GetLogger()
	log.SetMode(logger.DEBUG)
	log.Debug("Debug mode enabled for diagnostics")
}

func disableDebugMode() {
	log := logger.GetLogger()
	log.SetMode(logger.PROD)
	log.Info("Switched to PROD mode")
}
*/

// ExampleUsage демонстрирует основные операции логгера
//func ExampleUsage() {
//	log := Init(DEBUG, "example_logs.md")
//	defer log.Close()
//
//	fmt.Println("=== Logger Examples ===")
//
//	// 1. Различные уровни логирования
//	log.Info("Information message with param: %s", "value")
//	log.Debug("Debug message")
//	log.Warn("Warning message")
//	log.Error("Error message")
//
//	// 2. Функция Must для безопасной обработки ошибок
//	// Это не вызывает панику, так как err == nil
//	log.Must(nil, "No error occurred")
//
//	// 3. Проверка nil значений
//	value := "not nil"
//	log.MustNotNil(value, "Value check")
//
//	// 4. Изменение режима во время работы
//	log.SetMode(ERROR)
//	log.Info("This INFO won't be logged - we're in ERROR mode")
//	log.Error("But this ERROR will be logged")
//
//	fmt.Println("✓ Check example_logs.md to see the output")
//}
