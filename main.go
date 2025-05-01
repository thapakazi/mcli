package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/textinput"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/ssh"
    "github.com/charmbracelet/wish"
    "github.com/charmbracelet/wish/bubbletea"
    "github.com/joho/godotenv"
)

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
    TicketRemaining int         `ticketRemaining"`
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

// Model holds the application state.
type Model struct {
    events         []Event
    filteredEvents []Event
    cursor         int
    viewportTop    int
    viewportHeight int
    err            error
    loading        bool
    termHeight     int
    filtering      bool
    filterInput    textinput.Model
    filterText     string
    viewMode       string       // "list" or "details"
    selectedEvent  *EventDetail // Details of the selected event
}

// InitialModel initializes the application state.
func InitialModel() Model {
    ti := textinput.New()
    ti.Placeholder = "Filter by title..."
    ti.CharLimit = 50
    ti.Width = 50

    return Model{
        events:         []Event{},
        filteredEvents: []Event{},
        cursor:         0,
        viewportTop:    0,
        viewportHeight: 10,
        loading:        true,
        termHeight:     0,
        filtering:      false,
        filterInput:    ti,
        filterText:     "",
        viewMode:       "list",
        selectedEvent:  nil,
    }
}

// loadEnv loads the .env file and returns the API base URL.
func loadEnv() (string, error) {
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

// fetchEvents fetches events from the /events endpoint.
func fetchEvents() ([]Event, error) {
    apiBaseURL, err := loadEnv()
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

    var events []Event
    if err := json.Unmarshal(body, &events); err != nil {
        return nil, fmt.Errorf("failed to unmarshal events: %w", err)
    }

    return events, nil
}

// fetchEventDetail fetches the details of a specific event from /meetup/:id or /luma/:id.
func fetchEventDetail(event Event) (*EventDetail, error) {
    apiBaseURL, err := loadEnv()
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

    log.Printf("Fetching event details from: %s", apiURL)

    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    resp, err := client.Get(apiURL)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch event details: %w", err)
    }
    defer resp.Body.Close()

    log.Printf("Event details API response status: %s", resp.Status)

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch event details: status %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read event details response: %w", err)
    }

    log.Printf("Event details response body: %s", string(body))

    var eventDetail EventDetail
    if err := json.Unmarshal(body, &eventDetail); err != nil {
        return nil, fmt.Errorf("failed to unmarshal event details: %w", err)
    }

    return &eventDetail, nil
}

// Messages for Bubble Tea.
type errMsg struct{ err error }
type eventsMsg []Event
type eventDetailMsg struct {
    detail *EventDetail
    err    error
}

// Init starts the program by fetching events.
func (m Model) Init() tea.Cmd {
    return tea.Batch(func() tea.Msg {
        events, err := fetchEvents()
        if err != nil {
            return errMsg{err}
        }
        return eventsMsg(events)
    }, tea.WindowSize())
}

// Update handles user input and updates the model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.termHeight = msg.Height
        m.viewportHeight = msg.Height - 4
        if m.filtering {
            m.viewportHeight -= 2 // Reserve space for filter input
        }
        if m.viewportHeight < 1 {
            m.viewportHeight = 1
        }
        m.adjustViewport()
        return m, nil

    case eventsMsg:
        m.events = msg
        m.filteredEvents = m.events
        m.loading = false
        if len(m.events) == 0 {
            m.err = fmt.Errorf("no events found")
        }
        m.adjustViewport()
    case eventDetailMsg:
        m.loading = false
        if msg.err != nil {
            m.err = msg.err
            m.viewMode = "list"
            m.selectedEvent = nil
            return m, nil
        }
        m.selectedEvent = msg.detail
    case errMsg:
        m.err = msg.err
        m.loading = false
        return m, tea.Quit
    case tea.KeyMsg:
        if m.err != nil {
            if msg.String() == "q" || msg.String() == "ctrl+c" {
                return m, tea.Quit
            }
            return m, nil
        }

        if m.filtering {
            switch msg.String() {
            case "esc":
                m.filtering = false
                m.filterText = ""
                m.filterInput.SetValue("")
                m.filteredEvents = m.events
                m.cursor = 0
                m.viewportTop = 0
                m.viewportHeight = m.termHeight - 4
                m.adjustViewport()
            case "enter":
                m.filtering = false
                m.filterText = m.filterInput.Value()
                m.viewportHeight = m.termHeight - 4
                m.adjustViewport()
            default:
                var cmd tea.Cmd
                m.filterInput, cmd = m.filterInput.Update(msg)
                m.filterText = m.filterInput.Value()
                m.filteredEvents = filterEvents(m.events, m.filterText)
                m.cursor = 0
                m.viewportTop = 0
                m.adjustViewport()
                return m, cmd
            }
            return m, nil
        }

        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "/":
            if m.viewMode == "list" {
                m.filtering = true
                m.filterInput.Focus()
                m.viewportHeight = m.termHeight - 6
                m.adjustViewport()
                return m, textinput.Blink
            }
        case "up", "k":
            if m.viewMode == "list" {
                if m.cursor > 0 {
                    m.cursor--
                    m.adjustViewport()
                }
            }
        case "down", "j":
            if m.viewMode == "list" {
                if m.cursor < len(m.filteredEvents)-1 {
                    m.cursor++
                    m.adjustViewport()
                }
            }
        case "enter":
            if m.viewMode == "list" && len(m.filteredEvents) > 0 {
                m.viewMode = "details"
                m.loading = true
                selectedEvent := m.filteredEvents[m.cursor]
                return m, func() tea.Msg {
                    detail, err := fetchEventDetail(selectedEvent)
                    return eventDetailMsg{detail: detail, err: err}
                }
            } else if m.viewMode == "details" {
                m.viewMode = "list"
                m.selectedEvent = nil
                m.loading = false
                m.adjustViewport()
            }
        case "esc":
            if m.viewMode == "details" {
                m.viewMode = "list"
                m.selectedEvent = nil
                m.loading = false
                m.adjustViewport()
            }
        }
    }
    return m, nil
}

// filterEvents filters events based on a search term.
func filterEvents(events []Event, term string) []Event {
    if term == "" {
        return events
    }
    term = strings.ToLower(term)
    var filtered []Event
    for _, event := range events {
        if strings.Contains(strings.ToLower(event.Title), term) {
            filtered = append(filtered, event)
        }
    }
    return filtered
}

// adjustViewport ensures the cursor is visible within the viewport.
func (m *Model) adjustViewport() {
    viewportBottom := m.viewportTop + m.viewportHeight - 1

    if m.cursor < m.viewportTop {
        m.viewportTop = m.cursor
    }

    if m.cursor > viewportBottom {
        m.viewportTop = m.cursor - m.viewportHeight + 1
    }

    if m.viewportTop < 0 {
        m.viewportTop = 0
    }

    if m.viewportTop > len(m.filteredEvents)-m.viewportHeight && len(m.filteredEvents) >= m.viewportHeight {
        m.viewportTop = len(m.filteredEvents) - m.viewportHeight
    }
}

// formatTicketPrice converts ticketPrice to a string for display.
func formatTicketPrice(price interface{}) string {
    if price == nil {
        return "N/A"
    }
    switch v := price.(type) {
    case string:
        return v
    case float64:
        return fmt.Sprintf("$%.2f", v)
    case int:
        return fmt.Sprintf("$%d", v)
    default:
        return fmt.Sprintf("%v", v)
    }
}

// View renders the UI.
func (m Model) View() string {
    if m.loading {
        return "Loading...\n"
    }
    if m.err != nil {
        return fmt.Sprintf("Error: %v\nPress q to quit.", m.err)
    }

    if m.viewMode == "details" && m.selectedEvent != nil {
        return m.renderEventDetails()
    }

    return m.renderEventList()
}

// renderEventList renders the list of events.
func (m Model) renderEventList() string {
    if len(m.events) == 0 {
        return "No events available.\nPress q to quit."
    }

    normalStyle := lipgloss.NewStyle().
        PaddingLeft(2)

    highlightStyle := lipgloss.NewStyle().
        PaddingLeft(2).
        Background(lipgloss.Color("1")).
        Foreground(lipgloss.Color("15")).
        Width(80).
        Align(lipgloss.Left)

    s := strings.Builder{}
    s.WriteString("Events List (Use ↑/↓ to navigate, / to filter, Enter for details, q to quit)\n\n")

    if m.filtering {
        s.WriteString("Filter: " + m.filterInput.View() + "\n\n")
    }

    if len(m.filteredEvents) == 0 {
        s.WriteString("No events match the filter.\n")
        return s.String()
    }

    start := m.viewportTop
    end := m.viewportTop + m.viewportHeight
    if end > len(m.filteredEvents) {
        end = len(m.filteredEvents)
    }

    for i := start; i < end; i++ {
        event := m.filteredEvents[i]
        prefix := "  "
        if m.cursor == i {
            prefix = "> "
        }

        line := fmt.Sprintf("%s%s", prefix, event.Title)
        if m.cursor == i {
            line = highlightStyle.Render(line)
        } else {
            line = normalStyle.Render(line)
        }
        s.WriteString(line + "\n")
    }

    if m.viewportTop > 0 {
        s.WriteString("↑ More events above...\n")
    }
    if end < len(m.filteredEvents) {
        s.WriteString("↓ More events below...\n")
    }

    return s.String()
}

// renderEventDetails renders the details of a selected event with colors.
func (m Model) renderEventDetails() string {
    s := strings.Builder{}

    // Define styles based on the screenshot colors
    headerStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")). // White
        Bold(true)

    titleStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF3333")). // Red
        Bold(true)

    dateStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#00FFFF")) // Cyan

    metadataStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFF00")) // Yellow

    urlStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF33FF")). // Pinkish-purple for URLs (simulating gradient)
        Underline(true)

    bodyStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")). // White
        PaddingLeft(2)

    // Header
    s.WriteString(headerStyle.Render("Event Details (Press Enter or Esc to go back, q to quit)\n\n"))

    if m.selectedEvent == nil {
        return s.String()
    }

    d := m.selectedEvent
    details := []string{}

    // Title
    details = append(details, titleStyle.Render(fmt.Sprintf("Title: %s", d.Title)))

    // Group Name (Meetup-specific)
    if d.GroupName != "" {
        details = append(details, metadataStyle.Render(fmt.Sprintf("Group: %s", d.GroupName)))
    }

    // Date
    details = append(details, dateStyle.Render(fmt.Sprintf("Date: %s", d.DateTime)))

    // Event Type
    details = append(details, metadataStyle.Render(fmt.Sprintf("Type: %s", d.EventType)))

    // Venue
    venue := fmt.Sprintf("Venue: %s, %s, %s", d.VenueName, d.City, d.State)
    details = append(details, dateStyle.Render(venue))

    // Tickets (Luma-specific)
    if d.TicketCount > 0 {
        details = append(details, metadataStyle.Render(fmt.Sprintf("Tickets Remaining: %d/%d", d.TicketRemaining, d.TicketCount)))
    }

    // RSVPs (Meetup-specific)
    if d.RsvpsCount > 0 {
        details = append(details, metadataStyle.Render(fmt.Sprintf("RSVPs: %d", d.RsvpsCount)))
    }

    // Ticket Price
    if d.TicketPrice != nil {
        details = append(details, metadataStyle.Render(fmt.Sprintf("Price: %s", formatTicketPrice(d.TicketPrice))))
    }

    // URL
    details = append(details, urlStyle.Render(fmt.Sprintf("URL: %s", d.URL)))

    // Description
    if d.Description != "" {
        // Process description to apply bold/italic formatting for ** and * markers
        formattedDesc := formatDescription(d.Description)
        details = append(details, fmt.Sprintf("\nDescription:\n%s", formattedDesc))
    }

    // Render the details with body style
    s.WriteString(bodyStyle.Render(strings.Join(details, "\n")))

    return s.String()
}

// formatDescription applies bold/italic formatting to the description based on ** and * markers.
func formatDescription(desc string) string {
    // Define styles for bold and italic
    boldStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Bold(true)

    italicStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF")).
        Italic(true)

    normalStyle := lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFFFFF"))

    // Split the description into lines to preserve formatting
    lines := strings.Split(desc, "\n")
    var formattedLines []string

    for _, line := range lines {
        // Process **bold** markers
        parts := strings.Split(line, "**")
        var lineBuilder strings.Builder
        for i, part := range parts {
            if i%2 == 0 {
                // Outside ** markers, process *italic* markers
                italicParts := strings.Split(part, "*")
                for j, italicPart := range italicParts {
                    if j%2 == 0 {
                        lineBuilder.WriteString(normalStyle.Render(italicPart))
                    } else {
                        lineBuilder.WriteString(italicStyle.Render(italicPart))
                    }
                }
            } else {
                // Inside ** markers, apply bold
                lineBuilder.WriteString(boldStyle.Render(part))
            }
        }
        formattedLines = append(formattedLines, lineBuilder.String())
    }

    return strings.Join(formattedLines, "\n")
}

// teaHandler creates a Bubble Tea program for the Wish server.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
    m := InitialModel()
    opts := []tea.ProgramOption{
        tea.WithInput(s),
        tea.WithOutput(s),
    }
    return m, opts
}

// runWishServer starts a Charm Wish SSH server to serve the Bubble Tea app.
func runWishServer(host, port string) error {
    s, err := wish.NewServer(
        wish.WithAddress(fmt.Sprintf("%s:%s", host, port)),
        wish.WithHostKeyPath(".ssh/events_app_ed25519"),
        wish.WithMiddleware(
            bubbletea.Middleware(teaHandler),
        ),
    )
    if err != nil {
        return fmt.Errorf("could not start Wish server: %w", err)
    }

    log.Printf("Starting Wish SSH server on %s:%s", host, port)
    log.Println("Connect using: ssh -p", port, host)
    return s.ListenAndServe()
}

func main() {
    wishMode := flag.Bool("wish", false, "Run as a Charm Wish SSH server instead of CLI")
    host := flag.String("host", "localhost", "Host address for the Wish server")
    port := flag.String("port", "2222", "Port for the Wish server")
    flag.Parse()

    if *wishMode {
        if err := runWishServer(*host, *port); err != nil {
            log.Fatalf("Error running Wish server: %v", err)
        }
    } else {
        p := tea.NewProgram(InitialModel())
        if err := p.Start(); err != nil {
            fmt.Fprintf(os.Stderr, "Error running CLI program: %v\n", err)
            os.Exit(1)
        }
    }
}
