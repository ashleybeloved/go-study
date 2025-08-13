package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) >= 5 || len(os.Args) < 4 {
		fmt.Printf("Usage: go run main.go 'FROM' 'TO' 'AMOUNT'\n")
		fmt.Printf("Example: go run main.go usd eur 100\n\n")

		os.Exit(0)
	}

	from := os.Args[1]
	to := os.Args[2]
	amount := os.Args[3]

	apiUrl := fmt.Sprintf("https://api.frankfurter.app/latest?amount=%v&from=%v&to=%v", amount, from, to) // Это API к сожалению не поддерживает рубли, т.к. берёт информацию из https://www.ecb.europa.eu/stats/policy_and_exchange_rates/euro_reference_exchange_rates/html/index.en.html
	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal("Ошибка в получении JSON: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка при чтении Response Body", err)
	}

	var currency currencyResponse
	err = json.Unmarshal(body, &currency)
	if err != nil {
		log.Fatal("Ошибка парсинга JSON: ", err)
	}

	for key, value := range currency.Rates {
		fmt.Printf("%v %v = %v % v | на %v", currency.Amount, currency.Base, value, key, currency.Date)
	}
}

type currencyResponse struct {
	Amount float64            `json:"amount"`
	Base   string             `json:"base"`
	Date   string             `json:"date"`
	Rates  map[string]float64 `json:"rates"`
}
