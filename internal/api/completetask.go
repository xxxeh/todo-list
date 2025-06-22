package api

import (
	"net/http"
	"time"

	"github.com/xxxeh/todo-list/internal/db"
)

// completeTaskHandler обрабатывает запрос на завершение задачи.
// В зависимости от наличия условия повторения задачи, задача либо удаляется, либо обновляется с новой датой.
func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if len(id) == 0 {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	if len(task.Repeat) == 0 {
		err := db.DeleteTask(task.ID)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		} else {
			writeJson(w, struct{}{}, http.StatusOK)
		}
		return
	}

	task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	err = db.UpdateDate(task.Date, task.ID)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, struct{}{}, http.StatusOK)
}
