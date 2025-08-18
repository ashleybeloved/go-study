package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Result string `json:"result"`
}

func main() {
	db, err := sqliteLoad()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		req.Password = strings.TrimSpace(req.Password)
		req.Username = strings.ToLower(strings.TrimSpace(req.Username))

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Ошибка при хэшировании пароля", http.StatusInternalServerError)
			return
		}

		query := `INSERT OR IGNORE INTO users(username, password) VALUES (?, ?)`

		res, err := db.Exec(query, req.Username, string(hashedPassword))
		if err != nil {
			http.Error(w, "Ошибка записи в БД", http.StatusInternalServerError)
			return
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			http.Error(w, "Такой пользователь уже существует", http.StatusConflict)
			return
		}

		resp := Response{Result: "Пользователь зарегистрирован"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var req Request

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный JSON", http.StatusBadRequest)
			return
		}

		req.Password = strings.TrimSpace(req.Password)
		req.Username = strings.ToLower(strings.TrimSpace(req.Username))

		var storedPassword string
		query := `SELECT password FROM users WHERE username = ?`
		err := db.QueryRow(query, req.Username).Scan(&storedPassword)
		if err == sql.ErrNoRows {
			http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
		if err != nil {
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
		}

		respStr := fmt.Sprintf("Вы зашли в аккаунт %v", req.Username)
		resp := Response{Result: respStr}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	log.Println("Сервер запущен на :7777")
	http.ListenAndServe(":7777", r)
}

func sqliteLoad() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "users.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT,
		registerAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}

	log.Println("Подключение к БД успешно.")
	return db, nil
}
