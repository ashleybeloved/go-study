package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() {
	wordleText := color.New(color.FgGreen, color.Bold)

	fmt.Print("Добро пожаловать в ")
	wordleText.Print("Wordle-CLI!\n\n")

	fmt.Print("[1] Новая игра\n[2] Правила игры\n\n[0] Выход\n\n")
	fmt.Print("Выберите опцию: ")

	var choose string

	fmt.Scan(&choose)
	choose = strings.TrimSpace(choose)

	switch choose {
	case "1":
		RunGame()
	case "2":
		Help()
	case "0":
		Exit()
	default:
		errorInput()
	}
}

func RunGame() {
	wordsJson, err := os.Open("words.json")
	if err != nil {
		log.Fatal(err)
	}
	defer wordsJson.Close()

	var words []string
	err = json.NewDecoder(wordsJson).Decode(&words)
	if err != nil {
		log.Fatal(err)
	}

	var uword string
	try := 6

	randIndex := rand.Intn(len(words) - 1)
	word := words[randIndex]

	green := color.New(color.FgGreen, color.Bold)
	fmt.Print("\nЯ загадал ")
	green.Print("слово")
	fmt.Print(" из 5 букв! У тебя есть 6 попыток чтобы его угадать. Начинай!\n\n")

	for try >= 1 {
		fmt.Scan(&uword)

		if uword == word {
			Win(word)
		}

		if len(uword) != len(word) {
			deleteLine()
			fmt.Printf("%v \t| Слово должно состоять из 5 букв!\n", uword)
			continue
		}

		if !slices.Contains(words, uword) {
			deleteLine()
			fmt.Printf("%v \t| Такого слова не существует!\n", uword)
			continue
		}

		uword := strings.ToLower(uword)

		deleteLine()

		var wordslice []string
		var uwordslice []string

		for _, k := range word {
			wordslice = append(wordslice, string(k))
		}
		for _, k := range uword {
			uwordslice = append(uwordslice, string(k))
		}

		for i, k := range uwordslice {
			if wordslice[i] == k {
				green := color.New(color.BgGreen)
				green.Print(k)
			} else if slices.Contains(wordslice, k) {
				yellow := color.New(color.BgYellow)
				yellow.Print(k)
			} else {
				fmt.Print(k)
			}
		}
		try--
		fmt.Print("\t| Попыток осталось: ", try, "\n")
	}

	Lose(word)
}

func Win(word string) {
	green := color.New(color.FgGreen)
	green.Print("\n\nТы выиграл! Ты отгадал моё слово ")
	greenhi := color.New(color.FgHiGreen, color.Bold)
	greenhi.Print(word)
	greenhi.DisableColor()
	fmt.Print("\n\n")

	fmt.Print("Перенесу тебя в меню, через 10 секунд...\n\n")

	time.Sleep(10 * time.Second)
	main()
}

func Lose(word string) {
	red := color.New(color.FgRed)
	red.Print("\n\nТы проиграл! Загаданное слово: ")
	greenbg := color.New(color.FgGreen, color.Bold)
	greenbg.Print(word)
	greenbg.DisableColor()
	fmt.Print("\n\n")

	fmt.Print("Перенесу тебя в меню, через 10 секунд...\n\n")

	time.Sleep(10 * time.Second)
	main()
}

func deleteLine() {
	fmt.Print("\033[F\033[K")
}

func Help() {

	whitebold := color.New(color.FgHiWhite, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)

	fmt.Print("\n=-----------------------------------------=\n\n")

	whitebold.Print("Правила игры Wordle:\n\n")
	whitebold.Print("Цель игры: ")
	fmt.Print("Угадать загаданное слово, состоящее из 5 букв, за 6 попыток.\n")

	whitebold.Print("Начало игры: ")
	fmt.Print("Игрок вводит любое слово из 5 букв в качестве первой попытки.\n")

	whitebold.Print("Анализ: ")
	fmt.Print("- Если буква в слове на правильной позиции, она подсвечивается ")
	green.Print("зеленым.")
	fmt.Print("\n\t- Если буква есть в слове, но на неправильной позиции, она подсвечивается ")
	yellow.Print("жёлтым.")
	fmt.Print("\n\t- Если буква отсутствует в слове, она остаётся без подсветки.\n")

	whitebold.Print("Победа: ")
	fmt.Print("Игра заканчивается, когда игрок угадывает слово или использует все 6 попыток. В случае поражения игрок видит загаданное слово.\n")

	fmt.Print("\n=-----------------------------------------=\n\n")

	main()
}

func errorInput() {
	errorText := color.New(color.FgRed, color.Bold)
	errorText.Print("\nТакой опции не существует!\n\n")
	main()
}

func Exit() {
	fmt.Println("\nПока!")
	os.Exit(0)
}
