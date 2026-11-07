package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Confg struct {
	Port string
}

func Load() *Confg {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system vars")
	}
	return &Confg{Port: getEnv("PORT", "8080")}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
