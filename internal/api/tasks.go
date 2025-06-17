package api

import (
	"net/http"

	"github.com/xxxeh/todo-list/internal/db"
)

type tasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.FormValue("search")

	tasks, err := db.Tasks(search, tasksLimit)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJson(w, tasksResp{Tasks: tasks}, http.StatusOK)
}
