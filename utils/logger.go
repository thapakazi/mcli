package utils

import (
	"flag"
	"io"
	"log/slog"
	"os"
)

// Logger is the global structured logger instance.
var Logger *slog.Logger

// InitLogger initializes the global logger based on the --debug flag.
func InitLogger() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug logging to debug.log")
	flag.Parse()

	// Define handler options to exclude time and level
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug, // Accept all levels
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Drop time and level attributes
			if a.Key == slog.TimeKey || a.Key == slog.LevelKey {
				return slog.Attr{}
			}
			return a
		},
	}

	var handler slog.Handler
	if debug {
		// Open or create debug.log for writing
		file, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("Failed to open debug.log: " + err.Error())
		}
		// Use JSONHandler with custom options
		handler = slog.NewJSONHandler(file, opts)
	} else {
		// Use io.Discard for no output
		handler = slog.NewJSONHandler(io.Discard, opts)
	}

	// Initialize the global logger
	Logger = slog.New(handler)
}
