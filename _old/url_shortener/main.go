package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	URL string `json:"url"`
	ID  int    `json:"ID"`
}

var URLs []ShortenResponse
var nextID int

func main() {
	http.HandleFunc("/shorten", shortenHandle)
	http.HandleFunc("/", redirectHandle)

	fmt.Println("Сервер запущен на localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func shortenHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "This method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newURL ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&newURL); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if strings.HasPrefix(newURL.URL, "https://") || strings.HasPrefix(newURL.URL, "http://") {
		generateShortURL(w, newURL.URL)
	} else {
		http.Error(w, "URL must start with https:// or http://", http.StatusBadRequest)
		return
	}

}

func generateShortURL(w http.ResponseWriter, url string) {
	newURL := ShortenResponse{
		URL: url,
		ID:  nextID,
	}

	nextID++

	URLs = append(URLs, newURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newURL)
}

func redirectHandle(w http.ResponseWriter, r *http.Request) {
	idURLstr := r.URL.Path[1:]
	idURL, err := strconv.Atoi(idURLstr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, url := range URLs {
		if url.ID == idURL {
			fmt.Fprintf(w, "Redirecting to: %s", url.URL)
			w.Header().Set("Content-Type", "redirect")
			http.Redirect(w, r, url.URL, http.StatusFound)
			return
		}
	}

	http.Error(w, "Invalid ID", http.StatusNotFound)
}
