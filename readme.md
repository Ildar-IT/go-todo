# Todo REST API

Это простое REST API приложение для управления задачами с JWT-авторизацией.

## Установка

Клонируйте репозиторий:

```bash
git clone https://github.com/Ildar-IT/go-todo.git todo
```

## Запуск проекта

1. Установите зависимости:

```
go mod download
```

2. Запустите приложение:

```
go run ./cmd/main.go
```

Или соберите и запустите бинарный файл:

```
go build -o main ./cmd/main.go && ./main
```

## Миграция

1. Запустите миграции после поднятия базы

```
go run migrations/main.go
```

2. Или череp билд:

```
go build -o ./migrations/main ./migrations/main.go && ./migrations/main
```

## Документация API (Swagger)

1. Установите Swagger:

bash

```
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Сгенерируйте документацию:

```
swag init -d ./internal/ -g ../cmd/main.go
```

3. Запустите приложение и откройте в браузере:

```
http://localhost:${httpPort}/swagger/index.html
```

## Настройка Cron

1. Установите cron:

```
sudo apt-get install cron
```

2. Соберите cron-задачу:

```
go build -o tg_tasks_added_cron ./cmd/crons/tg_tasks_added_cron.go
```

3. Добавьте задачу в cron (редактируйте через crontab -e):

```
*/1 * * * * /todo/tg_tasks_added_cron --envPath=/todo/.env
```

## Запуск тестов

```
go test ./...
```

## Запуск через Docker

### Только PostgreSQL

```
docker compose -f 'docker-compose-db.yml' up -d --build
```

### Полный запуск (приложение + БД)

```
docker compose -f 'docker-compose.yml' up -d --build
```
