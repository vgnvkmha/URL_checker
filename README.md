# MyApp / Service Checker

Простой HTTP-чекер / API для мониторинга доступности сайтов.  
Сервис проверяет доступность сайтов и возвращает результат через REST API.

---

## 📦 Установка

### Локально
```bash
git clone https://github.com/username/myapp.git
cd myapp
export DATABASE_DSN="postgres://user:pass@localhost:5432/dbname?sslmode=disable"
go run cmd/myapp/main.go
