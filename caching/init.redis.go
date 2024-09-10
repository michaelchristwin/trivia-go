package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Question struct {
	ID       string   `json:"id"`
	Text     string   `json:"text"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
	Category string   `json:"category"`
}

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func CacheQuestion(question Question) error {
	jsonData, err := json.Marshal(question)
	if err != nil {
		return err
	}

	err = redisClient.Set(ctx, fmt.Sprintf("question:%s", question.ID), jsonData, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetQuestionFromCache(id string) (Question, error) {
	var question Question

	jsonData, err := redisClient.Get(ctx, fmt.Sprintf("question:%s", id)).Result()
	if err != nil {
		return question, err
	}

	err = json.Unmarshal([]byte(jsonData), &question)
	if err != nil {
		return question, err
	}

	return question, nil
}
