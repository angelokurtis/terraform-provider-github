package term

import (
	"fmt"

	"github.com/charmbracelet/glamour"
	yaml "github.com/goccy/go-yaml"
)

// Renderer interface defines methods for rendering different formats: Markdown, Code, and YAML
type Renderer interface {
	RenderMarkdown(markdown string) error
	RenderCode(code, lang string) error
	RenderYAML(value any) error
}

// MarkdownRenderer struct holds a reference to the glamour renderer to render markdown
type MarkdownRenderer struct {
	glamourRenderer *glamour.TermRenderer
}

// NewMarkdownRenderer creates a new MarkdownRenderer instance with the given glamour renderer
func NewMarkdownRenderer(renderer *glamour.TermRenderer) *MarkdownRenderer {
	return &MarkdownRenderer{glamourRenderer: renderer}
}

// RenderMarkdown renders a given markdown string and prints it to the terminal
func (m *MarkdownRenderer) RenderMarkdown(markdown string) error {
	// Render the markdown using the glamour renderer
	output, err := m.glamourRenderer.Render(markdown)
	if err != nil {
		return fmt.Errorf("failed to render markdown using glamour renderer: %w", err)
	}

	// Print the rendered output to the terminal
	_, err = fmt.Println(output)
	if err != nil {
		return fmt.Errorf("failed to print rendered markdown output: %w", err)
	}

	return nil
}

// RenderCode renders a given code string in a specified language using markdown formatting
func (m *MarkdownRenderer) RenderCode(code, lang string) error {
	// Format the code as markdown with code blocks
	markdown := fmt.Sprintf("```%s\n%s\n```", lang, code)
	return m.RenderMarkdown(markdown)
}

// RenderYAML renders the YAML representation of a given value and prints it as code
func (m *MarkdownRenderer) RenderYAML(value any) error {
	// Marshal the value into YAML format
	data, err := yaml.Marshal(value)
	if err != nil {
		return fmt.Errorf("YAML marshalling failed: %w", err)
	}

	// Render the YAML as code in markdown format
	return m.RenderCode(string(data), "yaml")
}
