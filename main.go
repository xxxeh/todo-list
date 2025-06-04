package main

import (
	"fmt"
	"go1f/pkg/server"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	err := server.Run()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
