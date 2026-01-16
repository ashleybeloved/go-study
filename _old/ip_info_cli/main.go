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
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Printf("Usage: go run main.go 'IP' 'LANGUAGE (OPTIONAL)'\n")
		fmt.Printf("Example: go run main.go 24.48.0.1 ru\n\n")
		os.Exit(0)
	}

	query := os.Args[1]
	lang := "en"
	if len(os.Args) > 2 && len(os.Args[2]) > 1 {
		lang = os.Args[2]
	}

	apiUrl := fmt.Sprintf("http://ip-api.com/json/%v?lang=%v", query, lang)
	resp, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal("Ошибка в получении JSON: ", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Ошибка при чтении Response Body", err)
	}

	var info IpInfo
	err = json.Unmarshal(body, &info)
	if err != nil {
		log.Fatal("Ошибка парсинга JSON: ", err)
	}

	fmt.Printf("IP: %v\n\nCountry: %v\nRegion: %v\nGeo: %v, %v\nISP: %v", info.Query, info.Country, info.RegionName, info.Lat, info.Lon, info.Isp)

}

type IpInfo struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}
