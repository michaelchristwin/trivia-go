package handlers

import (
	"fmt"
	"io"
	"net/http"
)

var sessionStore = map[string]string{}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Go world")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	fmt.Printf("Received POST request with body: %s\n", body)
	// Send a response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("POST request received"))
}
