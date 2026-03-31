# ВСЕ ДЛЯ ПК

Каркас проекта (без кода) для десктоп-приложения интернет-магазина на Go + PostgreSQL.

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
