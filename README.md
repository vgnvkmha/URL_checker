# 🚀 URL Checker

Простой HTTP API-сервис для мониторинга доступности сайтов.  
Сервис периодически проверяет URL и сохраняет результаты проверок в базу данных.

---

## 📦 Установка и запуск

### 🔹 Требования

- Go 1.21+
- PostgreSQL 14+
- Созданная база данных

---

## ▶ Локальный запуск

1. Запустить PostgreSQL  
2. Создать базу данных  
3. Применить схему:

```bash
psql -U postgres -d your_db -f schema.sql
```
4. Создать .env файл в корне проекта:

```bash
# PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_name
DB_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password
REDIS_DB=0
```
5. Запустить приложение из корня проекта, чтобы переменные окружения считались
```bash
go run ./cmd/api/main.go
```

---

## 📡 API Endpoints

### 🔹 Работа с URL (Targets)

#### Получить список всех URL
`GET /targets`

#### Получить URL по ID
`GET /targets/{id}`

#### Добавить URL
`POST /targets`

Пример body:
```json
{
  "url": "https://example.com",
  "interval_sec": 10,
  "timeout_ms": 1000,
  "active": true
}
```

### Обновить параметры URL
`PATCH /targets/{id}`

Пример body:
```
{
  "interval": 60,
  "timeout": 2000,
  "active": false
}
```

### Удалить URL по ID
`DELETE /targets/{id}`

### Получить активные URL
`GET /targets/active`

### 🔹 Работа с Checks

### Получить последнюю проверку по ID
`GET /targets/{id}/status`

### Получить 5 последних проверок по ID
`GET /targets/{id}/checks`