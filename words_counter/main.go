package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

func main() {
	// Читаем файл - переводим в строку

	bytes, err := os.ReadFile("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	text := string(bytes)

	// Можно было бы сделать и так, но есть крутой пакет RegExp, который чистит всё одной строкой кода.

	// cleanText := strings.ReplaceAll(text, ",", "")
	// cleanText = strings.ReplaceAll(cleanText, "!", "")
	// cleanText = strings.ReplaceAll(cleanText, ".", "")
	// cleanText = strings.ReplaceAll(cleanText, ";", "")

	re := regexp.MustCompile(`[^\p{L}\s]+`) // [^\p{L}\s] - удалить всё, что НЕ является буквой любого алфавита и НЕ является пробелом
	cleanText := re.ReplaceAllString(text, "")

	words := strings.Split(cleanText, " ")
	symbols := utf8.RuneCountInString(text)

	commonWords := make(map[string]int)
	for _, word := range words {
		if word != "" {
			commonWords[word]++
		}
	}

	fmt.Printf("Stats:\n\n- Symbols: %v\n- Words: %v\n\nThe most common words:\n", symbols, len(words))

	type wordStat struct {
		word  string
		count int
	}

	var stats []wordStat
	for k, v := range commonWords {
		stats = append(stats, wordStat{k, v})
	}

	// Сортируем слайс от большего к меньшему

	sort.Slice(stats, func(i, j int) bool {
		return stats[i].count > stats[j].count
	})

	// Вывод топ-5 слов

	for i := 0; i < len(stats) && i < 5; i++ {
		fmt.Printf("%d. %s: %d\n", i+1, stats[i].word, stats[i].count)
	}
}
