package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegistrationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	regexEmail := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

	if !regexEmail.MatchString(req.Email) {
		http.Error(w, "Неверный E-Mail", http.StatusBadRequest)
		return
	}

	req.Password = strings.TrimSpace(req.Password)
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))

	regexUsername := regexp.MustCompile(`^[a-z0-9_]+$`)

	if !regexUsername.MatchString(req.Username) {
		http.Error(w, "Юзернейм состоит из недопустимых символов", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 4 || len(req.Username) > 32 {
		http.Error(w, "Юзернейм должен быть от 4 до 32 символов", http.StatusBadRequest)
		return
	}

	err := RegUserInDB(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := Response{Result: "Пользователь успешно зарегистрирован"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}

	req.Password = strings.TrimSpace(req.Password)
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))

	regexUsername := regexp.MustCompile(`^[a-z0-9_]+$`)

	if !regexUsername.MatchString(req.Username) {
		http.Error(w, "Юзернейм состоит из недопустимых символов", http.StatusBadRequest)
		return
	}

	if len(req.Username) < 4 || len(req.Username) > 32 {
		http.Error(w, "Юзернейм должен быть от 4 до 32 символов", http.StatusBadRequest)
		return
	}

	err := CheckUserInDB(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respStr := fmt.Sprintf("Вы зашли в аккаунт %v", req.Username)
	resp := Response{Result: respStr}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
