// main.go
package main

import (
    "flag"
    "fmt"
    "log"
    "os"

    "github.com/thapakazi/mcli/api"
    "github.com/thapakazi/mcli/model"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/ssh"
    "github.com/charmbracelet/wish"
    "github.com/charmbracelet/wish/bubbletea"
)

// teaHandler creates a Bubble Tea program for the Wish server.
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
    m := model.NewModel()
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
    // Define command-line flags
    wishMode := flag.Bool("wish", false, "Run as a Charm Wish SSH server instead of CLI")
    host := flag.String("host", "localhost", "Host address for the Wish server")
    port := flag.String("port", "2222", "Port for the Wish server")
    flag.Parse()

    // Validate API connectivity before starting
    if _, err := api.LoadEnv(); err != nil {
        log.Fatalf("Error loading environment: %v", err)
    }

    if *wishMode {
        // Run as Wish SSH server
        if err := runWishServer(*host, *port); err != nil {
            log.Fatalf("Error running Wish server: %v", err)
        }
    } else {
        // Run as CLI
        p := tea.NewProgram(model.NewModel(),
            tea.WithInput(os.Stdin),
            tea.WithOutput(os.Stdout),
        )
        if _, err := p.Run(); err != nil {
            fmt.Fprintf(os.Stderr, "Error running CLI program: %v\n", err)
            os.Exit(1)
        }
    }
}
