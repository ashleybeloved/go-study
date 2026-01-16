package service

import (
	"database/sql"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BotLoad(cat_db *sql.DB) {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	log.Printf("Бот авторизован на: [%v]\n", bot.Self.UserName)

	for update := range updates {
		if update.Message != nil && update.Message.Text != "" {
			SaveUserSQLite(cat_db, update.Message.From)

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			go CommandHandlers(update, bot)
		}
	}
}
