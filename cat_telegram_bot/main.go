package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, переменные окружения могут быть не установлены")
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = false

	log.Printf("Авторизован на: [%v]", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for len(updates) > 0 {
		<-updates
	}

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			go handleUpdate(bot, update)
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch update.Message.Text {
	case "/cat":
		giveCat(bot, update)
		log.Printf("[%v] `Debug: Random cat sended.`", bot.Self.UserName)
	case "/gif":
		giveCatGif(bot, update)
		log.Printf("[%v] `Debug: Random cat gif sended.`", bot.Self.UserName)
	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🐈 Random /cat image and /gif")
		bot.Send(msg)
		log.Printf("[%v] 🐈 Random /cat image and /gif", bot.Self.UserName)
	}
}

func giveCatGif(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	url := fmt.Sprintf("https://cataas.com/cat/gif?%d", time.Now().UnixNano()) // Добавил рандом параметр, чтобы телеграм не кэшировал страницу и не отправлял один и тот же файл.

	gif := tgbotapi.NewAnimation(update.Message.Chat.ID, tgbotapi.FileURL(url))
	bot.Send(gif)
}

func giveCat(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	url := fmt.Sprintf("https://cataas.com/cat?%d", time.Now().UnixNano())

	photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileURL(url))
	bot.Send(photo)
}
