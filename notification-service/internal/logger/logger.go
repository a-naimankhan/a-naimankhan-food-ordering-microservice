package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type LogMode string

const (
	DEBUG LogMode = "DEBUG"
	ERROR LogMode = "ERROR"
	PROD  LogMode = "PROD"
)

type Logger struct {
	mode   LogMode
	logger *slog.Logger
	mu     sync.Mutex
}

var (
	instance    *Logger
	once        sync.Once
	globalLevel = new(slog.LevelVar)
)

// Init инициализирует логгер
func Init(mode LogMode) *Logger {

	once.Do(func() {
		// Определяем уровень логирования
		switch mode {
		case DEBUG:
			globalLevel.Set(slog.LevelDebug)
		case ERROR, PROD:
			globalLevel.Set(slog.LevelError)
		default:
			globalLevel.Set(slog.LevelInfo)
		}

		// Создаем обработчик с структурированным выводом в stdout
		opts := &slog.HandlerOptions{
			Level:     globalLevel,
			AddSource: true,
		}

		handler := slog.NewTextHandler(os.Stdout, opts)
		logger := slog.New(handler)

		instance = &Logger{
			mode:   mode,
			logger: logger,
		}

		instance.Info("Logger initialized", "mode", mode, "timestamp", time.Now().Format(time.RFC3339))
	})
	return instance
}

// GetLogger возвращает глобальный экземпляр логгера
func GetLogger() *Logger {
	if instance == nil {
		return Init(PROD)
	}
	return instance
}

// SetMode устанавливает режим логирования
func (l *Logger) SetMode(level LogMode) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.mode = level

	switch level {
	case DEBUG:
		globalLevel.Set(slog.LevelDebug)
	case ERROR, PROD:
		globalLevel.Set(slog.LevelError)
	default:
		globalLevel.Set(slog.LevelInfo)
	}
}

// GetMode возвращает текущий режим
func (l *Logger) GetMode() LogMode {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.mode
}

// Info логирует информационное сообщение
func (l *Logger) Info(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logger == nil {
		return
	}

	attrs := l.argsToAttrs(args)
	l.logger.LogAttrs(context.Background(), slog.LevelInfo, msg, attrs...)
}

// Debug логирует отладочное сообщение
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logger == nil {
		return
	}

	if l.mode != DEBUG {
		return
	}

	attrs := l.argsToAttrs(args)
	l.logger.LogAttrs(context.Background(), slog.LevelDebug, msg, attrs...)
}

// Warn логирует предупреждение
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logger == nil {
		return
	}

	if l.mode != DEBUG {
		return
	}

	attrs := l.argsToAttrs(args)
	l.logger.LogAttrs(context.Background(), slog.LevelWarn, msg, attrs...)
}

// Error логирует ошибку
func (l *Logger) Error(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logger == nil {
		return
	}

	attrs := l.argsToAttrs(args)
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)
}

// Fatal логирует критическую ошибку и паникует
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.logger == nil {
		panic(msg)
	}

	attrs := l.argsToAttrs(args)
	l.logger.LogAttrs(context.Background(), slog.LevelError, msg, attrs...)

	// Записываем в файл при панике
	l.writePanicLog(msg, args...)

	panic(msg)
}

// Must проверяет ошибку и паникует если она не nil
func (l *Logger) Must(err error, context string) {
	if err != nil {
		l.Error("Critical error", "context", context, "error", err.Error())

		// Записываем в файл при панике
		l.writePanicLog(context, "error", err.Error())

		panic(fmt.Sprintf("%s: %v", context, err))
	}
}

// MustNotNil проверяет что значение не nil
func (l *Logger) MustNotNil(value interface{}, context string) {
	if value == nil {
		l.Error("Nil value check failed", "context", context)

		// Записываем в файл при панике
		l.writePanicLog("Nil value", "context", context)

		panic(fmt.Sprintf("Value is nil: %s", context))
	}
}

// Close завершает работу логгера
func (l *Logger) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logger.Info("Logger closed", "timestamp", time.Now().Format(time.RFC3339))
}

// writePanicLog записывает информацию о панике в файл
func (l *Logger) writePanicLog(context string, args ...interface{}) {
	fileName := "prod_failed_logs.md"

	timestamp := time.Now().Format(time.RFC3339)

	// Конвертируем args в строку для причины
	reason := ""
	if len(args) > 0 {
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				reason += fmt.Sprintf("%v=%v ", args[i], args[i+1])
			}
		}
	}

	stackTrace := string(debug.Stack())

	// Форматируем содержимое файла
	content := fmt.Sprintf("# Production Failure Log\n\n**Timestamp:** %s\n\n**Context:** %s\n\n**Reason:** %s\n\n**Stack Trace:**\n```\n%s\n```\n\n---\n\n",
		timestamp, context, reason, stackTrace)

	// Если файл существует, добавляем в конец
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write panic log: %v\n", err)
		return
	}
	defer file.Close()

	_, _ = file.WriteString(content)
}

// argsToAttrs преобразует пары key-value в slog.Attr
func (l *Logger) argsToAttrs(args []interface{}) []slog.Attr {
	var attrs []slog.Attr
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key := fmt.Sprintf("%v", args[i])
			value := args[i+1]
			attrs = append(attrs, slog.Any(key, value))
		}
	}
	return attrs
}

// resetForTests сбрасывает синглтон (для тестирования)
func resetForTests() {
	once = sync.Once{}
	instance = nil
}
