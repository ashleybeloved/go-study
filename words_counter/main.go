package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

func main() {
	bytes, err := os.ReadFile("./input.txt")
	if err != nil {
		log.Fatal(err)
	}

	text := string(bytes)

	re := regexp.MustCompile(`\P{L}+`)
	cleanText := re.ReplaceAllString(text, " ")
	cleanText = strings.ToLower(cleanText)

	words := strings.Fields(cleanText)

	commonWords := make(map[string]int)

	for _, word := range words {
		commonWords[word]++
	}

	keys := make([]string, 0, len(commonWords))
	for k := range commonWords {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return commonWords[keys[i]] > commonWords[keys[j]]
	})

	for _, k := range keys {
		fmt.Printf("%s: %d\n", k, commonWords[k])
	}
}
