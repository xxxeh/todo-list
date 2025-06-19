package main

import (
	"fmt"

	"os"

	"github.com/xxxeh/todo-list/internal/db"
	"github.com/xxxeh/todo-list/internal/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if len(dbFile) == 0 {
		fmt.Println("Environment variable TODO_DBFILE is not defined")
	}

	err := db.Init(dbFile)
	if err != nil {
		panic(err)
	}

	err = server.Run()
	if err != nil {
		panic(err)
	}
}
