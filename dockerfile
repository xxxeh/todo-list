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

#Определяем переменные окружения.
#Переменные определенные здесь имеют приоритет выше, чем переменные из .env,
#т.к godotenv.Load() не перезаписывает существующие переменные окружения.
ENV TODO_PORT=7540
ENV TODO_DBFILE=data/scheduler.db
ENV TODO_PASSWORD=9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08
ENV TODO_SECRET_KEY=f0ff0c9c364bc8955438c040e6643ca37945052f4ed3e843f1363e8a78615197

EXPOSE 7540

CMD ["./todo-list"]