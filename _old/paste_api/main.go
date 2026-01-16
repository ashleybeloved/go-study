package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Paste struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Password  string    `json:"password"`
	CreatedAt time.Time `jsgon:"createdAt"`
}

type reqNewPaste struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Password string `json:"password"`
}

type reqFindPaste struct {
	Password string `json:"password"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	if err := godotenv.Load(".config"); err != nil {
		log.Println("File .env not found")
	}

	if err := os.MkdirAll("pastes", 0755); err != nil {
		log.Print("Failed to create folder")
	}

	port := os.Getenv("PORT")

	r.POST("/", newPaste)
	r.POST("/:id", findPaste)

	log.Println("Server running on " + port)
	r.Run(port)
}

func newPaste(c *gin.Context) {
	var req reqNewPaste
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if len(req.Title) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title must be filled"})
		return
	}

	if len(req.Content) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content must be filled"})
		return
	}

	digitsStr := os.Getenv("DIGITS")
	digits, err := strconv.Atoi(digitsStr)
	if err != nil {
		log.Println("Failed to convert string to int, set default to 10")
		digits = 10
	}

	id, err := gonanoid.New(digits)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate ID"})
		return
	}

	paste := Paste{
		ID:        id,
		Title:     req.Title,
		Content:   req.Content,
		Password:  req.Password,
		CreatedAt: time.Now(),
	}

	filename := filepath.Join("pastes", id+".json")
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create file %s: %v", filename, err)
		c.JSON(http.StatusInternalServerError, "Failed to create paste")
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(paste); err != nil {
		log.Printf("Failed to write JSON to file %s: %v", filename, err)
		c.JSON(http.StatusInternalServerError, "Failed to encode paste")
		return
	}

	c.JSON(http.StatusAccepted, paste)
}

func findPaste(c *gin.Context) {
	var req reqFindPaste
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	id := c.Request.URL.Path[1:]

	filename := filepath.Join("pastes", id+".json")
	if _, err := os.Stat(filename); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page Not Found"})
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open file %s: %v", filename, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read paste"})
		return
	}
	defer file.Close()

	var paste Paste
	if err := json.NewDecoder(file).Decode(&paste); err != nil {
		log.Printf("Failed to decode JSON from %s: %v", filename, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Corrupted paste data"})
		return
	}

	if req.Password == paste.Password {
		c.JSON(http.StatusAccepted, paste)
		return
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid password"})
	}
}
