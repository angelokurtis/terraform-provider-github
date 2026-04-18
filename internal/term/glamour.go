package term

import (
	"fmt"

	"github.com/charmbracelet/glamour"
)

// NewGlamourTermRenderer initializes and returns a new Glamour terminal renderer.
func NewGlamourTermRenderer() (*glamour.TermRenderer, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(), // Enable automatic style detection.
		glamour.WithWordWrap(0), // Disable word wrapping (0 means no wrapping).
	)
	if err != nil {
		return nil, fmt.Errorf("unable to create renderer: %w", err)
	}

	return renderer, nil
}
