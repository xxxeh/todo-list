package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xxxeh/todo-list/internal/api"
)

// Run запускает HTTP-сервер на порту, определённом в переменной окружения TODO_PORT.
// Функция инициализирует API и начинает прослушивание указанного порта для обработки входящих запросов.
func Run() error {
	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		return fmt.Errorf("Environment variable TODO_PORT is not defined")
	}

	r := api.Init()
	return http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
