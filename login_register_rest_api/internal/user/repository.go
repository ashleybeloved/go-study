package user

import (
	"database/sql"
	"fmt"
	"rest_api/internal/db"
	"rest_api/pkg/hash"

	_ "modernc.org/sqlite"
)

func RegUserInDB(req RegistrationRequest) error {
	query := `INSERT OR IGNORE INTO users(username, password, email) VALUES (?, ?, ?)`

	res, err := db.DB.Exec(query, req.Username, req.Password, req.Email)
	if err != nil {
		return fmt.Errorf("ошибка записи в БД")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("такой пользователь уже существует")
	}

	return nil
}

func CheckUserInDB(req LoginRequest) error {
	query := `SELECT password FROM users WHERE username = ?`

	var storedHash string
	err := db.DB.QueryRow(query, req.Username).Scan(&storedHash)
	if err == sql.ErrNoRows {
		return fmt.Errorf("пользователь не найден")
	} else if err != nil {
		return fmt.Errorf("ошибка запроса к БД")
	}

	err = hash.CheckPassword(storedHash, req.Password)
	if err != nil {
		return fmt.Errorf("неверный пароль")
	}

	return nil
}
