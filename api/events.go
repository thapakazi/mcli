// api/api.go
package api

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/thapakazi/mcli/types"

    "github.com/joho/godotenv"
)

// LoadEnv loads the .env file and returns the API base URL.
func LoadEnv() (string, error) {
    err := godotenv.Load()
    if err != nil {
        return "", fmt.Errorf("error loading .env file: %w", err)
    }

    apiBaseURL := os.Getenv("API_BASE_URL")
    if apiBaseURL == "" {
        return "", fmt.Errorf("API_BASE_URL not set in .env file")
    }

    return apiBaseURL, nil
}

// FetchEvents fetches events from the /events endpoint.
func FetchEvents() ([]types.Event, error) {
    apiBaseURL, err := LoadEnv()
    if err != nil {
        return nil, err
    }

    apiURL := apiBaseURL + "/events"
    log.Printf("Fetching events from: %s", apiURL)

    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    resp, err := client.Get(apiURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch events: %w", err)
    }
    defer resp.Body.Close()

    log.Printf("Events API response status: %s", resp.Status)

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch events: status %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read events response: %w", err)
    }

    var events []types.Event
    if err := json.Unmarshal(body, &events); err != nil {
        return nil, fmt.Errorf("failed to unmarshal events: %w", err)
    }

    return events, nil
}

// FetchEventDetail fetches the details of a specific event from /meetup/:id or /luma/:id.
func FetchEventDetail(event types.Event) (*types.EventDetail, error) {
    apiBaseURL, err := LoadEnv()
    if err != nil {
        return nil, err
    }

    var apiURL string
    switch event.Source {
    case "luma":
        apiURL = apiBaseURL + "/luma/" + event.ID
    default:
        apiURL = apiBaseURL + "/meetup/" + event.ID
    }

    //log.Printf("Fetching event details from: %s", apiURL)

    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    resp, err := client.Get(apiURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch event details: %w", err)
    }
    defer resp.Body.Close()

    //log.Printf("Event details API response status: %s", resp.Status)

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch event details: status %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read event details response: %w", err)
    }

    log.Printf("Event details response body: %s", string(body))

    var eventDetail types.EventDetail
    if err := json.Unmarshal(body, &eventDetail); err != nil {
        return nil, fmt.Errorf("failed to unmarshal event details: %w", err)
    }

    return &eventDetail, nil
}
