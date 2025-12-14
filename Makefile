.PHONY: dev up down logs migrate backend-test backend-build

dev: ## Запуск всего стека (frontend + backend + postgres)
	@docker compose up --build

up: ## Запуск в фоне
	@docker compose up -d --build

down: ## Остановка и удаление контейнеров
	@docker compose down

logs: ## Общие логи
	@docker compose logs -f backend frontend

migrate: ## Прогон всех миграций (up)
	@docker compose run --rm migrate

backend-test: ## Локальные go test (нужен установленный Go)
	cd backend && GOCACHE=$$(pwd)/.cache go test ./...

backend-build: ## Сборка Go бинарника
	cd backend && GOCACHE=$$(pwd)/.cache go build ./cmd/api
reboot:
	docker compose down -v && docker compose build --no-cache && docker compose up
db:
	docker compose exec postgres psql -U bot -d bothuini