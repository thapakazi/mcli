package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"mcli/types"

	"github.com/joho/godotenv"
)

// TODO: return config in future
func LoadEnv() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", fmt.Errorf("Error loading .env file, %w", err)
	}
	apiBaseUrl := os.Getenv("API_BASE_URL")
	if apiBaseUrl == "" {
		return "", fmt.Errorf("API_BASE_URL not set in .env file")
	}
	return apiBaseUrl, nil

}

func FetchEvents() ([]types.Event, error) {
	apiBaseUrl, err := LoadEnv()
	if err != nil {
		return nil, err
	}
	apiUrl := apiBaseUrl + "/events"
	log.Printf("Fetching events from :%s", apiUrl)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch events from %s, Error: %w", apiUrl, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read events response body: %w", err)
	}

	var events []types.Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("Failed to parse response, %w", err)
	}
	return events, nil
}
