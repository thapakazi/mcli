package profile

import (
	"mcli/internal/types"
	"time"
)

type UserProfile struct {
	UserID     string
	Location   string
	Bookmarks  []types.EventId
	ReadEvents []types.EventId
	Filters    map[string]string
	CreatedAt  time.Time
	UpdatedAt  time.Time
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
