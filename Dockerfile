## Этап 1: сборка
#FROM golang:1.23.3-alpine3.20 AS builder
#
## Устанавливаем рабочую директорию
#WORKDIR /app
#
## Копируем файлы зависимостей
#COPY go.mod go.sum ./
#
## Устанавливаем зависимости
#RUN go mod download
#
## Копируем оставшиеся файлы
#COPY . .
#
## Устанавливаем нужные библиотеки
#RUN apk add --no-cache libc6-compat
#
## Собираем бинарник
#RUN go build -o main_app ./cmd/app/main.go
#
## Этап 2: запуск
#FROM alpine:latest
#
## Устанавливаем рабочую директорию
#WORKDIR /root/
#
## Копируем бинарник из этапа сборки
#COPY --from=builder /app/main_app .
#COPY config/config.yaml ./config/config.yaml
#COPY data/order.json ./data/order.json
#RUN chmod +x main_app
#
## Указываем порт, который использует приложение
#EXPOSE 8000
#
## Запускаем приложение
#CMD ["./main_app"]
