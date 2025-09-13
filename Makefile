-include app.env

.PHONY: help
help:
	@echo ''
	@echo 'Usage:'
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
	@echo ''

## run: Start the Go server
.PHONY: run
run:
	@echo starting the Go server
	go run cmd/main.go

## migrations name=<name>: Create a new migration
.PHONY: migrations
migrations:
	@echo creating migration file
	migrate create -ext=sql -dir=./internal/db/migrations -seq $(name)

## migrations-up: Apply up migrations
.PHONY: migrations-up
migrations-up:
	@echo running up migrations
	migrate -database=$(DSN) -path=./internal/db/migrations -verbose up

## migrations-down: Apply down migrations
.PHONY: migrations-down
migrations-down:
	@echo running down migrations
	migrate -database=$(DSN) -path=./internal/db/migrations -verbose down

## migrations-force version=<version>: Force migrations to a version
.PHONY: migrations-force
migrations-force:
	@echo "Forcing migration version $(version)"
	migrate -database=$(DSN) -path=./internal/db/migrations force $(version)

## test: Run all unit tests
.PHONY: test
test:
	@echo running all unit tests
	go test -v -cover -count=1 ./...

## audit: Format, lint, test
.PHONY: audit
audit:
	@echo 'Formatting...'
	go fmt ./...

	@echo 'Linting...'
	golangci-lint run

	# @echo 'Running tests...'
	# go test -race -vet=off ./...

## vendor: Tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo Vendoring...
	go mod tidy
	go mod verify
	go mod vendor

## build: Build the Go binary
.PHONY: build
build:
	@echo Building the Go binary
	go build -o bin/app cmd/main.go

## clean: Remove build artifacts
.PHONY: clean
clean:
	@echo cleaning our binarys
	rm -rf ./bin/*

## sqlc: Generate code with sqlc
.PHONY: sqlc
sqlc:
	sqlc generate

## mock: filename=<name> interface-name=<iface> Generate mocks
.PHONY: mock
mock:
	mockgen -package mockdb -destination internal/db/mock/$(filename).go github.com/0xOnah/bank/internal/service $(interface-name)

## proto: generate proto files with grpc gateway included using relative path
.PHONY: proto
proto:
	rm -rf pb/*.go
	rm -rf doc/swagger/*.swagger.json	
	@protoc --proto_path=proto \
		--go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		 -I . --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge \
		proto/*.proto

## evans: run the evans grpc client
.PHONY: evans
evans:
	evans --host localhost -p 9090 -r repl

## redis: run the redis client   
.PHONY: redis
redis:
	docker run --name redis-bankapp -p 6379:6379 -d  redis:8.2.0-alpine

.PHONY: compose-build
compose-build:
	docker compose build --no-cache

.PHONY: compose-up 
compose-up:
	docker compose up -f docker-compose.yaml

.PHONY: compose-debug
compose-debug:
	docker compose -f docker-compose.yaml -f docker-compose-debug.yaml up

.PHONY: compose-down
compose-down:
	docker compose down

.PHONY: compose-test
compose-test:
	docker compose -f docker-compose.yaml -f docker-compose-test.yaml run --build simplebank


######################possible deletion#################################
## integration-test: Run integration tests
.PHONY: integration-test
integration-test:
	@echo starting integration test
	docker compose up -d db
	go test -tags=integration -count=1 -v ./...

## e2e: Run end-to-end tests
.PHONY: e2e
e2e:
	@echo starting end to end test
	docker compose up -d --build
	go test -tags=e2e -count=1 -v ./...