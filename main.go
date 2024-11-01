package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

// Entry represents a simple data structure with ID and Name.
type Entry struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Global variables to store entries and handle concurrent access.
var (
	entries = []Entry{
		{ID: 1, Name: "Book 1: One piece marineford arc"},
		{ID: 2, Name: "Book 2: Dragon ball super"},
		{ID: 3, Name: "Book 3: Bleach thousand year blood war"},
	}
	// Mutex to synchronize access to the entries slice.
	mutex = &sync.Mutex{}
)

// GetEntries handles GET requests to retrieve the list of entries.
func GetEntries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	mutex.Lock()
	defer mutex.Unlock()
	json.NewEncoder(w).Encode(entries)
}

// CreateEntry handles POST requests to create a new entry.
func CreateEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newEntry Entry
	if err := json.NewDecoder(r.Body).Decode(&newEntry); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the new entry has a unique ID
	mutex.Lock()
	newEntry.ID = len(entries) + 1
	entries = append(entries, newEntry)
	mutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEntry)
}

func main() {
	// Set up HTTP handlers for the API.
	http.HandleFunc("/entries", GetEntries)
	http.HandleFunc("/create", CreateEntry)

	// Start the server.
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
