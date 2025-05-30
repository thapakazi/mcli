package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

	"mcli/types"

	tea "github.com/charmbracelet/bubbletea"
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

func fetchEvents() ([]types.Event, error) {
	apiBaseUrl, err := LoadEnv()
	if err != nil {
		return nil, err
	}
	apiUrl := apiBaseUrl + "/events"
	Logger.Info("Fetching events from", apiUrl, apiUrl)

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

type FetchErrorMsg struct {
	Err error
}

type FetchSuccessMsg struct {
	Events []types.Event
}

func FetchEventCmd() tea.Msg {
	events, err := fetchEvents()
	if err != nil {
		return FetchErrorMsg{Err: err}
	}
	// sort events prior to returning
	sortedEvents := sortByDate(events)
	return FetchSuccessMsg{Events: sortedEvents}
}

// sortByDate sorts a slice of Events by DateTime in descending order
func sortByDate(events types.Events) types.Events {
	// Use sort.Slice to sort events in place
	sort.Slice(events, func(i, j int) bool {
		// Parse DateTime for both events
		timeI, errI := parseDateTime(events[i].DateTime)
		timeJ, errJ := parseDateTime(events[j].DateTime)

		// Handle parsing errors: invalid dates go to the end
		if errI != nil && errJ != nil {
			return false
		}
		if errI != nil {
			return false
		}
		if errJ != nil {
			return true
		}

		// Sort by date in descending order
		return timeI.After(timeJ)
	})

	return events
}
