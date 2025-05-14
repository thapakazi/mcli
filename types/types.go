package types

type EventId string

type Event struct {
	ID          EventId `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Url         string  `json:"url"`
	Location    string  `json:"venueAddress"`
	DateTime    string  `json:"dateTime"`
	Source      string  `json:"source"`
}

type Events []Event

type ErrMsg struct{ Err error }

type EventsMsg struct {
	Events []Event
	Err    error
}
