package main

import (
	config "cat/config"
	"cat/service"
)

func main() {
	config.EnvLoad()              // load .env
	cat_db := config.SQLiteLoad() // cat_sqlite.db
	service.BotLoad(cat_db)       // load telegram bot with api key from .env
}
