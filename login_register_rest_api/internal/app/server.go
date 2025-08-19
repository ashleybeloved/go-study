package app

import (
	"log"
	"net/http"

	"rest_api/internal/db"
	"rest_api/internal/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Run() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db.Connect()

	r.Post("/register", user.Register)
	r.Post("/login", user.Login)

	log.Println("Сервер запущен на localhost:1337")
	http.ListenAndServe(":1337", r)
}
