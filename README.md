# Описание todo-list
Небольшое приложения для отслеживания задач. Позволяет просматривать, добавлять и редактировать задачи.

## Дополнительные задания повышенной сложности
Выполнены все:
* Работа с переменными окружения
* Правила повторения задач по дням недели и месяцам
* Получение списка ближайших задач с дополнительным фильром по дате или части названия/комментария
* Аутентификация
* Создание Docker-образа

## Рекомендации по запуску приложения локально
### Переменные окружения и .env
Приложение для корректной работы требует ряд переменных окружения.
В проекте используется пакет dotenv, все переменные окружения должны быть указаны в файле .env. Они инициализируются в функции init() в main.go

Список переменных:

* `TODO_PORT` - порт, который должен слушать сервер
* `TODO_DBFILE` - путь до файла базы данных
* `TODO_PASSWORD` - пароль для авторизации, зашифрованный методом sha256
* `TODO_SECRET_KEY` - секрет для подписания JSON Web токена

Пример файла `.env` (именно такой файл используется сейчас в проекте)

```
TODO_PORT=7540
TODO_DBFILE=data/scheduler.db
TODO_PASSWORD=8341a67de42ecf564e578e27e539e91cefeb336da97ad8f31d8e1887d82ab972
TODO_SECRET_KEY=f0ff0c9c364bc8955438c040e6643ca37945052f4ed3e843f1363e8a78615197
```

**Пароль для авторизации - VeryStrongPassword**

### Запуск приложения
Из директории проекта командой:
```bash
go run .
```

Можно сначала собрать исполняемый файл:
```bash
go build -o todo-list
```

Затем запустить его:
```bash
./todo-list
```

### Использование
После запуска приложения веб-интерфейс будет доступен по адресу `http://localhost:7540/`
Порт **7540** указан в текущем .env файле в переменной окружения `TODO_PORT`. Если вы изменили значение переменной, следует указать новый порт.
Для авторизации необходимо указать пароль, который соответствует паролю в `TODO_PASSWORD`. В текущем .env пароль **VeryStrongPassword**

## Тестирование
Для удобства тестирования файле `tests/settings.go` не использует переменные окружения.
Рекомендуется использовать текущий файл `tests/settings.go` из проекта:

```go
package tests

var Port = 7540
var DBFile = "../data/scheduler.db"
var FullNextDate = true
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJoYXNoIjoiODM0MWE2N2RlNDJlY2Y1NjRlNTc4ZTI3ZTUzOWU5MWNlZmViMzM2ZGE5N2FkOGYzMWQ4ZTE4ODdkODJhYjk3MiJ9.jW5yYhpcx1XxBd-8rhYkXhvMFiFeYLUkj8_xzOl_PSc`
```

### Запуск тестов
Тесты запускаются из директории проекта командой:
```bash
go test ./tests
```

## Cборка и запуск проекта через Doker
Для сборки Doker образа можно использовать следующую команду из директории проекта:
```
docker build -t todo-list .
```

После сборки для запуска контейнера (если не менялись файлы `.env` и `dokerfile`) можно испольовать команду:
```
docker run -p 7540:7540 todo-list
```

где `-p port1:port2` это параметр, который связывает порт контейнера с портом хост-системы.
В `dockerfile` командой `EXPOSE` должен быть указан порт, который слушает http-сервер.
`port1` - порт, с которого должно быть доступно приложение, `port2` - порт контейнера, на котором работает приложение.
Это необходимо учесть в случае изменения конфигурационных файлов приложения.