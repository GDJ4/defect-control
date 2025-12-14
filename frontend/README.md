# Defect Tracker Frontend

Single Page Application на Vue 3 + Vite для работы с API учёта дефектов.

## Скрипты

```bash
npm install           # установка зависимостей
npm run dev -- --host # режим разработки (порт 5173)
npm run build         # production сборка (dist/)
npm run preview       # локальный предпросмотр сборки
```

API базовый URL задаётся переменной `VITE_API_BASE_URL` (см. `.env.example` в корне). По умолчанию берётся `http://localhost:8080/api/v1`.

## Структура

```
src/
  router/       # Маршруты приложения
  stores/       # Pinia (дефекты)
  services/     # Axios-клиент
  views/        # Основные страницы (Dashboard/Defects/Projects)
  components/   # UI-виджеты
```

Базовый UI подключается через Docker (см. `frontend/Dockerfile`) либо отдельным `npm run dev`.
