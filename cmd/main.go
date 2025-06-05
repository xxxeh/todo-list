package main

import (
	"fmt"
	"go1f/pkg/server"

	"github.com/joho/godotenv"
)

// Инициализация переменных окружения из файла .env
func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
