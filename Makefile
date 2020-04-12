.PHONY: openapi
openapi:
	@./hack/openapi/openapi-generate.sh

.PHONY: lint
lint:
	@./hack/lint.sh

.PHONY: setup
setup:
	@go mod download
	@go mod verify

.PHONY: mod-download
mod-download:
	@go mod download
	
.PHONY: mod-tidy
mod-tidy:
	@go mod tidy

.PHONY: install
install:
	@go install github.com/gold-kou/go-housework/app/cmd/...

.PHONY: cross-install
cross-install:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -installsuffix cgo -ldflags '-w -extldflags "-static"' github.com/gold-kou/go-housework/app/cmd/...

.PHONY: test
# serviceとhandlerを同時に走らせると、PostgresqlへのTRUNCATEがバッティングしてしまうため、テストの同時実行数(-p)を1に設定
test:
	@go test -p=1 -covermode=count -coverprofile=cover.out github.com/gold-kou/go-housework/app/...
	@go tool cover -html=cover.out -o coverage.html

.PHONY: migrate
migrate:
	@./app/sql/migrate/migrate.linux-amd64_v425 -path ./app/sql/migrate -database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_NAME}?sslmode=disable up
