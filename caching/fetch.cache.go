package caching

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/michaelchristwin/trivia-go/middleware"
)

var tokenManager = middleware.NewTokenManager()
var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// Shuffle the options in place
func shuffleOptions(options []string) {
	for i := len(options) - 1; i > 0; i-- {
		j := rnd.Intn(i + 1)
		options[i], options[j] = options[j], options[i]
	}
}

func fetchAndCacheQuestions(amount int, category string) error {
	// Construct URL for Open Trivia DB API
	token, err := tokenManager.GetToken()
	if err != nil {
		return fmt.Errorf("error getting token: %v", err)
	}
	url := fmt.Sprintf("https://opentdb.com/api.php?amount=%d&category=%s&&difficulty=easy&type=multiple&token=%s", amount, category, token)

	// Fetch questions from API
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching questions: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			Question  string   `json:"question"`
			Correct   string   `json:"correct_answer"`
			Incorrect []string `json:"incorrect_answers"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	// Prepare questions for caching
	var questions []Question
	for i, q := range result.Results {
		options := append(q.Incorrect, q.Correct) // Combine correct and incorrect answers
		shuffleOptions(options)                   // Shuffle options

		question := Question{
			ID:       fmt.Sprintf("%s_%d", category, i),
			Text:     q.Question,
			Options:  options,
			Answer:   q.Correct,
			Category: category,
		}
		questions = append(questions, question)
	}

	// Cache all questions at once
	if err := CacheQuestions(questions); err != nil {
		return fmt.Errorf("error caching questions: %v", err)
	}

	return nil
}
