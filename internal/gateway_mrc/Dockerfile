# Build stage
FROM golang:1.23-alpine3.19 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем необходимые файлы для сборки
COPY pkg ./pkg
COPY cmd/gateway_mrc/main.go ./cmd/gateway_mrc/main.go
COPY internal/gateway_mrc ./internal/gateway_mrc
COPY go.mod go.sum ./

# Устанавливаем зависимости
RUN go mod tidy
RUN go mod download


# Выполняем сборку приложения
RUN go build -o /gateway_mrc ./cmd/gateway_mrc

# Runtime stage
FROM alpine:3.19

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранный бинарный файл из этапа сборки
COPY --from=builder /gateway_mrc /app/gateway_mrc

ENV APP_ENV_PATH=/app/app.env


# Открываем порт, если нужно
EXPOSE 8081

# Команда для запуска
CMD ["/app/gateway_mrc"]