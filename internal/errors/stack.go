package errors

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"runtime"
	"strings"
)

var functionPattern = regexp.MustCompile(`^(github.com/angelokurtis/terraform-provider-github|main\.)`)

// Stack represents a call stack, a slice of function call program counter (PC) values
type Stack []uintptr

// callers function retrieves the call stack information
func callers() *Stack {
	s := make(Stack, 64)
	n := runtime.Callers(3, s)
	s = s[:n]

	return &s
}

// String method on Stack type formats the call stack into a human-readable string
func (s *Stack) String() string {
	if s == nil {
		return ""
	}

	// Create a new stackBuilder to accumulate formatted frame information
	builder := new(stackBuilder)
	frames := runtime.CallersFrames(*s)

	for more := true; more; {
		var frame runtime.Frame

		frame, more = frames.Next()

		// Skip frames that don't match the function pattern
		if !functionPattern.MatchString(frame.Function) {
			continue
		}

		// Add the frame to the builder if it passes the pattern check
		builder.AddCallerFrame(frame)
	}

	return builder.Build()
}

// StackFormatOptions specifies options for formatting a call stack.
type StackFormatOptions struct {
	FunctionPattern *regexp.Regexp
}

// stackBuilder is a helper struct to build the formatted stack trace string
type stackBuilder struct {
	builder strings.Builder
}

// AddCallerFrame method on stackBuilder adds a single frame information to the string
func (s *stackBuilder) AddCallerFrame(frame runtime.Frame) {
	if _, err := s.builder.WriteString(fmt.Sprintf(
		"\n%s\n\t%s:%d",
		frame.Function,
		frame.File,
		frame.Line,
	)); err != nil {
		slog.WarnContext(context.TODO(), "Error formatting stack trace", "error", err)
	}
}

// Build method on stackBuilder returns the final formatted stack trace string
func (s *stackBuilder) Build() string {
	return s.builder.String()
}
