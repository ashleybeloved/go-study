package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

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
		if update.Message != nil {

			err := sqliteSaveUser(db, update.Message.From)
			if err != nil {
				panic(err)
			}

			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🌍 Я бот погоды для Telegram. Напиши свой город, чтобы узнать погоду.")
				bot.Send(msg)
				log.Printf("[%v] 🌍 Я бот погоды для Telegram. Напиши свой город, чтобы узнать погоду.", bot.Self.UserName)
				continue
			}

			if update.Message.Location != nil {
				lat := update.Message.Location.Latitude
				lon := update.Message.Location.Longitude

				weatherInfo, err := getWeatherByLocation(lat, lon)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: "+err.Error()))
					continue
				}

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, weatherInfo))
				log.Printf("[%v] %v", bot.Self.UserName, weatherInfo)
			} else {
				city := update.Message.Text
				weatherInfo, err := getWeather(city)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: "+err.Error()))
					continue
				}

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, weatherInfo))
				log.Printf("[%v] %v", bot.Self.UserName, weatherInfo)
			}

		}
	}
}

func getWeatherByLocation(lat float64, lon float64) (string, error) {
	API_KEY := os.Getenv("API_WEATHER")
	if API_KEY == "" {
		return "", fmt.Errorf("API ключ не найден в переменных окружения")
	}
	apiUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%v&q=%f,%f&aqi=no", API_KEY, lat, lon)

	weather, err := getJson(apiUrl)
	if err != nil {
		return "", err
	}

	weatherInfo := fmt.Sprintf("🌍 Погода сейчас в %v, %v\n\n🕖 Локальное время: %v\n🌡 Температура: %v°C\n🤔 Ощущается как: %v°C\n💨 Ветер: %v км/ч\n", weather.Location.Name, weather.Location.Country, weather.Location.Localtime, weather.Current.TempC, weather.Current.FeelslikeC, weather.Current.WindKPH)

	if weather.Location.Name == "" {
		weatherInfo = "Такого города не существует, попробуй снова."
	}

	return weatherInfo, err
}

func getWeather(city string) (string, error) {
	escapedCity := url.QueryEscape(city)

	API_KEY := os.Getenv("API_WEATHER")
	if API_KEY == "" {
		return "", fmt.Errorf("API ключ не найден в переменных окружения")
	}
	apiUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%v&q=%v&aqi=no", API_KEY, escapedCity)

	weather, err := getJson(apiUrl)
	if err != nil {
		return "", err
	}

	if weather.Location.Name == "" {
		return "Такого города не существует, попробуй снова.", err
	} else {
		weatherInfo := fmt.Sprintf("🌍 Погода сейчас в %v, %v\n\n🕖 Локальное время: %v\n🌡 Температура: %v°C\n🤔 Ощущается как: %v°C\n💨 Ветер: %v км/ч\n", weather.Location.Name, weather.Location.Country, weather.Location.Localtime, weather.Current.TempC, weather.Current.FeelslikeC, weather.Current.WindKPH)
		return weatherInfo, err
	}

}

func getJson(apiUrl string) (WeatherResponse, error) {
	var weather WeatherResponse

	resp, err := http.Get(apiUrl)
	if err != nil {
		return weather, fmt.Errorf("ошибка при GET-запросе: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weather, fmt.Errorf("ошибка при GET-запросе: %v", err)
	}

	err = json.Unmarshal(body, &weather)
	if err != nil {
		err = fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	return weather, err
}

func sqliteLoad() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "telegram_bot.db")
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

	log.Println("Подключение к БД успешно.")
	return db, nil
}

func sqliteSaveUser(db *sql.DB, user *tgbotapi.User) error {
	query := `INSERT OR IGNORE INTO users(user_id, username, first_name, last_name) VALUES (?, ?, ?, ?)`

	_, err := db.Exec(query, user.ID, user.UserName, user.FirstName, user.LastName)
	return err
}

type WeatherResponse struct {
	Location struct {
		Name      string `json:"name"`
		Country   string `json:"country"`
		Localtime string `json:"localtime"`
	} `json:"location"`
	Current struct {
		TempC      float64 `json:"temp_c"`
		FeelslikeC float64 `json:"feelslike_c"`
		WindKPH    float64 `json:"wind_kph"`
	} `json:"current"`
}
