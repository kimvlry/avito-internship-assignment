up:
	@echo "запуск приложения в docker контейнере..."
	docker-compose up
	@echo "приложение запущено на порту 8080."
	@echo "для генерации токенов используйте:"
	@echo "  make token-admin"
	@echo "  make token-user"

down:
	@echo "остановка docker контейнера..."
	docker-compose down

#------------------------------------
token-admin:
	@echo "генерируется jwt для админа..."
	@go run cmd/token/main.go -user admin1 -admin

token-user:
	@echo "генерируется jwt для юзера..."
	@go run cmd/token/main.go -user user1

#------------------------------------
INTEGRATION_LOGS ?= 0

test:
	@echo "запуск всех тестов..."
	@INTEGRATION_LOGS=0 go test -race -v -cover ./...

repositories_integration_test:
	@echo "запуск интеграционного теста postgres репозиториев..."
	@INTEGRATION_LOGS=1 go test -race -v -cover ./internal/infrastructure/postgres/...

#------------------------------------
generate:
	@echo "генерация кода api с помощью oapi-codegen..."
	oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
	go mod tidy
	@echo "готово"