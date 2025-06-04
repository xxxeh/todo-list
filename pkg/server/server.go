package server

import (
	"fmt"
	"net/http"
	"os"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		return fmt.Errorf("Environment variable TODO_PORT is not defined")
	}

	http.Handle("/", http.FileServer(http.Dir("web")))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
