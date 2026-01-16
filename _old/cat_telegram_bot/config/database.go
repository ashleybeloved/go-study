package config

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func SQLiteLoad() (cat_db *sql.DB) {
	_, err := os.Stat("/cat/data")
	if os.IsNotExist(err) {
		err = os.Mkdir("/cat/data", 0775)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Папка data создана")
	} else {
		log.Println("Папка data уже создана")
	}

	log.Println("Подключаюсь к базе данных")

	cat_db, err = sql.Open("sqlite", "/cat/data/cat_sqlite.db")
	if err != nil {
		log.Fatal(err)
	}

	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY,
		username TEXT,
		first_name TEXT NOT NULL,
		last_name TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = cat_db.Exec(usersTable)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Подключение к базе данных произошло успешно")

	return cat_db
}
