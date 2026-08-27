package logger

import "time"

func DoNothing() {
	time.Sleep(1 * time.Millisecond)
}

//
//import (
//	"fmt"
//	"os"
//	"sync"
//	"time"
//)
//
//type LogMode string
//
//const (
//	DEBUG LogMode = "DEBUG"
//	ERROR LogMode = "ERROR"
//	PROD  LogMode = "PROD"
//)
//
//type Logger struct {
//	mode     LogMode
//	filePath string
//	file     *os.File
//	mu       sync.Mutex
//}
//
//var (
//	instance *Logger
//	once     sync.Once
//	mu       sync.Mutex
//)
//
//// GetLogger возвращает глобальный экземпляр логгера (singleton)
//func GetLogger() *Logger {
//	once.Do(func() {
//		instance = &Logger{
//			mode:     PROD,
//			filePath: "logs_test.md",
//		}
//	})
//	return instance
//}
//
//// Init инициализирует логгер с режимом и путём к файлу
//func Init(mode LogMode, filePath string) *Logger {
//	once.Do(func() {
//		file, err := os.Create(filePath)
//		if err != nil {
//			panic(fmt.Sprintf("Failed to create log file: %v", err))
//		}
//
//		instance = &Logger{
//			mode:     mode,
//			filePath: filePath,
//			file:     file,
//		}
//
//		header := fmt.Sprintf("# Logs - %s Mode\n\n**Started at:** %s\n\n", mode, time.Now().Format(time.RFC3339))
//		_, err = instance.file.WriteString(header)
//		if err != nil {
//			panic(fmt.Sprintf("Failed to write log file header: %v", err))
//		}
//
//		fmt.Printf("Logger initialized in %s mode | Logs: %s\n\n", mode, filePath)
//	})
//	return instance
//}
//
//// resetForTests сбрасывает синглтон (только для тестирования)
//func resetForTests() {
//	mu.Lock()
//	defer mu.Unlock()
//	once = sync.Once{}
//	instance = nil
//}
//
//// SetMode устанавливает режим логирования
//func (l *Logger) SetMode(mode LogMode) {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	l.mode = mode
//}
//
//// GetMode возвращает текущий режим
//func (l *Logger) GetMode() LogMode {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//	return l.mode
//}
//
//// writeToAll пишет одновременно в консоль и в файл
//func (l *Logger) writeToAll(level string, formattedEntry string) {
//	// Консоль
//	fmt.Print(formattedEntry)
//
//	// Файл
//	if l.file != nil {
//		l.file.WriteString(formattedEntry)
//	}
//}
//
//// log пишет сообщение в файл и консоль в зависимости от режима
//func (l *Logger) log(level string, message string) {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//
//	shouldLog := false
//	switch l.mode {
//	case DEBUG:
//		shouldLog = true
//	case ERROR:
//		shouldLog = level == "ERROR" || level == "FATAL"
//	case PROD:
//		shouldLog = level == "ERROR" || level == "FATAL"
//	}
//
//	if !shouldLog {
//		return
//	}
//
//	timestamp := time.Now().Format("2006-01-02 15:04:05")
//
//	// Emoji для разных уровней
//	emoji := ""
//	switch level {
//	case "INFO":
//		emoji = "ℹ️ "
//	case "DEBUG":
//		emoji = "🔍 "
//	case "WARN":
//		emoji = "⚠️ "
//	case "ERROR":
//		emoji = "❌ "
//	case "FATAL":
//		emoji = "💥 "
//	}
//
//	// Формат для консоли
//	consoleEntry := fmt.Sprintf("%s[%s] %s: %s\n", emoji, timestamp, level, message)
//
//	// Формат для MD файла
//	mdEntry := fmt.Sprintf("- **[%s] %s:** %s\n", timestamp, level, message)
//
//	// Пишем в оба места
//	fmt.Print(consoleEntry)
//	if l.file != nil {
//		l.file.WriteString(mdEntry)
//	}
//}
//
//// Info логирует информационное сообщение (только в DEBUG режиме)
//func (l *Logger) Info(message string, args ...interface{}) {
//	msg := fmt.Sprintf(message, args...)
//	l.log("INFO", msg)
//}
//
//// Debug логирует сообщение отладки (только в DEBUG режиме)
//func (l *Logger) Debug(message string, args ...interface{}) {
//	msg := fmt.Sprintf(message, args...)
//	l.log("DEBUG", msg)
//}
//
//// Error логирует сообщение об ошибке
//func (l *Logger) Error(message string, args ...interface{}) {
//	msg := fmt.Sprintf(message, args...)
//	l.log("ERROR", msg)
//}
//
//// Warn логирует предупреждение (только в DEBUG режиме)
//func (l *Logger) Warn(message string, args ...interface{}) {
//	msg := fmt.Sprintf(message, args...)
//	l.log("WARN", msg)
//}
//
//// Fatal логирует критическую ошибку и паникует
//func (l *Logger) Fatal(message string, args ...interface{}) {
//	msg := fmt.Sprintf(message, args...)
//	l.log("FATAL", msg)
//	panic(msg)
//}
//
//// Must проверяет ошибку и паникует если она не nil
//// Если логгер в режиме DEBUG, перед паникой выводит стек вызовов
//func (l *Logger) Must(err error, context string) {
//	if err != nil {
//		errorMsg := fmt.Sprintf("%s: %v", context, err)
//		l.Error("%s", errorMsg)
//
//		if l.mode == DEBUG {
//			panic(fmt.Sprintf("CRITICAL ERROR in %s: %v", context, err))
//		}
//		panic(errorMsg)
//	}
//}
//
//// MustNotNil проверяет что значение не nil, иначе паникует
//func (l *Logger) MustNotNil(value interface{}, context string) {
//	if value == nil {
//		errorMsg := fmt.Sprintf("Value is nil: %s", context)
//		l.Error("%s", errorMsg)
//		panic(errorMsg)
//	}
//}
//
//// Close записывает завершающую информацию
//func (l *Logger) Close() {
//	l.mu.Lock()
//	defer l.mu.Unlock()
//
//	timestamp := time.Now().Format("2006-01-02 15:04:05")
//	footer := fmt.Sprintf("\n**Finished at:** %s\n", timestamp)
//
//	if l.file != nil {
//		l.file.WriteString(footer)
//		l.file.Close()
//	}
//
//	fmt.Printf("\n👋 Logger closed at %s\n", timestamp)
//}
