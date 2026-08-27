# Production Failure Log

**Timestamp:** 2026-08-27T23:53:35+05:00

**Context:** Failed to connect to RabbitMQ

**Reason:** error=dial tcp 127.0.0.1:5672: connect: connection refused 

**Stack Trace:**
```
goroutine 1 [running]:
runtime/debug.Stack()
	/home/aibar/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64/src/runtime/debug/stack.go:26 +0x5e
order-service/internal/logger.(*Logger).writePanicLog(0x1cd891ac140?, {0x10cbb05, 0x1d}, {0x1cd89107d28, 0x2, 0x0?})
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:213 +0x9a
order-service/internal/logger.(*Logger).Must(0x1cd8909f280, {0x11084a0, 0x1cd8910a9b0}, {0x10cbb05, 0x1d})
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:171 +0x1b5
main.main()
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/cmd/api/main.go:34 +0x2d9

```

---

# Production Failure Log

**Timestamp:** 2026-08-28T00:55:13+05:00

**Context:** Failed to connect to RabbitMQ

**Reason:** error=dial tcp 127.0.0.1:5672: connect: connection refused 

**Stack Trace:**
```
goroutine 1 [running]:
runtime/debug.Stack()
	/home/aibar/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64/src/runtime/debug/stack.go:26 +0x5e
order-service/internal/logger.(*Logger).writePanicLog(0x6dec3644200?, {0x10cbb05, 0x1d}, {0x6dec37a9d28, 0x2, 0x0?})
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:213 +0x9a
order-service/internal/logger.(*Logger).Must(0x6dec36b12a0, {0x11084a0, 0x6dec37ac9b0}, {0x10cbb05, 0x1d})
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:171 +0x1b5
main.main()
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/cmd/api/main.go:34 +0x2d9

```

---

# Production Failure Log

**Timestamp:** 2026-08-28T00:55:54+05:00

**Context:** Failed to connect to RabbitMQ

**Reason:** error=dial tcp [::1]:5672: connectex: No connection could be made because the target machine actively refused it. 

**Stack Trace:**
```
goroutine 1 [running]:
runtime/debug.Stack()
	C:/Program Files/Go/src/runtime/debug/stack.go:26 +0x5e
order-service/internal/logger.(*Logger).writePanicLog(0x3ab7d1d6c930?, {0x7ff629762b2e, 0x1d}, {0x3ab7d1e71d28, 0x2, 0x0?})
	C:/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:213 +0x9a
order-service/internal/logger.(*Logger).Must(0x3ab7d1d5b700, {0x7ff62979fac0, 0x3ab7d1e426e0}, {0x7ff629762b2e, 0x1d})
	C:/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:171 +0x1b5
main.main()
	C:/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/cmd/api/main.go:34 +0x2d9

```

---

# Production Failure Log

**Timestamp:** 2026-08-28T00:59:12+05:00

**Context:** Failed to connect to RabbitMQ

**Reason:** error=Exception (501) Reason: "EOF" 

**Stack Trace:**
```
goroutine 1 [running]:
runtime/debug.Stack()
	C:/Program Files/Go/src/runtime/debug/stack.go:26 +0x5e
order-service/internal/logger.(*Logger).writePanicLog(0x27242500efc0?, {0x7ff659ed2b2e, 0x1d}, {0x27242517dd28, 0x2, 0x0?})
	C:/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:213 +0x9a
order-service/internal/logger.(*Logger).Must(0x272425067700, {0x7ff659f0fb80, 0x27242520e180}, {0x7ff659ed2b2e, 0x1d})
	C:/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:171 +0x1b5
main.main()
	C:/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/cmd/api/main.go:34 +0x2d9

```

---

