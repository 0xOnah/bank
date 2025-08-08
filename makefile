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

## migrations name=<name>: Create a new migration
.PHONY: migrations
migrations:
	migrate create -ext=sql -dir=./internal/db/migrations -seq $(name)

## migrations-up: Apply up migrations
.PHONY: migrations-up
migrations-up:
	migrate -database=$(DSN) -path=./internal/db/migrations -verbose up

## migrations-down: Apply down migrations
.PHONY: migrations-down
migrations-down:
	migrate -database=$(DSN) -path=./internal/db/migrations -verbose down

## migrations-force version=<version>: Force migrations to a version
.PHONY: migrations-force
migrations-force:
	@echo "Forcing migration version $(version)"
	migrate -database=$(DSN) -path=./internal/db/migrations force $(version)


## test: Run all unit tests
.PHONY: test
test:
	go test -v -cover -count=1 ./...

## audit: Format, vet, test
.PHONY: audit
audit:
	@echo 'Formatting...'
	go fmt ./...

	# @echo 'Vetting...'
	# go vet ./...

	# @echo 'Running tests...'
	# go test -race -vet=off ./...

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

.PHONY: evans
evans:
	evans --host localhost -p 9090 -r repl