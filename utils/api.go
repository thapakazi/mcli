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

func fetchEventByLocation(location string) error {
	apiBaseUrl, err := LoadEnv()
	if err != nil {
		return err
	}
	apiUrl := apiBaseUrl + "/fetch" + "?location=" + location
	Logger.Info("Fetching events for", location, location)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(apiUrl)
	if err != nil {
		return fmt.Errorf("Failed to fetch events for %s, Error: %w", location, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read events response body: %w", err)
	}

	// ugly grok shit
	// Unmarshal the response into a map to avoid struct
	var response map[string]string
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Log the message field
	if message, ok := response["message"]; ok {
		Logger.Info("Fetch response:", message)
	} else {
		return fmt.Errorf("response does not contain message field")
	}

	// TODO: reply user to refresh ui in a while
	return nil
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

func FetchEventByLocationCmd(location string) tea.Msg {
	err := fetchEventByLocation(location)
	if err != nil {
		return FetchErrorMsg{Err: err}
	}
	return nil
}

// sortByDate sorts a slice of Events by DateTime in descending order

// sortByDate sorts a slice of Events by DateTime in descending order
func sortByDate(events types.Events) types.Events {
	// Get current date at midnight for comparison
	currentDate := time.Now().Truncate(24 * time.Hour)

	// Filter out events older than today
	filteredEvents := types.Events{}
	for _, event := range events {
		eventTime, err := parseDateTime(event.DateTime)
		if err != nil {
			continue // Skip invalid dates
		}
		if !eventTime.Before(currentDate) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	// Sort filtered events in ascending order
	sort.Slice(filteredEvents, func(i, j int) bool {
		timeI, errI := parseDateTime(filteredEvents[i].DateTime)
		timeJ, errJ := parseDateTime(filteredEvents[j].DateTime)

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

		// Sort by date in ascending order
		return timeI.Before(timeJ)
	})

	return filteredEvents
}
