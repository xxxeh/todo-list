// Пакет db содержит функции для работы с базой данных SQLite.
package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const createSchedulerTable string = `CREATE TABLE scheduler (
									id INTEGER PRIMARY KEY AUTOINCREMENT,
									date char(8) NOT NULL DEFAULT "",
									title varchar NOT NULL DEFAULT "",
									comment TEXT NOT NULL DEFAULT "",
									repeat varchar(128) NOT NULL DEFAULT "");
									CREATE INDEX scheduler_date on scheduler (date);`

var db *sql.DB

// Init инициализирует подключение к базе данных SQLite и создаёт таблицу scheduler, если база данных ещё не существует.
//
// Параметры:
//
//	dbFile - путь до файла БД, должен быть указан в переменной окружения TODO_DBFILE
//
// Возвращаемые значение:
//
//	error - ошибка, если не удалось установить подключение или выполнить sql-запрос.
func Init(dbFile string) error {
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	//Подключение не закрываем, т.к. оно должно быть открыто постоянно, пока работает сервис.

	if install {
		_, err := db.Exec(createSchedulerTable)
		if err != nil {
			return err
		}
	}
	return nil
}
