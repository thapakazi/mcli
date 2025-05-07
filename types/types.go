package types

type Event struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Location    string `json:"venueAddress"`
	DateTime    string `json:"dateTime"`
	Source      string `json:"source"`
}

type ErrMsg struct{ Err error }

type EventsMsg struct {
	Events []Event
	Err    error
}
