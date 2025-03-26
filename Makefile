.SILENT:

# Основные переменные
BINARY_NAME=main
CRON_BINARY=tg_tasks_added_cron
MIGRATION_BINARY=./migrations/main
SWAGGER_CMD=swag
GO_CMD=go
DOCKER_COMPOSE_DB=docker-compose-db.yml
DOCKER_COMPOSE_FULL=docker-compose.yml


# Запуск приложения
run:
	$(GO_CMD) run ./cmd/main.go

# Сборка приложения
build:
	$(GO_CMD) build -o $(BINARY_NAME) ./cmd/main.go

# Запуск миграций
migrate:
	$(GO_CMD) run $(MIGRATION_BINARY).go

# Генерация Swagger документации
swagger:
	$(SWAGGER_CMD) init -d ./internal/ -g ../cmd/main.go

# Установка Swagger (если не установлен)
install-swag:
	$(GO_CMD) install github.com/swaggo/swag/cmd/swag@latest

# Запуск тестов
test:
	$(GO_CMD) test ./...

# Сборка cron задачи
cron:
	$(GO_CMD) build -o $(CRON_BINARY) ./cmd/crons/tg_tasks_added_cron.go

# Запуск только БД в Docker
docker-up-db:
	docker compose -f $(DOCKER_COMPOSE_DB) up -d --build

# Полный запуск (приложение + БД) в Docker
docker-up:
	docker compose -f $(DOCKER_COMPOSE_FULL) up -d --build

