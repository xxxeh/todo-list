package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в базу данных.
// Параметры:
//
//	task - указатель на структуру Task, содержащую данные задачи.
//
// Возвращаемые значения:
//
//	int64 - идентификатор добавленной задачи.
//	error - ошибка, которая могла возникнуть в ходе работы.
func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// Tasks выполняет поиск задач в базе данных.
// Параметры:
//
//	search - параметр, по которому фильтруются задачи (дата или часть названия/комментария задачи). Может быть пустой строкой, если фильтрация не требуется.
//	limit - максимальное количество задач в результате.
//
// Возвращаемые значения:
//
//	[]*Task - список найденных задач.
//	error - ошибка, которая могла возникнуть в ходе работы.
func Tasks(search string, limit int) ([]*Task, error) {
	var tasks []*Task

	query := `SELECT * FROM scheduler ORDER BY date LIMIT :limit`
	if len(search) > 0 {
		date, err := time.Parse("02.01.2006", search)
		if err == nil {
			search = date.Format("20060102")
			query = `SELECT * FROM scheduler WHERE date = :search ORDER BY date LIMIT :limit`
		} else {
			query = `SELECT * FROM scheduler WHERE title LIKE '%' || :search || '%' OR comment LIKE '%' || :search || '%' ORDER BY date LIMIT :limit`
		}
	}

	rows, err := db.Query(query, sql.Named("limit", limit), sql.Named("search", search))
	if err != nil {
		return tasks, err
	}

	defer rows.Close()

	for rows.Next() {
		t := &Task{}
		err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}

	if tasks == nil {
		return []*Task{}, nil
	}

	return tasks, nil
}

// GetTask выполняет поиск задачи в базе данных по заданному идентификатору.
// Параметры:
//
//	id - идентификатор задачи.
//
// Возвращеаемы значения:
//
//	*Task - найденная задача.
//	error - ошибка, которая могла возникнуть в ходе работы.
func GetTask(id string) (*Task, error) {
	t := &Task{}

	query := `SELECT * FROM scheduler WHERE id = :id`
	row := db.QueryRow(query, sql.Named("id", id))
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)

	return t, err
}

// UpdateTask обновляет информацию о задаче в базе данных.
// Параметры:
//
//	task - указатель на структуру Task, содержащую данные задачи.
//
// Возвращаемые значения:
//
//	error - ошибка, которая могла возникнуть в ходе работы.
func UpadteTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

// DeleteTask удалаяет задачу из базы данных.
// Параметры:
//
//	id - идентификатор задачи.
//
// Возвращаемые значения:
//
//	error - ошибка, которая могла возникнуть в ходе работы.
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = :id`
	res, err := db.Exec(query, sql.Named("id", id))

	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

// UpdateDate обновляет дату задачи.
// Параметры:
//
//	date - новая дата.
//	id - идентификатор задачи.
//
// Возвращаемые значения:
//
//	error - ошибка, которая могла возникнуть в ходе работы.
func UpdateDate(date string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := db.Exec(query, sql.Named("date", date), sql.Named("id", id))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}
