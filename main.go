package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/xxxeh/todo-list/internal/db"
	"github.com/xxxeh/todo-list/internal/server"
)

// init инициализирует переменные окружения из файла .env
func init() {
	if err := godotenv.Load(); err != nil {
		log.Panic("Не обнаружен файл .env")
	}
}

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if len(dbFile) == 0 {
		log.Panic("Не определена переменная окружения TODO_DBFILE")
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Panic(err)
	}

	err = server.Run()
	if err != nil {
		log.Panic(err)
	}
}
