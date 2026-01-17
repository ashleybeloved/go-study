package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var sites []string

const path = "./sites.txt"

func main() {
	bytes, err := os.ReadFile(path)
	if err != nil {
		os.WriteFile(path, []byte(""), 0644)
	} else {
		if len(string(bytes)) > 0 {
			sites = strings.Fields(string(bytes))
		}
	}

	fmt.Print("Welcome to Site Checker\n\n")
	Start()
}

func Start() {
	fmt.Print("1. Run checker\n2. List of my sites\n3. Add new site\n4. Delete site\n0. Exit\n\nChoose your option: ")

	var choose string
	fmt.Scanln(&choose)

	switch choose {
	case "1":
		Run()
	case "2":
		List()
	case "3":
		NewSite()
	case "4":
		DelSite()
	case "0":
		log.Println("Bye!")
		os.Exit(1)
	default:
		log.Fatalln("Unknown command")
	}
}

func Run() {
	if len(sites) == 0 {
		fmt.Println("You don't have any sites to check")
		Start()
	}

	results := make(chan string, len(sites))

	for _, url := range sites {
		go func(u string) {
			client := http.Client{Timeout: 5 * time.Second}

			if !strings.HasPrefix(u, "http") {
				u = "http://" + u
			}

			resp, err := client.Get(u)
			if err != nil {
				results <- fmt.Sprintf("[DOWN] %s (Error: %v)", u, err)
				return
			}

			defer resp.Body.Close()

			results <- fmt.Sprintf("[%d] %s", resp.StatusCode, u)
		}(url)
	}

	for i := 0; i < len(sites); i++ {
		fmt.Println(<-results)
	}

	Start()
}

func List() {
	if len(sites) == 0 {
		fmt.Println("You don't have any sites")
		Start()
	}

	fmt.Println("Your sites:")

	for i, site := range sites {
		fmt.Printf("%d. %v\n", i+1, site)
	}

	fmt.Println()

	Start()
}

func NewSite() {
	var url string
	fmt.Print("URL: ")
	fmt.Scanln(&url)

	sites = append(sites, url)

	data := strings.Join(sites, "\n")
	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("You added '%v' to checker\n", url)

	Start()
}

func DelSite() {
	var idStr string
	fmt.Print("Write ID site to delete: ")
	fmt.Scanln(&idStr)

	id, _ := strconv.Atoi(idStr)
	id -= 1

	if id >= len(sites) || id < 0 {
		fmt.Println("Invalid ID")
		Start()
	}

	siteUrl := sites[id]

	sites = append(sites[:id], sites[id+1:]...)

	data := strings.Join(sites, "\n")
	err := os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nSite '%v' successfully deleted\n", siteUrl)
	Start()
}
