package utils

import (
	"io"
	"os"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const maxDebugLines = 50 // Increased to allow more logs to be captured
const viewportHeight = 7 // Explicitly set viewport height to 7 lines

var (
	// Define the label style for "Logs"
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true).
			Padding(0, 1)
	logViewPortStyle = lipgloss.NewStyle() // reset to defaults
)

// memoryWriter is a custom io.Writer that captures logs in memory and writes to a file
type memoryWriter struct {
	file     *os.File
	lines    []string
	maxLines int
	mu       sync.Mutex // For thread safety
}

func (w *memoryWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Write to file if file is set
	if w.file != nil {
		n, err = w.file.Write(p)
		if err != nil {
			return n, err
		}
	}

	// Capture in memory
	line := strings.TrimSpace(string(p))
	w.lines = append(w.lines, line)
	if len(w.lines) > w.maxLines {
		w.lines = w.lines[len(w.lines)-w.maxLines:]
	}

	return len(p), nil
}

// GetLines returns the captured log lines
func (w *memoryWriter) GetLines() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.lines
}

// Logger encapsulates all logging-related functionality
type Logger struct {
	logger     *log.Logger
	debug      bool
	showDebug  bool
	writer     *memoryWriter
	viewport   viewport.Model // For scrollable debug panel
	labelStyle lipgloss.Style // For the "Logs" label
}

// NewLogger initializes a new Logger instance
func NewLogger(debug bool) *Logger {
	writer := &memoryWriter{
		lines:    make([]string, 0, maxDebugLines),
		maxLines: maxDebugLines,
	}

	var file *os.File
	var err error
	if debug {
		// Create or open the log file
		file, err = os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		writer.file = file
	}

	// Use only the memoryWriter as the output
	var output io.Writer = writer
	if !debug {
		output = io.Discard
	}

	// Initialize the logger
	logger := log.NewWithOptions(output, log.Options{
		ReportTimestamp: true,
		ReportCaller:    true,
		Level:           log.DebugLevel,
	})
	if !debug {
		logger.SetLevel(log.FatalLevel)
	}

	if debug {
		logger.Info("Application started")
	}

	// Initialize the viewport
	vp := viewport.New(0, viewportHeight) // Width will be set dynamically, height is 7 lines
	vp.Style = lipgloss.NewStyle().
		Background(lipgloss.Color("#222222")).
		Foreground(lipgloss.Color("#FFFFFF"))

	return &Logger{
		logger:     logger,
		debug:      debug,
		showDebug:  false,
		writer:     writer,
		labelStyle: labelStyle,
		viewport:   vp,
	}
}

// Close closes the log file if it is open
func (l *Logger) Close() error {
	if l.writer.file != nil {
		return l.writer.file.Close()
	}
	return nil
}

// GetLogger returns the logger instance
func (l *Logger) GetLogger() *log.Logger {
	return l.logger
}

// IsDebugEnabled returns whether debugging is enabled
func (l *Logger) IsDebugEnabled() bool {
	return l.debug
}

// ToggleDebugView toggles the debug view on or off
func (l *Logger) ToggleDebugView() {
	l.showDebug = !l.showDebug
	l.logger.Info("Toggled debug view", "showDebug", l.showDebug)
	if l.showDebug {
		// Ensure the viewport content is updated when shown
		l.UpdateViewportContent()
	}
}

// IsDebugViewShown returns whether the debug view is currently shown
func (l *Logger) IsDebugViewShown() bool {
	return l.showDebug
}

// DebugPanelHeight returns the height of the debug panel (fixed height for viewport)
func (l *Logger) GetDebugPanelHeight() int {
	return viewportHeight
}

// UpdateViewportContent updates the viewport content with the latest logs
func (l *Logger) UpdateViewportContent() {
	debugContent := strings.Join(l.writer.GetLines(), "\n")
	l.viewport.SetContent(debugContent)
	// Ensure the viewport scrolls to the bottom to show the latest logs
	l.viewport.SetYOffset(max(0, len(l.writer.GetLines())-viewportHeight))
}

// UpdateViewport updates the viewport dimensions and content
func (l *Logger) UpdateViewport(width int) {
	l.viewport.Width = width - 2 // Account for borders
	l.UpdateViewportContent()
}

// ScrollUp scrolls the debug panel up
func (l *Logger) ScrollUp() {
	l.viewport.ScrollUp(1)
}

// ScrollDown scrolls the debug panel down
func (l *Logger) ScrollDown() {
	l.viewport.ScrollDown(1)
}

// RenderDebugPanel renders the debug panel with the specified width and a labeled border
func (l *Logger) RenderDebugPanel(width int) string {
	l.UpdateViewport(width)
	// Render the viewport content
	content := l.viewport.View()
	l.viewport.Style = logViewPortStyle
	// Create the labeled border with "Logs"
	label := l.labelStyle.Render("Logs")
	borderWidth := width - 2 // Adjust for left and right border characters
	leftPadding := 1
	rightPadding := borderWidth - lipgloss.Width(label) - leftPadding
	if rightPadding < 0 {
		rightPadding = 0
	}
	topBorder := "┏" + strings.Repeat("━", leftPadding) + label + strings.Repeat("━", rightPadding) + "┓"
	return lipgloss.JoinVertical(lipgloss.Center, topBorder, content)
}
