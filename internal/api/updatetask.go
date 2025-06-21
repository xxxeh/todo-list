package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/xxxeh/todo-list/internal/db"
)

// tasksHandler обрабатывает запросы на изменение задачи.
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	err = db.UpdateTask(&task)

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}
