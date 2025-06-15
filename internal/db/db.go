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
	defer db.Close()

	if install {
		_, err := db.Exec(createSchedulerTable)
		if err != nil {
			return err
		}
	}
	return nil
}
