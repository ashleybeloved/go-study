package main

import (
	"log"
	"net/http"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/gin-gonic/gin"
)

type Paste struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type reqNewPaste struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Password string `json:"password"`
}

type reqFindPaste struct {
	Password string `json:"password"`
}

var pastes []Paste

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/", newPaste)
	r.POST("/:id", findPaste)

	log.Println("Server running on :8080")
	r.Run(":8080")
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

	id, err := gonanoid.New(10)
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

	pastes = append(pastes, paste)

	c.JSON(http.StatusAccepted, paste)
}

func findPaste(c *gin.Context) {
	var req reqFindPaste
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	id := c.Request.URL.Path[1:]

	for _, paste := range pastes {
		if paste.ID == id {
			if paste.Password == req.Password {
				c.JSON(http.StatusAccepted, paste)
				return
			} else {
				c.JSON(http.StatusForbidden, gin.H{"error": "Invalid Password"})
				return
			}
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Page Not Found"})
}
