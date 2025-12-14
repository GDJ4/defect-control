# Defect Tracker Platform

Монолитный backend на Go (Gin) + SPA на Vue 3 для системы управления дефектами на строительных объектах. Репозиторий организован по каталогам `backend/` и `frontend/`, сверху — инфраструктура (Docker Compose, Makefile, env-файлы).

## Быстрый старт (Docker)

```bash
cp .env.example .env   # при необходимости поправь логины/пароли
make up                # docker compose up -d --build
```

Сервисы:

- API: http://localhost:8080 (проверьте `GET /healthz`)
- Frontend: http://localhost:5173
- PostgreSQL: localhost:5432 (`POSTGRES_USER/PASSWORD` из `.env`)

Миграции (включая сиды с демонстрационными пользователями/проектами) лежат в `backend/migrations`. Прогнать их можно командой `make migrate`. При необходимости добавлены:

- `003_refresh_tokens` — таблица refresh-токенов.
- `004_more_data` — дополнительные проекты/дефекты/комментарии для демонстрации.
- `005_demo_data` — расширенные демо-пользователи/проекты и набор свежих дефектов с комментариями/историей.

После запуска можно либо зарегистрировать нового пользователя запросом `POST /api/v1/auth/register`, либо авторизоваться одной из готовых учётных записей:

- Менеджер: `manager@systemacontrola.ru` / `password`
- Менеджер проектного офиса: `chief@systemacontrola.ru` / `password`
- Инженер: `engineer@systemacontrola.ru` / `password`
- Инженер ОТК: `qa@systemacontrola.ru` / `password`
- Наблюдатель по технике безопасности: `safety@systemacontrola.ru` / `password`

JWT-параметры настраиваются через `.env` (`JWT_SECRET`, `JWT_ACCESS_TTL`, `JWT_REFRESH_TTL`). За хранение файлов отвечает `STORAGE_DRIVER`:

- `local` (по умолчанию) — файлы падают в `storage/uploads`, скачивание идёт через `GET /api/v1/defects/:id/attachments/:attachmentId`.
- `s3` — используется встроенный MinIO (`docker-compose` поднимает `minio` + `minio-setup`). Настройки (`STORAGE_S3_ENDPOINT`, `STORAGE_S3_BUCKET`, `STORAGE_S3_ACCESS_KEY`, `...`) берутся из `.env`, а фронт получает готовые pre-signed URL из API.

## Структура

```
backend/   # Go monolith (cmd/, internal/, migrations/, Dockerfile)
frontend/  # Vue 3 + Vite SPA
docker-compose.yml
Makefile
.env.example
```

## Локальная разработка

- Backend:
  ```bash
  cd backend
  export GOCACHE=$(pwd)/.cache
  go run ./cmd/api
  ```
- Frontend:
  ```bash
  cd frontend
  npm install
  npm run dev -- --host
  ```

API базовый URL для фронта задаётся переменной `VITE_API_BASE_URL` (по умолчанию `http://localhost:8080/api/v1`).

## Что развёрнуто

- PostgreSQL 15 + миграции `001_init` и `002_seed_data` с базовыми пользователями/проектами/дефектами.
- Backend (Go + Gin + pgx + JWT): `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`, `POST /api/v1/auth/logout`, CRUD по проектам (`GET/POST /projects`), дефектам (`GET/POST /defects`, `GET /defects/:id`), комментариям/вложениям, смена статусов (`PATCH /defects/:id/status`) и скачивание файлов (`GET /defects/:id/attachments/:attachmentId`).
- Frontend (Vue 3 + Pinia): авторизация, создание проектов/дефектов, фильтрация, комментарии и загрузка/скачивание вложений в UI.
