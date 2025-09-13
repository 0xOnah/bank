ARG SERVER_GO_VERSION=1.25-alpine
ARG OS_VERSION=3.21

#base
FROM golang:${SERVER_GO_VERSION} AS build-base
WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,id=gomods,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    go mod download

# dev-mode
FROM build-base AS dev

RUN go install github.com/air-verse/air@latest && \ 
    go install github.com/go-delve/delve/cmd/dlv@latest

ENV DSN=postgresql://postgres:secret@postgres-bankapp:5432/bank?sslmode=disable

CMD ["air", "-c", ".air.toml"]

#production-build
FROM build-base AS build-production

WORKDIR /app

COPY ./cmd/ ./cmd/
COPY .//doc/ ./doc/
COPY ./internal/config ./internal/config
COPY ./internal/db/migrate.go ./internal/db/migrate.go
COPY ./internal/db/migrations ./internal/db/migrations
COPY ./internal/db/repo ./internal/db/repo
COPY ./internal/db/sqlc ./internal/db/sqlc
COPY ./internal/entity ./internal/entity
COPY ./internal/sdk ./internal/sdk
COPY ./internal/service/ ./internal/service/
COPY ./internal/transport/grpc ./internal/transport/grpc
COPY ./internal/transport/http ./internal/transport/http
COPY ./internal/transport/sdk ./internal/transport/sdk
COPY ./app.env ./app.env
COPY ./pb ./pb

COPY ./start.sh ./start.sh
RUN chmod +x start.sh

RUN --mount=type=cache,id=gomods,target=/go/pkg/mod \
    --mount=type=cache,id=gobuild,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -x -ldflags="-s -w" -o ./bin/main cmd/main.go

#runtime
FROM alpine:${OS_VERSION} AS runtime
WORKDIR /app

COPY --from=build-production /app/bin/main /app/
COPY --from=build-production /app/app.env /app/app.env
COPY --from=build-production /app/start.sh /app/start.sh

USER 1001:1001

ENV LOG_LEVEL=info
ENV ENVIRONMENT=production
ENV GIN_MODE=release

EXPOSE 9090
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]