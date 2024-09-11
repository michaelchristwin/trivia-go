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

func CacheQuestions(questions []Question) error {
	for _, question := range questions {
		// Marshal each question to JSON
		jsonData, err := json.Marshal(question)
		if err != nil {
			return err
		}

		// Cache the question in Redis
		err = redisClient.Set(ctx, fmt.Sprintf("question:%s", question.ID), jsonData, 24*time.Hour).Err()
		if err != nil {
			return err
		}

		// Add the question ID to the Redis Set for the category
		err = redisClient.SAdd(ctx, fmt.Sprintf("questions:%s", question.Category), question.ID).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func getAndPopRandomQuestion(category string) (*Question, error) {
	// Construct the Redis key for the category's question set
	questionsKey := fmt.Sprintf("questions:%s", category)

	// Fetch a random question ID from the set
	randomQuestionID, err := redisClient.SRandMember(ctx, questionsKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error fetching random question ID: %v", err)
	}

	if randomQuestionID == "" {
		return nil, fmt.Errorf("no questions available")
	}

	// Remove the question ID from the set
	_, err = redisClient.SRem(ctx, questionsKey, randomQuestionID).Result()
	if err != nil {
		return nil, fmt.Errorf("error removing question ID from set: %v", err)
	}

	// Get the question data from Redis
	questionJSON, err := redisClient.Get(ctx, fmt.Sprintf("question:%s", randomQuestionID)).Result()
	if err != nil {
		return nil, fmt.Errorf("error fetching question data: %v", err)
	}

	var question Question
	if err := json.Unmarshal([]byte(questionJSON), &question); err != nil {
		return nil, fmt.Errorf("error unmarshalling question data: %v", err)
	}

	return &question, nil
}
