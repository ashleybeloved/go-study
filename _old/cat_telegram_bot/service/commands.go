package service

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CommandHandlers(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	switch update.Message.Text {
	case "/cat":
		CatImage(update, bot)
		log.Printf("[%v] `Debug: Random cat-image sended.`", bot.Self.UserName)
	case "/gif":
		CatGif(update, bot)
		log.Printf("[%v] `Debug: Random cat-gif sended.`", bot.Self.UserName)
	case "/fact":
		CatFact(update, bot)
		log.Printf("[%v] `Debug: Random cat-fact sended.`", bot.Self.UserName)
	default:
		UnknownCommand(update, bot)
	}
}
