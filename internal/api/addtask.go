package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/xxxeh/todo-list/internal/db"
)

// addTaskHandler обрабатывает запросы на добавление новой задачи.
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]int64{"id": id}, http.StatusCreated)

}
