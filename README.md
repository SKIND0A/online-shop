# ВСЕ ДЛЯ ПК

Каркас десктоп-приложения интернет-магазина на Go + PostgreSQL.

## Быстрый старт

1. Запусти PostgreSQL в Docker:
   - `docker compose up -d`
2. При необходимости создай локальный `.env` на основе `configs/.env.example`.
3. Запусти приложение:
   - `go run ./cmd/app`
4. Проверь healthcheck:
   - `http://localhost:8081/health`

Пример успешного ответа:

```json
{"success":true,"data":{"status":"ok","db":"ok"},"error":null}
```

## Переменные окружения

- `HTTP_ADDR` (по умолчанию `:8081`)
- `DATABASE_URL` (по умолчанию `postgres://postgres:1234@localhost:5432/postgres?sslmode=disable`)

## Структура

- `cmd/app` - точка входа приложения.
- `internal/config` - конфигурация.
- `internal/domain` - доменные сущности.
- `internal/usecase` - бизнес-логика.
- `internal/repository/postgres` - слой работы с PostgreSQL.
- `internal/delivery/http/handlers` - обработчики API.
- `migrations` - SQL-миграции.
- `ui` - десктоп-интерфейс.
- `configs` - шаблоны конфигурации.
- `scripts` - служебные скрипты.
- `docs` - документация проекта.

Файл `PLAN_VSE_DLYA_PK.md` содержит подробный план разработки.
