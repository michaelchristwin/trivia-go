package main

import (
	"fmt"

	"github.com/michaelchristwin/trivia-go/db"
)

func main() {
	fmt.Println("This is Go baby")
	db.ConnectDB()
}
