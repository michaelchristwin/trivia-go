package caching

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func fetchAndCacheQuestions(amount int, category string) error {
	// Construct URL for Open Trivia DB API
	url := fmt.Sprintf("https://opentdb.com/api.php?amount=%d&category=%s&type=multiple", amount, category)

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

	// Cache each question
	for i, q := range result.Results {
		question := Question{
			ID:       fmt.Sprintf("%s_%d", category, i),
			Text:     q.Question,
			Options:  append(q.Incorrect, q.Correct), // Note: You might want to shuffle these
			Answer:   q.Correct,
			Category: category,
		}

		if err := CacheQuestion(question); err != nil {
			fmt.Printf("Error caching question %s: %v\n", question.ID, err)
			// Decide whether to continue or return based on your error handling strategy
		}
	}

	return nil
}
