package main

import (
	"io"
	"log"
	"net/http"
	"os"

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

	for update := range updates {
		if update.Message != nil {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			switch update.Message.Text {
			case "/cat":
				giveCat(bot, update)
				log.Printf("[%v] `Debug: Random cat sended.`", bot.Self.UserName)
				continue
			case "/gif":
				giveCatGif(bot, update)
				log.Printf("[%v] `Debug: Random cat gif sended.`", bot.Self.UserName)
				continue
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🐈 Random /cat image and /gif")
				bot.Send(msg)
				log.Printf("[%v] 🐈 Random /cat image and /gif", bot.Self.UserName)
				continue
			}
		}
	}
}

func giveCat(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	url := "https://cataas.com/cat"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка при скачивании:", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении данных:", err)
		return
	}

	photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileBytes{
		Name:  "cat.jpg",
		Bytes: data,
	})
	bot.Send(photo)
}

func giveCatGif(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	url := "https://cataas.com/cat/gif"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Ошибка при скачивании:", err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Ошибка при чтении данных:", err)
		return
	}

	gif := tgbotapi.NewAnimation(update.Message.Chat.ID, tgbotapi.FileBytes{
		Name:  "cat.gif",
		Bytes: data,
	})
	bot.Send(gif)
}
