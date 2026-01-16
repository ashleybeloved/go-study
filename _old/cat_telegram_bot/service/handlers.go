package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CatImage(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	url := fmt.Sprintf("https://cataas.com/cat?%v", time.Now().UnixNano())
	catphoto := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileURL(url))
	bot.Send(catphoto)
}

func CatGif(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	url := fmt.Sprintf("https://cataas.com/cat/gif?%v", time.Now().UnixNano())
	catanimation := tgbotapi.NewAnimation(update.Message.Chat.ID, tgbotapi.FileURL(url))
	bot.Send(catanimation)
}

func CatFact(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	url := "https://meowfacts.herokuapp.com/?lang=rus"

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка при GET-запросе: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении Response Body: %v", err)
	}

	var fact FactResponse
	err = json.Unmarshal(body, &fact)
	if err != nil {
		log.Printf("Ошибка парсинга JSON: %v\nОтвет сервера: %v", err, string(body))
	}

	if len(fact.Data) == 0 {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "😿 Не удалось получить факт о котах. Попробуй позже."))
		return
	}

	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, fact.Data[0]))
}

func UnknownCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "🐈 Я отправляю /cat, /gif и /fact"))
}

type FactResponse struct {
	Data []string `json:"data"`
}
