package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			text := update.Message.Text

			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🌍 Я бот погоды для Telegram. Напиши свой город, чтобы узнать погоду.")
				bot.Send(msg)
				continue
			}

			if update.Message.Location != nil && text == "" {
				lat := update.Message.Location.Latitude
				lon := update.Message.Location.Longitude

				weatherInfo, err := getWeatherByLocation(lat, lon)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: "+err.Error()))
					continue
				}

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, weatherInfo))
			} else {
				city := update.Message.Text
				weatherInfo, err := getWeather(city, text)

				if err != nil {
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка: "+err.Error()))
					continue
				}

				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, weatherInfo))
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

	weatherInfo := fmt.Sprintf("🌍 Погода сейчас в %v, %v\n\n🌡 Температура: %v°C\n🤔 Ощущается как: %v°C\n💨 Ветер: %v км/ч\n", weather.Location.Name, weather.Location.Country, weather.Current.TempC, weather.Current.FeelslikeC, weather.Current.WindKPH)

	if weather.Location.Name == "" {
		weatherInfo = "Такого города не существует, попробуй снова."
	}

	return weatherInfo, err
}

func getWeather(city string, text string) (string, error) {
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

	if weather.Location.Name == "" && text != "" {
		return "Такого города не существует, попробуй снова.", err
	} else {
		weatherInfo := fmt.Sprintf("🌍 Погода сейчас в %v, %v\n\n🌡 Температура: %v°C\n🤔 Ощущается как: %v°C\n💨 Ветер: %v км/ч\n", weather.Location.Name, weather.Location.Country, weather.Current.TempC, weather.Current.FeelslikeC, weather.Current.WindKPH)
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

type WeatherResponse struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC      float64 `json:"temp_c"`
		FeelslikeC float64 `json:"feelslike_c"`
		WindKPH    float64 `json:"wind_kph"`
	} `json:"current"`
}
