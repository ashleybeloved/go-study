package service

import (
	"database/sql"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SaveUserSQLite(cat_db *sql.DB, user *tgbotapi.User) {
	query := `INSERT OR IGNORE INTO users(user_id, username, first_name, last_name) VALUES (?, ?, ?, ?)`

	_, err := cat_db.Exec(query, user.ID, user.UserName, user.FirstName, user.LastName)
	if err != nil {
		log.Println(err)
	}
}
