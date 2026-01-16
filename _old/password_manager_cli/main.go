package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

// map - хэш-таблица для хранения паролей (name:pass)
var passwordDatabase = map[string]string{}

func main() {
	loadPasswords()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("[1] 🔑 Создать новый пароль")
		fmt.Println("[2] ❌ Удалить пароль")
		fmt.Println("[3] 📜 Посмотреть список паролей")
		fmt.Println("[4] 🚪 Выход")
		fmt.Print("\npassword_manager_cli | Выберите опцию: ")

		choose := readInt(reader)

		switch choose {
		case 1:
			passwordCreate(reader)
		case 2:
			passwordDel(reader)
		case 3:
			passwordList()
		case 4:
			bye()
		default:
			fmt.Print("\n❌ Неверный выбор. Попробуйте снова.\n\n")
		}
	}
}

func passwordCreate(reader *bufio.Reader) {
	fmt.Println("\n[1] 🛠 Создать пароль вручную")
	fmt.Println("[2] 🎲 Сгенерировать пароль")
	fmt.Println("[3] 🔙 Назад")
	fmt.Print("\npassword_manager_cli | Выберите опцию: ")

	choose := readInt(reader)

	switch choose {
	case 1:
		fmt.Print("\n✏ Введите пароль: ")
		passwordNew := readLine(reader)

		fmt.Print("🏷️ Введите его название (для чего он): ")
		passwordName := readLine(reader)

		passwordDatabase[passwordName] = passwordNew
		fmt.Println("\n✅ Пароль сохранён!")

		savePasswords()
	case 2:
		fmt.Print("\n📏 Введите длину пароля: ")
		length := readInt(reader)
		if length <= 0 {
			fmt.Println("⚠ Введите корректное положительное число.")
			return
		}

		passwordNew := generatePassword(length)

		fmt.Printf("\n🎯 Вот ваш пароль, сохраните его: %v\n", passwordNew)

		fmt.Print("🏷️ Введите его название (для чего он): ")
		passwordName := readLine(reader)

		passwordDatabase[passwordName] = passwordNew
		fmt.Println("\n✅ Пароль сохранён!")

		savePasswords()
	case 3:
		return
	default:
		fmt.Println("❌ Неверный выбор. Попробуйте снова.")
	}
}

func passwordList() {
	fmt.Print("\n📜 Вот список паролей: \n\n")

	if len(passwordDatabase) == 0 {
		fmt.Print("📂 Ваш список паролей пуст.\n\n")
		return
	}

	for passwordName, password := range passwordDatabase { // перебор по key:value
		fmt.Printf("🔹 %v: %v\n", passwordName, password)
	}
	fmt.Println()
}

func generatePassword(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?"
	password := make([]byte, length) // создание СРЕЗА байтов длиной length (1 символ - 1 байт)
	for i := range password {
		password[i] = chars[rand.Intn(len(chars))] // перебор по индексам среза
	}
	return string(password)
}

func bye() {
	fmt.Println("\n👋 Выход...")
	os.Exit(0)
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func readInt(reader *bufio.Reader) int {
	line := readLine(reader)
	var num int
	fmt.Sscanf(line, "%d", &num)
	return num
}

const dataFile = "passwords.json"

func loadPasswords() {
	file, err := os.Open(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			passwordDatabase = map[string]string{}
			return
		}
		fmt.Println("Ошибка чтения файла:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&passwordDatabase)
	if err != nil {
		fmt.Println("Ошибка парсинга JSON:", err)
	}
}

func savePasswords() {
	file, err := os.Create(dataFile)
	if err != nil {
		fmt.Println("Ошибка сохранения файла:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(passwordDatabase)
	if err != nil {
		fmt.Println("Ошибка кодирования JSON:", err)
	}
}

func passwordDel(reader *bufio.Reader) {
	if len(passwordDatabase) == 0 {
		fmt.Println("📂 Ваш список паролей пуст.")
		return
	}

	fmt.Println("\n📜 Вот список паролей:")
	for name := range passwordDatabase {
		fmt.Printf("🔹 %v\n", name)
	}

	fmt.Print("\n❓ Введите название пароля для удаления: ")
	name := readLine(reader)

	if _, exists := passwordDatabase[name]; exists {
		delete(passwordDatabase, name)
		fmt.Println("✅ Пароль успешно удалён!")
	} else {
		fmt.Println("⚠ Пароль с таким названием не найден.")
	}

	savePasswords()
}
