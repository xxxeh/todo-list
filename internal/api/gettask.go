package api

import (
	"net/http"

	"github.com/xxxeh/todo-list/internal/db"
)

// getTaskHandler обрабатывает запрос на получение задачи по идентификатору.
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if len(id) == 0 {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if task == nil {
			writeJson(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		} else {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		}
		return
	}

	writeJson(w, task, http.StatusOK)
}
