// types/types.go
package types

// Event represents the structure of an event from the /events endpoint.
type Event struct {
    ID              string      `json:"id"`
    Title           string      `json:"title"`
    Description     *string     `json:"description"`
    DateTime        string      `json:"dateTime"`
    EventType       string      `json:"eventType"`
    URL             string      `json:"url"`
    ImageURL        *string     `json:"imageUrl"`
    Source          string      `json:"source"`
    VenueName       string      `json:"venueName"`
    VenueAddress    string      `json:"venueAddress"`
    VenueCity       string      `json:"venueCity"`
    VenueCountry    string      `json:"venueCountry"`
    VenueState      string      `json:"venueState"`
    VenueLat        float64     `json:"venueLat"`
    VenueLong       float64     `json:"venueLong"`
    OrganizerID     string      `json:"organizerId"`
    OrganizerName   string      `json:"organizerName"`
    TicketCount     int         `json:"ticketCount"`
    TicketRemaining int         `json:"ticketRemaining"`
    TicketPrice     interface{} `json:"ticketPrice"`
}

// EventDetail represents the detailed response from /meetup/:id or /luma/:id.
type EventDetail struct {
    ID              string      `json:"id"`
    Title           string      `json:"title"`
    GroupName       string      `json:"groupName,omitempty"` // Meetup-specific
    Description     string      `json:"description"`
    DateTime        string      `json:"dateTime"`
    EventType       string      `json:"eventType"`
    URL             string      `json:"url"`
    ImageURL        *string     `json:"imageUrl"`
    Source          string      `json:"source,omitempty"` // Luma-specific
    VenueName       string      `json:"venueName"`
    VenueAddress    string      `json:"venueAddress"`
    City            string      `json:"city"`      // Meetup uses "city", Luma uses "venueCity"
    State           string      `json:"state"`     // Meetup uses "state", Luma uses "venueState"
    Country         string      `json:"country"`   // Meetup uses "country", Luma uses "venueCountry"
    VenueLat        float64     `json:"venueLat,omitempty"`
    VenueLong       float64     `json:"venueLong,omitempty"`
    OrganizerID     string      `json:"organizerId,omitempty"` // Luma-specific
    OrganizerName   string      `json:"organizerName,omitempty"`
    TicketCount     int         `json:"ticketCount,omitempty"`
    TicketRemaining int         `json:"ticketRemaining,omitempty"`
    TicketPrice     interface{} `json:"ticketPrice,omitempty"`
    RsvpsCount      int         `json:"rsvpsCount,omitempty"` // Meetup-specific
}

// Messages for Bubble Tea.
type ErrMsg struct{ Err error }
type EventsMsg []Event
type EventDetailMsg struct {
    Detail *EventDetail
    Err    error
}
