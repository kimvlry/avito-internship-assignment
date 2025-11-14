.PHONY: generate
generate:
	@echo "Generating code from OpenAPI spec..."
	oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml
	go mod tidy
	@echo "Code generated"

.PHONY: generate-check
generate-check:
	@echo "Checking if generated code is up to date..."
	oapi-codegen -config api/oapi-codegen.yaml api/openapi.yaml > /tmp/generated.go
	diff internal/delivery/http/generated.go /tmp/generated.go || \
		(echo "Generated code is outdated. Run 'make generate'" && exit 1)
	@echo "Generated code is up to date"