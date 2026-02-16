package cmdprompt

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// CommandHandler is a function type for processing commands.
// It takes the entered command string and returns the output string and a tea.Cmd.
type CommandHandler func(string) (string, tea.Cmd)

// CommandPrompt is a reusable component for a tmux-like command prompt.
type CommandPrompt struct {
	textInput     textinput.Model
	active        bool
	output        string
	activationKey string
}

// New creates a new CommandPrompt with default settings.
// activationKey is the key to toggle the command prompt (e.g., ":").
// handler is the function to process entered commands.
func New(activationKey string, handler CommandHandler) *CommandPrompt {
	ti := textinput.New()
	ti.Placeholder = "Enter command or press ESC to cancel"
	ti.Prompt = "☯︎: "
	ti.CharLimit = 156
	ti.Width = 50
	return &CommandPrompt{
		textInput:     ti,
		active:        false,
		output:        "",
		activationKey: activationKey,
	}
}

// Init initializes the CommandPrompt.
func (c *CommandPrompt) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages for the CommandPrompt.
// It returns whether the message was consumed, the updated CommandPrompt, and any tea.Cmd.
// The parent model should call this in its Update function.
func (c *CommandPrompt) Update(msg tea.Msg, handler CommandHandler) (bool, *CommandPrompt, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !c.active {
			// Check for activation key
			if msg.String() == c.activationKey {
				c.active = true
				c.textInput.Reset()
				c.textInput.Focus()
				return true, c, nil
			}
			return false, c, nil
		}
		// Command mode active
		switch msg.String() {
		case "enter":
			// Process the command
			command := c.textInput.Value()
			c.active = false
			c.textInput.Blur()
			if handler != nil {
				output, cmd := handler(command)
				c.output = output
				return true, c, cmd
			}
			c.output = "No command handler provided"
			return true, c, nil
		case "esc":
			// Cancel command mode
			c.active = false
			c.textInput.Blur()
			c.output = "Command cancelled"
			return true, c, nil
		}
		// Update text input
		var cmd tea.Cmd
		c.textInput, cmd = c.textInput.Update(msg)
		return true, c, cmd
	}
	return false, c, nil
}

// View renders the CommandPrompt.
// It returns the rendered string to be included in the parent model's View.
func (c *CommandPrompt) View() string {
	if c.active {
		return fmt.Sprintf("%s", c.textInput.View())
	}
	if c.output != "" {
		return fmt.Sprintf("🚀: %s", c.output)
	}
	return ""
}

// IsActive returns whether the command prompt is currently active.
func (c *CommandPrompt) IsActive() bool {
	return c.active
}

// SetOutput sets the output message manually.
func (c *CommandPrompt) SetOutput(output string) {
	c.output = output
}

// SetOutput sets the output message manually.
func (c *CommandPrompt) GetOutput() string {
	return c.output
}
