package profile

import (
	"encoding/json"
	"fmt"
	"mcli/internal/types"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultDataDir = "data/profiles"

type UserProfile struct {
	UserID     string            `json:"userId"`
	Location   string            `json:"location,omitempty"`
	Bookmarks  []types.EventId   `json:"bookmarks,omitempty"`
	ReadEvents []types.EventId   `json:"readEvents,omitempty"`
	Filters    map[string]string `json:"filters,omitempty"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

// New creates a fresh profile for the given user ID
func New(userID string) *UserProfile {
	now := time.Now()
	return &UserProfile{
		UserID:     userID,
		Bookmarks:  []types.EventId{},
		ReadEvents: []types.EventId{},
		Filters:    map[string]string{},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// profilePath returns the file path for a user's profile JSON
func profilePath(userID string) string {
	// SHA256 fingerprints contain colons and slashes — sanitize for filename
	safe := strings.NewReplacer("/", "_", ":", "_").Replace(userID)
	return filepath.Join(defaultDataDir, safe+".json")
}

// Load reads a profile from disk, or creates a new one if it doesn't exist
func Load(userID string) (*UserProfile, error) {
	path := profilePath(userID)

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return New(userID), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read profile %s: %w", path, err)
	}

	var p UserProfile
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("failed to parse profile %s: %w", path, err)
	}
	return &p, nil
}

// Save writes the profile to disk
func (p *UserProfile) Save() error {
	p.UpdatedAt = time.Now()
	path := profilePath(p.UserID)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create profile dir: %w", err)
	}

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	return os.WriteFile(path, data, 0644)
}

// ToggleBookmark adds or removes an event ID from bookmarks
func (p *UserProfile) ToggleBookmark(eventID types.EventId) bool {
	for i, id := range p.Bookmarks {
		if id == eventID {
			p.Bookmarks = append(p.Bookmarks[:i], p.Bookmarks[i+1:]...)
			return false // removed
		}
	}
	p.Bookmarks = append(p.Bookmarks, eventID)
	return true // added
}

// IsBookmarked checks if an event is bookmarked
func (p *UserProfile) IsBookmarked(eventID types.EventId) bool {
	for _, id := range p.Bookmarks {
		if id == eventID {
			return true
		}
	}
	return false
}

// MarkRead adds an event ID to the read list
func (p *UserProfile) MarkRead(eventID types.EventId) {
	for _, id := range p.ReadEvents {
		if id == eventID {
			return
		}
	}
	p.ReadEvents = append(p.ReadEvents, eventID)
}

// IsRead checks if an event has been read
func (p *UserProfile) IsRead(eventID types.EventId) bool {
	for _, id := range p.ReadEvents {
		if id == eventID {
			return true
		}
	}
	return false
}
