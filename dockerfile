# Этап сборки
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем файлы go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код и компилируем приложение
COPY . . 
RUN GOOS=linux GOARCH=amd64 go build -o /app/srest-api ./cmd/sapi

# Финальный образ
FROM alpine:latest

# Устанавливаем Nginx
RUN apk add --no-cache nginx

# Копируем собранное приложение
COPY --from=builder /app/srest-api /usr/local/bin/srest-api

# Копируем конфигурационные файлы приложения
COPY config /usr/local/bin/config

# Копируем миграции
COPY internal/database/migrations /usr/local/bin/internal/database/migrations

# Копируем конфигурацию Nginx
COPY ./nginx.conf /etc/nginx/nginx.conf

# Устанавливаем переменную окружения
ENV CONFIG_PATH=/usr/local/bin/config/prod.yaml

# Открываем порты для приложения и Nginx
EXPOSE 80
EXPOSE 443

# Запускаем Nginx и приложение
CMD ["sh", "-c", "nginx && /usr/local/bin/srest-api"]
