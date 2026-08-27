# Production Failure Log

**Timestamp:** 2026-08-27T23:41:05+05:00

**Context:** test context

**Reason:** error=test error 

**Stack Trace:**
```
goroutine 25 [running]:
runtime/debug.Stack()
	/home/aibar/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64/src/runtime/debug/stack.go:26 +0x5e
order-service/internal/logger.(*Logger).writePanicLog(0x5b85fb?, {0x5b8eb0, 0xc}, {0x2e7f6b5f1e88, 0x2, 0x7322f97309e0?})
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:213 +0x9a
order-service/internal/logger.(*Logger).Must(0x2e7f6b5b6360, {0x5c88e8, 0x2e7f6b5ae640}, {0x5b8eb0, 0xc})
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:171 +0x1b5
order-service/internal/logger.TestLoggerMust(0x2e7f6b5e6fc8)
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_test.go:113 +0xfd
testing.tRunner(0x2e7f6b5e6fc8, 0x5c5dd0)
	/home/aibar/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 1
	/home/aibar/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64/src/testing/testing.go:2101 +0x4c5

```

---

# Production Failure Log

**Timestamp:** 2026-08-27T23:41:05+05:00

**Context:** Nil value

**Reason:** context=nil check 

**Stack Trace:**
```
goroutine 26 [running]:
runtime/debug.Stack()
	/home/aibar/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64/src/runtime/debug/stack.go:26 +0x5e
order-service/internal/logger.(*Logger).writePanicLog(0x5b8043?, {0x5b8031, 0x9}, {0x2e7f6b5f1ed0, 0x2, 0x2e7f6b5f1ef0?})
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:213 +0x9a
order-service/internal/logger.(*Logger).MustNotNil(0x2e7f6b5b63a0, {0x0?, 0x0?}, {0x5b8043, 0x9})
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_slog.go:183 +0x11b
order-service/internal/logger.TestLoggerMustNotNil(0x2e7f6b5e7208)
	/mnt/c/Users/ISO/Desktop/Aibar/Go/petProjects/mirco-service-project/order-service/internal/logger/logger_test.go:133 +0xc5
testing.tRunner(0x2e7f6b5e7208, 0x5c5dd8)
	/home/aibar/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64/src/testing/testing.go:2036 +0xea
created by testing.(*T).Run in goroutine 1
	/home/aibar/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.4.linux-amd64/src/testing/testing.go:2101 +0x4c5

```

---

