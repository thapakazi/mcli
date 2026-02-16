package types

type EventId string

type Location struct {
	VenueAddress string `json:"venueAddress"`
	VenueName    string `json:"venueName"`
}
type Event struct {
	ID          EventId `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Url         string  `json:"url"`
	DateTime    string  `json:"dateTime"`
	Source      string  `json:"source"`
	EventMeta
	Location
}

type EventMeta struct {
	Status     string `json:"status"`
	Type       string `json:"eventType"`
	RsvpsCount int    `json:"rsvpCount"`
}

type Events []Event

type ErrMsg struct{ Err error }

type EventsMsg struct {
	Events []Event
	Err    error
}
