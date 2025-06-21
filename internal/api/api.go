package api

//Файл содержит функцию инициализации маршрутизатора и регистрации хэндлеров,
//а также несколько вспомогательных функций, которые используются в нескольких хэндлерах.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/xxxeh/todo-list/internal/db"
)

const (
	dateFormat string = "20060102"
	tasksLimit int    = 30
)

// Init инициализирует и настраивает HTTP-сервер с маршрутами для работы с задачами.
//
// Возвращаемое значение:
//
//	*chi.Mux - маршрутизатор chi с зарегистрированными обработчиками маршрутов.
func Init() *chi.Mux {
	r := chi.NewRouter()
	r.Use(logger)

	r.Handle("/*", http.FileServer(http.Dir("web")))
	r.Get("/api/nextdate", nextDateHandler)
	r.Get("/api/tasks", auth(tasksHandler))
	r.Get("/api/task", auth(getTaskHandler))
	r.Put("/api/task", auth(updateTaskHandler))
	r.Post("/api/task", auth(addTaskHandler))
	r.Post("/api/task/done", auth(completeTaskHandler))
	r.Post("/api/signin", authHandler)
	r.Delete("/api/task", auth(deleteTaskHandler))

	return r
}

// writeJson записывает данные в формате json в ответ HTTP-сервера.
//
// Параметры:
//
//	 w - http.ResponseWriter, используемый для записи ответа клиенту.
//		data - данные, которые будут преобразованы в json и записаны в тело ответа.
//		status - код ответа HTTP-сервера, который будет записан в заголовок.
func writeJson(w http.ResponseWriter, data any, status int) {
	resp, err := json.Marshal(data)
	if err != nil {
		//На данном этапе ошибки возникнуть не можем, поэтому опропускаем ее в возвращаемом параметре
		resp, _ = json.Marshal(map[string]string{"error": err.Error()})
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	w.Write(resp)

	log.Printf("Sending response with status %d - %s", status, string(resp))
}

// checkDate рассчитывает и сохраняет корректную дату, в которую должна быть назначена задача.
//
// Параметры:
//
//	task - указатель на структуру Task, содержащую данные задачи.
//
// Возвращаемые значения:
//
//	error - ошибка, которая могла возникнуть в ходе работы.
func checkDate(task *db.Task) error {
	now := time.Now()
	if len(task.Date) == 0 {
		//Если дата изначально не указана в задаче, записываем текущую и возвращаем nil в качестве ошибки.
		task.Date = now.Format(dateFormat)
		return nil
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("Неверный формат даты")
	}

	if after(now, t) {
		if len(task.Repeat) == 0 {
			//Если текущая дата больше чем дата в задаче и нет условия повторения, то записываем текущую дату.
			task.Date = now.Format(dateFormat)
		} else {
			//Если текущая дата больше чем дата в задаче и есть условие повторения, то вычисляем и новую дату.
			task.Date, err = NextDate(now, task.Date, task.Repeat)
		}
	}
	return err
}

// after проверяет, является ли первая дата (date1) более поздней, чем вторая дата (date2).
// Даты сравниваются с усечением до начала суток, что позволяет игнорировать время и сравнивать только даты.
//
// Параметры:
//
//	date1 - первая дата для сравнения.
//	date2 - вторая дата для сравнения.
//
// Возвращаемое значение:
//
//	bool - true, если date1 позже date2; false в противном случае.
func after(date1, date2 time.Time) bool {
	return date1.Truncate(24 * time.Hour).After(date2.Truncate(24 * time.Hour))
}
