# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и компилируем приложение
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /app/srest-api ./cmd/sapi

# Финальный образ
FROM alpine:latest

# Копируем собранное приложение
COPY --from=builder /app/srest-api /usr/local/bin/srest-api

# Копируем конфигурационные файлы приложения
COPY config /usr/local/bin/config

# Устанавливаем переменную окружения
ENV CONFIG_PATH=/usr/local/bin/config/local.yaml

# Устанавливаем рабочую директорию
WORKDIR /usr/local/bin

# Порт для приложения
EXPOSE 8082

# Запускаем приложение
CMD ["srest-api"]
