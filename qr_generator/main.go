package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skip2/go-qrcode"
)

type genRequest struct {
	Content string `json:"content"`
}

func main() {
	http.HandleFunc("/gen", genHandle)

	fmt.Println("Сервер на localhost:8080 запущен.")
	http.ListenAndServe(":8080", nil)
}

func genHandle(w http.ResponseWriter, r *http.Request) {
	var req genRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if len(req.Content) == 0 {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	generateQR(w, r, req.Content)
}

func generateQR(w http.ResponseWriter, r *http.Request, content string) {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	w.Write(png)
}
