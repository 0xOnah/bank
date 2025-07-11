#Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 
RUN apk update && apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

COPY . .
RUN go build -o ./main cmd/main.go

#Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY internal/db/migrations ./migrations
COPY app.env .
COPY start.sh .
RUN chmod +x start.sh

EXPOSE 8080
CMD [ "/app/main" ] 
ENTRYPOINT [ "/app/start.sh" ]