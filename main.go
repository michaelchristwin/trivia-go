package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/michaelchristwin/trivia-go/db"
	"github.com/michaelchristwin/trivia-go/handlers"
)

func main() {
	fmt.Println("This is Go baby")
	db.ConnectDB()
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.ListenAndServe(":8080", nil)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
