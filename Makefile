BOLD = \033[1m
CLEAR = \033[0m
CYAN = \033[36m

API_SPEC := https://github.com/apache/iceberg/raw/main/open-api/rest-catalog-open-api.yaml
API_YAML := api.yaml

help: ## Display this help
	@awk '\
		BEGIN {FS = ":.*##"; printf "Usage: make $(CYAN)<target>$(CLEAR)\n"} \
		/^[a-z\-0-9]+([\/]%)?([\/](%-)?[a-z\-0-9%]+)*:.*? ##/ { printf "  $(CYAN)%-15s$(CLEAR) %s\n", $$1, $$2 } \
		/^##@/ { printf "\n$(BOLD)%s$(CLEAR)\n", substr($$0, 5) }' \
		$(MAKEFILE_LIST)

##@: Development

.PHONY: build
build: ## Build this B
	@go build ./...

.PHONY: test
test: ## Run the smoke test
	@go test ./... -v -cover

catalog: ## Run the catalog server used by tests
	@docker run --rm -it -p "8181:8181" tabulario/iceberg-rest

##@: Release

.PHONY: release
release: ## Release (tag) a new version (VERSION=<version> make release)
	@git tag -sm "Release v$(VERSION)" "v$(VERSION)"
	@git push --atomic origin main "v$(VERSION)"

##@: Generate

generate: $(API_YAML) ## (Re)Generate the client
	@go mod tidy
	@go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml $(API_YAML)

.PHONY: update-api
update-api: ## Fetch the latest Iceberg OpenAPI Spec
	@curl -o $(API_YAML) -fsSL $(API_SPEC)
	@sed -i '' 's|{prefix}/||g' $(API_YAML)
	@sed -i '' "s|.*: '#/components/parameters/prefix'||g" $(API_YAML)
	@sed -i '' 's/    Namespace:/    Namespace:\n      x-go-name: Namespaces/' $(API_YAML)
