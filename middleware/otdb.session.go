package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type TokenManager struct {
    token           string
    lastRefresh     time.Time
    mutex           sync.Mutex
    refreshInterval time.Duration
}

func NewTokenManager() *TokenManager {
    return &TokenManager{
        refreshInterval: 6 * time.Hour, // Refresh token every 6 hours
    }
}

func (tm *TokenManager) GetToken() (string, error) {
    tm.mutex.Lock()
    defer tm.mutex.Unlock()

    if tm.token == "" || time.Since(tm.lastRefresh) > tm.refreshInterval {
        if err := tm.refreshToken(); err != nil {
            return "", err
        }
    }

    return tm.token, nil
}

func (tm *TokenManager) refreshToken() error {
    resp, err := http.Get("https://opentdb.com/api_token.php?command=request")
    if err != nil {
        return fmt.Errorf("error fetching token: %v", err)
    }
    defer resp.Body.Close()

    var result struct {
        ResponseCode int    `json:"response_code"`
        Token        string `json:"token"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return fmt.Errorf("error decoding token response: %v", err)
    }

    if result.ResponseCode != 0 {
        return fmt.Errorf("error response from API: %d", result.ResponseCode)
    }

    tm.token = result.Token
    tm.lastRefresh = time.Now()
    return nil
}
