## Include environment variables
include app.env

## help: Show this help message
.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ''

## confirm: Ask before continuing
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ "$$ans" = y ]


## run/server: Start the Go server
.PHONY: run/server
run/server:
	go run cmd/main.go

## run/docker: Start Docker Compose
.PHONY: run/docker
run/docker:
	docker compose up --build

## integration-test: Run integration tests
.PHONY: integration-test
integration-test:
	docker compose up -d db
	go test -tags=integration -count=1 -v ./...

## e2e: Run end-to-end tests
.PHONY: e2e
e2e:
	docker compose up -d --build
	go test -tags=e2e -count=1 -v ./...

## db/migrations name=<name>: Create a new migration
.PHONY: db/migrations
db/migrations:
	migrate create -ext=sql -dir=./db/migrations -seq $(name)

## db/migrations/up: Apply up migrations
.PHONY: db/migrations/up
db/migrations/up:
	migrate -database=$(DSN) -path=./internal/db/migrations -verbose up

## db/migrations/down: Apply down migrations
.PHONY: db/migrations/down
db/migrations/down:
	migrate -database=$(DSN) -path=./internal/db/migrations -verbose down

## db/migrations/force version=<version>: Force migrations to a version
.PHONY: db/migrations/force
db/migrations/force:
	@echo "Forcing migration version $(version)"
	migrate -database=$(DSN) -path=./db/migrations force $(version)


## test: Run all unit tests
.PHONY: test
test:
	go test -v -cover -count=1 ./...

## lint: Run golangci-lint
.PHONY: lint
lint:
	golangci-lint run

## audit: Format, vet, test
.PHONY: audit
audit: vendor
	@echo 'Formatting...'
	go fmt ./...

	@echo 'Vetting...'
	go vet ./...

	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: Tidy and vendor dependencies
.PHONY: vendor
vendor:
	go mod tidy
	go mod verify
	go mod vendor


## build: Build the Go binary
.PHONY: build
build:
	go build -o bin/app cmd/main.go

## clean: Remove build artifacts
.PHONY: clean
clean:
	rm -rf ./bin/*


## sqlc: Generate code with sqlc
.PHONY: sqlc
sqlc:
	sqlc generate

## mock filename=<name> interface-name=<iface>: Generate mocks
.PHONY: mock
mock:
	mockgen -package mockdb -destination internal/db/mock/$(filename).go github.com/onahvictor/bank/internal/service $(interface-name)
