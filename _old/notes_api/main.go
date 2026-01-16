package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

var notes []Note
var nextID = 1

func main() {
	http.HandleFunc("/notes", notesHandler)
	http.HandleFunc("/notes/", idHandler)

	fmt.Println("Сервер на localhost:8080 запущен.")
	http.ListenAndServe(":8080", nil)
}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		viewNotes(w, r, notes)
	case http.MethodPost:
		addNote(w, r)
	default:
		http.Error(w, "This method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func idHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		delNote(w, id)
	case http.MethodGet:
		idStr := strings.TrimPrefix(r.URL.Path, "/notes/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		getNote(w, r, id)

	default:
		http.Error(w, "This method not allowed", http.StatusMethodNotAllowed)
	}

}

func viewNotes(w http.ResponseWriter, r *http.Request, notes []Note) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func addNote(w http.ResponseWriter, r *http.Request) {
	var newNote Note
	if err := json.NewDecoder(r.Body).Decode(&newNote); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	newNote.ID = nextID
	nextID++

	newNote.CreatedAt = time.Now()

	notes = append(notes, newNote)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newNote)
}

func delNote(w http.ResponseWriter, id int) {
	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

func getNote(w http.ResponseWriter, r *http.Request, id int) {
	for _, note := range notes {
		if note.ID == id {

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(note)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}
