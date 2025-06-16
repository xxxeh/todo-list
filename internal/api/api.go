package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const dateFormat string = "20060102"

func Init() *chi.Mux {
	r := chi.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("web")))
	r.Get("/api/nextdate", nextDateHandler)
	r.Post("/api/task", addTaskHandler)

	return r
}

func writeJson(w http.ResponseWriter, data any, status int) {
	resp, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write([]byte(resp))
}
