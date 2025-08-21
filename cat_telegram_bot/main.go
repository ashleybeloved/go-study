package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Не удалось загрузить .env файл, переменные окружения могут быть не установлены")
	}

	db, err := sqliteLoad()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
		if update.Message != nil && update.Message.Text != "" {
			if update.Message.From != nil {
				err = sqliteSaveUser(db, update.Message.From)
				if err != nil {
					log.Printf("Ошибка сохранения юзера: %v", err)
				}
			}

			if update.Message.Chat != nil {
				err = sqliteSaveChat(db, update.Message.Chat)
				if err != nil {
					log.Printf("Ошибка сохранения чата: %v", err)
				}
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			go handleUpdate(bot, update, db)
		}

		if update.Message != nil && update.Message.NewChatMembers != nil {
			for _, user := range update.Message.NewChatMembers {
				sendWelcomeMessage(bot, update, user)
			}
		}
	}

}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	switch update.Message.Text {
	case "/cat":
		giveCat(bot, update)
		sqliteLogCommand(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)
		log.Printf("[%v] `Debug: Random cat sended.`", bot.Self.UserName)
	case "/gif":
		giveCatGif(bot, update)
		sqliteLogCommand(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)
		log.Printf("[%v] `Debug: Random cat gif sended.`", bot.Self.UserName)
	case "/fact":
		giveCatFact(bot, update)
		sqliteLogCommand(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)
		log.Printf("[%v] `Debug: Random cat fact sended.`", bot.Self.UserName)
	case "/help":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🐈 Мои команды: /cat, /gif, /fact, /help")
		bot.Send(msg)
		sqliteLogCommand(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)
		log.Printf("[%v] 🐈 Мои команды: /cat, /gif, /fact, /help", bot.Self.UserName)
	case "/cat@ashley_cats_bot":
		giveCat(bot, update)
		sqliteLogCommand(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)
		log.Printf("[%v] `Debug: Random cat sended.`", bot.Self.UserName)
	case "/gif@ashley_cats_bot":
		giveCatGif(bot, update)
		sqliteLogCommand(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)
		log.Printf("[%v] `Debug: Random cat gif sended.`", bot.Self.UserName)
	case "/fact@ashley_cats_bot":
		giveCatFact(bot, update)
		sqliteLogCommand(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)
		log.Printf("[%v] `Debug: Random cat fact sended.`", bot.Self.UserName)
	case "/help@ashley_cats_bot":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🐈 Мои команды: /cat, /gif, /fact, /help")
		bot.Send(msg)
		sqliteLogCommand(db, update.Message.From.ID, update.Message.Chat.ID, update.Message.Text)
		log.Printf("[%v] 🐈 Мои команды: /cat, /gif, /fact, /help", bot.Self.UserName)
	}
}

func giveCatGif(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	url := fmt.Sprintf("https://cataas.com/cat/gif?%d", time.Now().UnixNano())

	gif := tgbotapi.NewAnimation(update.Message.Chat.ID, tgbotapi.FileURL(url))
	bot.Send(gif)
}

func giveCat(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var photo tgbotapi.PhotoConfig

	if rand.Intn(1000) < 1 { // 0.1% шанс
		photo = tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath("dog.jpg"))
		photo.Caption = "you find @fuckcensor (dev) dog"
	} else {
		url := fmt.Sprintf("https://cataas.com/cat?%d", time.Now().UnixNano())
		photo = tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FileURL(url))
	}

	bot.Send(photo)
}

func giveCatFact(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
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

func sendWelcomeMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, User tgbotapi.User) {

	url := fmt.Sprintf("https://cataas.com/cat/gif?%d", time.Now().UnixNano())
	gif := tgbotapi.NewAnimation(update.Message.Chat.ID, tgbotapi.FileURL(url))

	name := User.UserName
	message := fmt.Sprintf("мяу, @%s! добро пожаловать в группу!", name)
	if User.UserName == "" {
		name = User.FirstName
		message = fmt.Sprintf("мяу, %s! добро пожаловать в группу!", name)
	}

	gif.Caption = message

	bot.Send(gif)

}

func sqliteLoad() (*sql.DB, error) {
	if _, err := os.Stat("/data"); os.IsNotExist(err) {
		err := os.Mkdir("/data", 0755)
		if err != nil {
			return nil, fmt.Errorf("не удалось создать папку data: %v", err)
		}
	}

	db, err := sql.Open("sqlite", "/data/cat_users.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY,
		username TEXT,
		first_name TEXT NOT NULL,
		last_name TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}

	createChatsTableSQL := `
	CREATE TABLE IF NOT EXISTS chats (
		chat_id INTEGER PRIMARY KEY,
		title TEXT,
		type TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createChatsTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}

	createLogsTableSQL := `CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		chat_id INTEGER,
		command TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createLogsTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}

	log.Println("Подключение к БД успешно (папка data).")
	return db, nil
}

func sqliteSaveUser(db *sql.DB, user *tgbotapi.User) error {
	query := `INSERT OR IGNORE INTO users(user_id, username, first_name, last_name) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, user.ID, user.UserName, user.FirstName, user.LastName)
	return err
}

func sqliteSaveChat(db *sql.DB, chat *tgbotapi.Chat) error {
	query := `INSERT OR IGNORE INTO chats(chat_id, title, type) VALUES (?, ?, ?)`
	_, err := db.Exec(query, chat.ID, chat.Title, chat.Type)
	return err
}

func sqliteLogCommand(db *sql.DB, userID int64, chatID int64, command string) {
	_, _ = db.Exec(
		"INSERT INTO logs(user_id, chat_id, command) VALUES (?, ?, ?)",
		userID, chatID, command,
	)
}

type FactResponse struct {
	Data []string `json:"data"`
}
