# Logger Package

Простой и мощный логгер для Go микросервиса с поддержкой трёх режимов: DEBUG, ERROR и PROD.

## Особенности

- **Три режима логирования:**
  - `DEBUG`: логирует всё (Info, Debug, Warn, Error)
  - `ERROR`: логирует только ошибки и критические события
  - `PROD`: логирует только ошибки и критические события (как ERROR)

- **Функция Must**: проверяет ошибку и паникует если она не nil
- **Функция MustNotNil**: проверяет что значение не nil, иначе паникует
- **Сохранение логов**: все логи сохраняются в markdown файл
- **Thread-safe**: использует мьютексы для безопасной работы в многопоточной среде
- **Singleton паттерн**: одна очередь инстанция логгера на приложение

## Использование

### Инициализация

```go
package main

import (
	"order-service/internal/logger"
)

func main() {
	// Инициализация логгера в DEBUG режиме
	log := logger.Init(logger.DEBUG, "logs_test.md")
	defer log.Close()

	log.Info("Application started")
	log.Debug("Debug information")
	log.Error("An error occurred")
}
```

### Логирование сообщений

```go
log := logger.GetLogger()

// Информационное сообщение (только DEBUG режим)
log.Info("User %s logged in", username)

// Отладочное сообщение (только DEBUG режим)
log.Debug("Database connection established to %s", dbHost)

// Ошибка (все режимы)
log.Error("Failed to fetch user: %v", err)

// Предупреждение (только DEBUG режим)
log.Warn("Low memory: %dMB", memLeft)

// Критическая ошибка (паника)
log.Fatal("Critical system failure: %v", err)
```

### Функция Must для обработки ошибок

```go
// Must проверяет ошибку и паникует если она не nil
log := logger.GetLogger()

file, err := os.Open("data.json")
log.Must(err, "Failed to open file")  // Если err != nil, паника!

// Проверка что значение не nil
user := fetchUser()
log.MustNotNil(user, "User cannot be nil")  // Если user == nil, паника!
```

### Изменение режима во время работы

```go
log := logger.GetLogger()

// Начинаем в PROD режиме
log.SetMode(logger.PROD)

// Переключаемся в DEBUG для диагностики
log.SetMode(logger.DEBUG)
log.Debug("Detailed debug info")

// Возвращаемся в PROD
log.SetMode(logger.PROD)
```

## Режимы логирования

### DEBUG
Логирует всё: Info, Debug, Warn, Error, Fatal

```markdown
# Logs - DEBUG Mode

**Started at:** 2026-08-26T21:49:09+05:00

- **[2026-08-26 21:49:09] INFO:** User john logged in
- **[2026-08-26 21:49:10] DEBUG:** Database query took 45ms
- **[2026-08-26 21:49:11] WARN:** Cache miss for key user:123
- **[2026-08-26 21:49:12] ERROR:** Failed to send email
```

### ERROR & PROD
Логируют только ошибки и критические события

```markdown
# Logs - ERROR Mode

**Started at:** 2026-08-26T21:49:09+05:00

- **[2026-08-26 21:49:12] ERROR:** Failed to send email
- **[2026-08-26 21:49:13] FATAL:** Database connection lost
```

## Структура логов

Каждая запись в логе имеет формат:
```
- **[TIMESTAMP] LEVEL:** Message
```

Где:
- `TIMESTAMP`: время события в формате `2006-01-02 15:04:05`
- `LEVEL`: уровень логирования (INFO, DEBUG, WARN, ERROR, FATAL)
- `Message`: сообщение лога с форматированием как в `fmt.Printf`

## Обработка ошибок с Must

Функция `Must` - это "fail fast" механизм. Если ошибка не nil:
1. Логирует ошибку в файл
2. Паникует с сообщением об ошибке

```go
// Пример 1: Обработка ошибки файловой системы
data, err := os.ReadFile("config.json")
log.Must(err, "Failed to read config")  // Паника если файл не существует

// Пример 2: Обработка ошибки базы данных
user, err := db.GetUser(ctx, id)
log.Must(err, "Failed to fetch user from database")  // Паника если ошибка

// Пример 3: Проверка nil указателя
user := fetchUser()
log.MustNotNil(user, "User pointer cannot be nil")  // Паника если user == nil
```

## Тестирование

Запуск тестов:
```bash
go test ./internal/logger -v
```

## API

### Функции пакета

- `Init(mode LogMode, filePath string) *Logger` - инициализирует логгер
- `GetLogger() *Logger` - возвращает глобальный экземпляр

### Методы Logger

- `SetMode(mode LogMode)` - устанавливает режим логирования
- `GetMode() LogMode` - возвращает текущий режим
- `Info(message string, args ...interface{})` - информационное сообщение
- `Debug(message string, args ...interface{})` - отладочное сообщение
- `Warn(message string, args ...interface{})` - предупреждение
- `Error(message string, args ...interface{})` - ошибка
- `Fatal(message string, args ...interface{})` - критическая ошибка (паника)
- `Must(err error, context string)` - проверка ошибки с паникой
- `MustNotNil(value interface{}, context string)` - проверка на nil с паникой
- `Close()` - закрытие логгера (запись времени завершения)

## Пример интеграции в main.go

```go
package main

import (
	"order-service/internal/logger"
	"os"
)

func main() {
	// Определяем режим по переменной окружения
	mode := logger.PROD
	if os.Getenv("DEBUG") == "true" {
		mode = logger.DEBUG
	}

	// Инициализируем логгер
	log := logger.Init(mode, "logs_test.md")
	defer log.Close()

	log.Info("Application starting in %s mode", mode)

	// Остальной код приложения...
}
```

## Безопасность

- **Thread-safe**: все операции защищены мьютексом
- **Гарантированная паника**: функция Must гарантирует остановку приложения при критической ошибке
- **Нет потери логов**: логи синхронно записываются на диск при каждом вызове
