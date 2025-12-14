# Система управления дефектами

Стартовая заготовка монолитного backend-а на Go для системы учёта дефектов на строительных объектах. Код расположен в этом каталоге (`/backend`), что упрощает совместную работу с фронтом и инфраструктурой.

## Быстрый старт

```bash
cd backend
export GOCACHE=$(pwd)/.cache # чтобы кэш помещался в рабочую директорию
DATABASE_DSN="postgres://defect:defect@localhost:5432/defect_db?sslmode=disable" go run ./cmd/api
```

HTTP сервер поднимется на `http://localhost:8080`, проверка — `GET /healthz` или `GET /api/v1/ping`. Endpoint `GET /api/v1/defects` берёт данные из Postgres (фильтры `status`, `priority`, `limit`).

### Миграции

Полный SQL init лежит в `migrations/001_init.up.sql` (понижающая миграция — `001_init.down.sql`). В `002_seed_data` добавлены базовые пользователи/проекты/дефекты, в `004_more_data` и `005_demo_data` — расширенные данные для витрины (доп. проекты, комментарии, история). Для прогонки можно использовать `golang-migrate`:

```bash
migrate -path migrations -database "$DATABASE_DSN" up
```

## Текущая структура

```
cmd/api              # точка входа
internal/pkg/config  # работа с ENV
internal/pkg/logger  # инициализация zap
internal/pkg/server  # обёртка над http.Server
internal/transport   # HTTP-роуты (Gin)
migrations           # SQL init (users/projects/defects/...)
```

## Ближайшие шаги

1. Вынести хранение вложений в MinIO/S3 и отдавать pre-signed URL (сейчас файлы лежат локально в `storage/uploads`).
2. Ввести аудит действий (event log) и уведомления (email/Telegram) при изменении статусов дефектов.
3. Покрыть сервисы unit/integration тестами и добавить e2e (Cypress) для фронта.

### REST API (текущая реализация)

- `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`, `POST /api/v1/auth/logout`, `POST /api/v1/auth/password`
- `GET /api/v1/projects`, `POST /api/v1/projects`
- `GET /api/v1/defects`, `POST /api/v1/defects`, `GET /api/v1/defects/:id`, `PATCH /api/v1/defects/:id/status`
- `GET /api/v1/defects/:id/comments`, `POST /api/v1/defects/:id/comments`
- `POST /api/v1/defects/:id/attachments`, `GET /api/v1/defects/:id/attachments/:attachmentId`
