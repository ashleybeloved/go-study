package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect() (*sql.DB, error) {
	var err error
	DB, err = sql.Open("sqlite", "users.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT,
        email TEXT UNIQUE,
		registerAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		DB.Close()
		return nil, err
	}

	log.Println("Подключение к users.db успешно")
	return DB, nil
}
