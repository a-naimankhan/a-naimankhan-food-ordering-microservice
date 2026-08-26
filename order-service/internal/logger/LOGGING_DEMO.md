# 🎬 Демонстрация работы логгера

## 📺 Живой пример логирования

### Сценарий 1: Запуск сервиса

**Консоль (stdout):**
```
🚀 Logger initialized in DEBUG mode | Logs: logs.md

🚀 [2026-08-27 01:04:00] INFO: Order Service starting...
🔍 [2026-08-27 01:04:00] DEBUG: Connecting to PostgreSQL database
✅ [2026-08-27 01:04:01] INFO: PostgreSQL connected successfully
🔍 [2026-08-27 01:04:01] DEBUG: Connecting to RabbitMQ broker
✅ [2026-08-27 01:04:02] INFO: RabbitMQ connected successfully
🔍 [2026-08-27 01:04:02] DEBUG: Initializing repository, service, and handler
✅ [2026-08-27 01:04:03] INFO: All components initialized successfully
🔍 [2026-08-27 01:04:04] DEBUG: Registering HTTP routes
✅ [2026-08-27 01:04:04] INFO: HTTP routes registered
🎧 [2026-08-27 01:04:04] INFO: Listening on :8080
```

**MD файл (logs.md):**
```markdown
# Logs - DEBUG Mode

**Started at:** 2026-08-27T01:04:00+05:00

- **[2026-08-27 01:04:00] INFO:** Order Service starting...
- **[2026-08-27 01:04:00] DEBUG:** Connecting to PostgreSQL database
- **[2026-08-27 01:04:01] INFO:** PostgreSQL connected successfully
- **[2026-08-27 01:04:01] DEBUG:** Connecting to RabbitMQ broker
- **[2026-08-27 01:04:02] INFO:** RabbitMQ connected successfully
- **[2026-08-27 01:04:02] DEBUG:** Initializing repository, service, and handler
- **[2026-08-27 01:04:03] INFO:** All components initialized successfully
- **[2026-08-27 01:04:04] DEBUG:** Registering HTTP routes
- **[2026-08-27 01:04:04] INFO:** HTTP routes registered
- **[2026-08-27 01:04:04] INFO:** Listening on :8080
```

---

### Сценарий 2: Создание заказа (успешно)

**HTTP Request:**
```bash
POST /api/v1/orders
{
  "customer_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 299.99,
  "status": "pending"
}
```

**Консоль (stdout) - полный путь запроса:**
```
🔍 [2026-08-27 01:04:10] DEBUG: HTTP: POST /api/v1/orders received | RemoteAddr=127.0.0.1
🔍 [2026-08-27 01:04:10] DEBUG: HTTP: Request parsed | CustomerID=550e8400-e29b-41d4-a716-446655440000 | Amount=299.99 | Status=pending
ℹ️ [2026-08-27 01:04:10] INFO: HTTP: Creating order | CustomerID=550e8400-e29b-41d4-a716-446655440000 | Amount=299.99

🔍 [2026-08-27 01:04:10] DEBUG: Service: Creating order | CustomerID=550e8400-e29b-41d4-a716-446655440000 | Amount=299.99 | Status=pending
🔍 [2026-08-27 01:04:10] DEBUG: Service: Generated new Order ID | ID=abc-def-123-456
🔍 [2026-08-27 01:04:10] DEBUG: Service: Set creation timestamp | CreatedAt=2026-08-27 01:04:10
ℹ️ [2026-08-27 01:04:10] INFO: Service: Validations passed, saving order to repository | ID=abc-def-123-456

🔍 [2026-08-27 01:04:10] DEBUG: Creating order | ID=abc-def-123-456 | Customer=550e8400-e29b-41d4-a716-446655440000 | Amount=299.99 | Status=pending
✅ [2026-08-27 01:04:10] INFO: Order created successfully | ID=abc-def-123-456 | Rows affected: 1

✅ [2026-08-27 01:04:10] INFO: Service: Order saved to database | ID=abc-def-123-456
🔍 [2026-08-27 01:04:10] DEBUG: Service: Publishing order.created event | ID=abc-def-123-456

🔍 [2026-08-27 01:04:10] DEBUG: RabbitMQ: Publishing event | Topic=order.created | Exchange=orders_events
🔍 [2026-08-27 01:04:10] DEBUG: RabbitMQ: Serialized payload | Topic=order.created | Size=512 bytes
✅ [2026-08-27 01:04:10] INFO: RabbitMQ: Event published successfully | Topic=order.created | Exchange=orders_events

✅ [2026-08-27 01:04:10] INFO: Service: Order created successfully | ID=abc-def-123-456 | CustomerID=550e8400-e29b-41d4-a716-446655440000
✅ [2026-08-27 01:04:10] INFO: HTTP: Order created successfully | ID=abc-def-123-456 | StatusCode=201
```

**Отклик HTTP:**
```json
{
  "id": "abc-def-123-456",
  "customer_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 299.99,
  "status": "pending",
  "created_at": "2026-08-27T01:04:10Z"
}
```

---

### Сценарий 3: Ошибка валидации

**HTTP Request (неправильный JSON):**
```bash
POST /api/v1/orders
{
  "amount": 299.99,
  "status": "pending"
  # Пропущен customer_id!
}
```

**Консоль (stdout):**
```
🔍 [2026-08-27 01:04:15] DEBUG: HTTP: POST /api/v1/orders received | RemoteAddr=127.0.0.1
❌ [2026-08-27 01:04:15] ERROR: HTTP: Invalid JSON request | Error: Key: 'OrderRequest.CustomerID' Error:Field validation for 'CustomerID' failed on the 'required' tag
```

**Отклик HTTP:**
```json
{
  "error": "Key: 'OrderRequest.CustomerID' Error:Field validation for 'CustomerID' failed on the 'required' tag"
}
```

---

### Сценарий 4: Получение заказа

**HTTP Request:**
```bash
GET /api/v1/orders/abc-def-123-456
```

**Консоль (stdout):**
```
🔍 [2026-08-27 01:04:20] DEBUG: HTTP: GET /api/v1/orders/:id received | ID=abc-def-123-456 | RemoteAddr=127.0.0.1
🔍 [2026-08-27 01:04:20] DEBUG: Service: Getting order | ID=abc-def-123-456

🔍 [2026-08-27 01:04:20] DEBUG: Fetching order from database | ID=abc-def-123-456
🔍 [2026-08-27 01:04:20] DEBUG: Order fetched successfully | ID=abc-def-123-456 | Customer=550e8400-e29b-41d4-a716-446655440000 | Amount=299.99

✅ [2026-08-27 01:04:20] INFO: Service: Order retrieved | ID=abc-def-123-456 | Status=pending | Amount=299.99
✅ [2026-08-27 01:04:20] INFO: HTTP: Order retrieved successfully | ID=abc-def-123-456 | StatusCode=200
```

**Отклик HTTP:**
```json
{
  "id": "abc-def-123-456",
  "customer_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 299.99,
  "status": "pending",
  "created_at": "2026-08-27T01:04:10Z"
}
```

---

### Сценарий 5: Ошибка - заказ не найден

**HTTP Request:**
```bash
GET /api/v1/orders/nonexistent-id
```

**Консоль (stdout):**
```
🔍 [2026-08-27 01:04:25] DEBUG: HTTP: GET /api/v1/orders/:id received | ID=nonexistent-id | RemoteAddr=127.0.0.1
🔍 [2026-08-27 01:04:25] DEBUG: Service: Getting order | ID=nonexistent-id

🔍 [2026-08-27 01:04:25] DEBUG: Fetching order from database | ID=nonexistent-id
❌ [2026-08-27 01:04:25] ERROR: Failed to fetch order | ID=nonexistent-id | Error: sql: no rows in result set

❌ [2026-08-27 01:04:25] ERROR: Service: Failed to get order | ID=nonexistent-id | Error: sql: no rows in result set
❌ [2026-08-27 01:04:25] ERROR: HTTP: Order not found | ID=nonexistent-id | Error: sql: no rows in result set
```

**Отклик HTTP:**
```json
{
  "error": "order not found"
}
```

---

### Сценарий 6: Ping запрос

**HTTP Request:**
```bash
GET /api/v1/ping
```

**Консоль (stdout):**
```
🔍 [2026-08-27 01:04:30] DEBUG: HTTP: GET /api/v1/ping received | RemoteAddr=127.0.0.1
🔍 [2026-08-27 01:04:30] DEBUG: HTTP: Ping response sent | StatusCode=200
```

**Отклик HTTP:**
```json
{
  "message": "pong"
}
```

---

### Сценарий 7: Завершение сервиса

**Консоль (stdout):**
```
📪 [2026-08-27 01:04:40] INFO: Shutdown signal received
🔒 [2026-08-27 01:04:40] INFO: Closing database connection
🔌 [2026-08-27 01:04:40] INFO: Closing RabbitMQ connection
✋ [2026-08-27 01:04:40] INFO: Server shutdown complete

👋 Logger closed at 2026-08-27 01:04:40
```

**MD файл после завершения:**
```markdown
...
- **[2026-08-27 01:04:40] INFO:** Shutdown signal received
- **[2026-08-27 01:04:40] INFO:** Closing database connection
- **[2026-08-27 01:04:40] INFO:** Closing RabbitMQ connection
- **[2026-08-27 01:04:40] INFO:** Server shutdown complete

**Finished at:** 2026-08-27T01:04:40+05:00
```

---

## 📊 Логирование по слоям

### Main (запуск/остановка)
```
🚀 Order Service starting
✅ PostgreSQL connected
✅ RabbitMQ connected
📪 Shutdown signal
```

### Delivery (HTTP)
```
🔍 POST /api/v1/orders received
ℹ️ Creating order | CustomerID=...
✅ Order created | StatusCode=201
❌ Invalid JSON request
```

### Service (бизнес логика)
```
🔍 Service: Creating order
🔍 Generated new Order ID
🔍 Validations passed
✅ Service: Order created successfully
❌ Service: Invalid amount
```

### Repository (БД)
```
🔍 Creating order | ID=... | Amount=...
✅ Order created successfully
❌ Failed to create order | Error=...
```

### Infrastructure (RabbitMQ)
```
🔍 RabbitMQ: Publishing event | Topic=order.created
✅ RabbitMQ: Event published successfully
❌ RabbitMQ: Failed to publish event
```

---

## 🎨 Цветовая легенда (эмодзи)

| Эмодзи | Уровень | Назначение |
|--------|---------|-----------|
| 🚀 | INFO | Запуск/инициализация |
| ℹ️ | INFO | Информационные сообщения |
| 🔍 | DEBUG | Отладочная информация |
| ⚠️ | WARN | Предупреждения |
| ❌ | ERROR | Ошибки |
| 💥 | FATAL | Критические ошибки (паника) |
| ✅ | SUCCESS | Успешная операция |
| 🎧 | SERVER | Сервер запущен |
| 📪 | SIGNAL | Сигнал завершения |
| 🔒 | CLOSE | Закрытие соединения |

---

## 📈 Структура логов для мониторинга

### Prometheus
```
# Количество ERROR логов
order_service_errors_total{level="error"} 15

# Количество создано заказов
order_service_orders_created_total 1023
```

### ELK Stack
```json
{
  "@timestamp": "2026-08-27T01:04:10Z",
  "level": "INFO",
  "message": "Order created successfully",
  "order_id": "abc-def-123-456",
  "customer_id": "550e8400-e29b-41d4-a716-446655440000",
  "amount": 299.99,
  "service": "order-service"
}
```

---

## 🔎 Поиск в логах

### Найти все ошибки
```bash
grep "❌" logs_test.md
grep "\[ERROR\]" logs_test.md
```

### Найти ошибки для определенного заказа
```bash
grep "abc-def-123-456" logs_test.md
```

### Найти время обработки запроса
```bash
grep -A5 "POST /api/v1/orders" logs_test.md
```

---

## 💾 Размер логов

- **Консоль:** 0 МБ (выводится по мере)
- **MD файл:** ~1 КБ на 1000 операций
- **После 8 часов работы:** ~100 МБ (в DEBUG режиме)

---

## 🎓 Лучшие практики чтения логов

1. **Начните с ошибок (❌):**
   ```bash
   grep "❌" logs_test.md | head -20
   ```

2. **Проследите цепочку для конкретного заказа:**
   ```bash
   grep "order-id" logs_test.md
   ```

3. **Проверьте время выполнения операции:**
   ```bash
   # Время между первым и последним логом операции
   grep "Creating order" logs_test.md
   grep "Order created" logs_test.md
   ```

4. **Анализируйте паттерны ошибок:**
   ```bash
   grep "❌" logs_test.md | wc -l  # Количество ошибок
   grep "Invalid" logs_test.md    # Ошибки валидации
   ```

---

## ✨ Перспективы

- ✅ Логирование в консоль и файл
- ⏳ Вращение логов (log rotation)
- ⏳ Отправка логов в Elasticsearch
- ⏳ Метрики производительности
- ⏳ Трассировка транзакций (tracing)
- ⏳ Алерты по ошибкам
