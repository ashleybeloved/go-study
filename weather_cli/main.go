package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠ Не удалось загрузить .env файл", err)
	}

	fmt.Print("weather_cli | Чтобы узнать погоду, напишите ваш город: ")

	reader := bufio.NewReader(os.Stdin)
	city, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("⚠ Ошибка при чтении ввода:", err)
	}
	escapedCity := url.QueryEscape(city)

	API_KEY := os.Getenv("API_WEATHER")
	apiUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", API_KEY, escapedCity)

	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal("⚠ Ошибка при GET-запросе:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("⚠ Ошибка при чтении Response Body", err)
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		log.Fatalf("⚠ Ошибка парсинга JSON: %v\nОтвет сервера: %v", err, string(body))
	}

	fmt.Printf("🌍 Погода сейчас в %v, %v\n\n", weather.Location.Name, weather.Location.Country)
	fmt.Printf("🌡 Температура: %v°C\n", weather.Current.TempC)
	fmt.Printf("🤔 Ощущается как: %v°C\n", weather.Current.FeelslikeC)
	fmt.Printf("💨 Ветер: %v км/ч\n", weather.Current.WindKPH)
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
