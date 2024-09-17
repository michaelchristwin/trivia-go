package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/michaelchristwin/trivia-go/db"
	"github.com/michaelchristwin/trivia-go/handlers"
	"github.com/michaelchristwin/trivia-go/middleware"
)

func main() {
	fmt.Println("This is Go baby")

	if err := db.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handlers.HomeHandler))
	mux.Handle("/login", http.HandlerFunc(handlers.LoginHandler))
	mux.Handle("/signup", http.HandlerFunc(handlers.AddUser))
	wrappedMux := middleware.CORSmiddleware(mux)
	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", wrappedMux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
