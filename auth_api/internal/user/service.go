package user

import (
	"fmt"
	"regexp"
	"strings"

	"rest_api/pkg/hash"
)

func RegisterUser(req RegistrationRequest) error {
	req.Email = strings.TrimSpace(req.Email)
	regexEmail := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

	if !regexEmail.MatchString(req.Email) {
		return fmt.Errorf("неверный E-Mail")
	}

	req.Password = strings.TrimSpace(req.Password)

	hashPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("ошибка при хэшировании пароля")
	}

	req.Password = hashPassword

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))

	regexUsername := regexp.MustCompile(`^[a-z0-9_]+$`)

	if !regexUsername.MatchString(req.Username) {
		return fmt.Errorf("юзернейм состоит из недопустимых символов")
	}

	if len(req.Username) < 4 || len(req.Username) > 32 {
		return fmt.Errorf("юзернейм должен быть от 4 до 32 символов")
	}

	return RegUserInDB(req)
}

func LoginUser(req LoginRequest) error {
	req.Password = strings.TrimSpace(req.Password)

	req.Username = strings.ToLower(strings.TrimSpace(req.Username))

	regexUsername := regexp.MustCompile(`^[a-z0-9_]+$`)

	if !regexUsername.MatchString(req.Username) {
		return fmt.Errorf("юзернейм состоит из недопустимых символов")
	}

	if len(req.Username) < 4 || len(req.Username) > 32 {
		return fmt.Errorf("юзернейм должен быть от 4 до 32 символов")
	}

	return CheckUserInDB(req)
}
