up:
	@echo "запуск приложения в docker контейнере..."
	docker-compose up
	@echo "приложение запущено на порту 8080."
	@echo "для генерации токенов используйте:"
	@echo "  make token-admin"
	@echo "  make token-user"

down:
	@echo "Stopping docker container..."
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
	@echo "Running all tests..."
	@INTEGRATION_LOGS=0 go test -race -v -cover ./...

repositories_integration_test:
	@echo "Running integration tests for Postgres repos with logs..."
	@INTEGRATION_LOGS=1 go test -race -v -cover ./internal/infrastructure/postgres/...

#------------------------------------
generate:
	@echo "Generating code from OpenAPI spec..."
	oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
	go mod tidy
	@echo "Code generated"

generate-check:
	@echo "Checking if generated code is up to date..."
	oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml > /tmp/generated.go
	diff api/generated.go /tmp/generated.go || \
		(echo "Generated code is outdated. Run 'make generate'" && exit 1)
	@echo "Generated code is up to date"