package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/michaelchristwin/trivia-go/db"
	"github.com/michaelchristwin/trivia-go/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserFactory struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	ID        string `bson:"_id"`
	FirstName string `bson:"firstName"`
	LastName  string `bson:"lastName"`
	Email     string `bson:"email"`
	Password  string `bson:"password"`
}

var sessionStore = map[string]string{}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello Go world")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Cookie exists, validate session
		sessionID := cookie.Value
		userID, sessionExists := sessionStore[sessionID]
		if sessionExists {
			// Session is valid, user is already logged in
			fmt.Fprintf(w, "Welcome back, User %s!", userID)
			return
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	time.Sleep(time.Millisecond * 300)

	var loginReq UserFactory
	if err := json.Unmarshal(body, &loginReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	filter := bson.D{{Key: "email", Value: loginReq.Email}}
	collection := db.MongoClient.Database("users").Collection("users")
	var result User
	err = collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		} else {
			log.Printf("Error finding document: %v", err)
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		}
		return
	}

	if err := utils.CheckPasswordHash(loginReq.Password, result.Password); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionStore[sessionID] = result.ID
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		MaxAge:   24 * 60 * 60, // 24 hours
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	data := struct {
		ID        string
		FirstName string
		LastName  string
		Email     string
	}{
		ID:        result.ID,
		FirstName: result.FirstName,
		LastName:  result.LastName,
		Email:     result.Email,
	}
	response := map[string]interface{}{
		"message": "Login successful",
		"data":    data, // Pass `data` as an object, not as a JSON string
	}
	responseBody, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var siginReq UserFactory
	if err := json.Unmarshal(body, &siginReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	hashedPassword, err := utils.HashPassword(siginReq.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	document := bson.D{{Key: "email", Value: siginReq.Email}, {Key: "password", Value: hashedPassword}, {Key: "firstName", Value: siginReq.FirstName}, {Key: "lastName", Value: siginReq.LastName}}
	collection := db.MongoClient.Database("users").Collection("users")
	result, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	fmt.Printf("Inserted document with ID: %v\n", result.InsertedID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User added successfully"))
}
