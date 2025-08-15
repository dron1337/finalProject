package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title varchar(100),
    repeat varchar(128),
    comment text
);
CREATE INDEX idx_date ON scheduler(date);`

func Init() error {
	var err error
	_, err = os.Stat(os.Getenv("TODO_DBFILE"))
	install := os.IsNotExist(err)

	DB, err = sql.Open("sqlite", os.Getenv("TODO_DBFILE"))
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	if install {
		if _, err = DB.Exec(schema); err != nil {
			return fmt.Errorf("failed to create database schema: %v", err)
		}
	}

	return nil
}
