# Этап сборки
FROM golang:1.23 AS builder

WORKDIR /app

# Скопировать go.mod/go.sum отдельно — кэшируются зависимости
COPY go.mod go.sum ./
RUN go mod download

# Скопировать весь исходный код
COPY . .

# Собрать бинарник
RUN go build -o worker-service ./cmd

# Финальный минимальный контейнер
FROM debian:bookworm-slim

WORKDIR /app

# Копируем бинарник из билда
COPY --from=builder /app/worker-service .
# Swagger JSON (чтобы UI видел spec)
COPY --from=builder /app/docs ./docs

EXPOSE 8080
CMD ["./worker-service"]