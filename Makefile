INTEGRATION_LOGS ?= 0

.PHONY: test repositories_integration_test
test:
	@echo "Running all tests..."
	@INTEGRATION_LOGS=0 go test -v -cover ./...

repositories_integration_test:
	@echo "Running integration tests for Postgres repos with logs..."
	@INTEGRATION_LOGS=1 go test -v -cover ./internal/infrastructure/postgres/...


generate:
	@echo "Generating code from OpenAPI spec..."
	oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
	go mod tidy
	@echo "Code generated"

.PHONY: generate-check
generate-check:
	@echo "Checking if generated code is up to date..."
	oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml > /tmp/generated.go
	diff api/generated.go /tmp/generated.go || \
		(echo "Generated code is outdated. Run 'make generate'" && exit 1)
	@echo "Generated code is up to date"