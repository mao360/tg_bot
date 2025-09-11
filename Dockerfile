# Используем официальный образ Go для сборки
FROM golang:1.23-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Используем минимальный образ Alpine для финального контейнера
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates

# Создаем пользователя для безопасности
RUN adduser -D -s /bin/sh appuser

# Устанавливаем рабочую директорию
WORKDIR /root/

# Копируем собранное приложение из builder стадии
COPY --from=builder /app/main .

# Меняем владельца файла
RUN chown appuser:appuser main

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт (если понадобится в будущем)
# EXPOSE 8080

# Запускаем приложение
CMD ["./main"]

