# Сборка приложения на основе образа golang
FROM golang:1.24.3 AS builder

WORKDIR /app

COPY . /app/

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o todo-list

# Создание образа на основе ubuntu
FROM ubuntu:latest

WORKDIR /app

COPY .env /app/

COPY web /app/web

COPY --from=builder /app/todo-list .

EXPOSE 7540

CMD ["./todo-list"]