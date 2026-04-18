package errors

import (
	"errors"
	"fmt"
	"io"
)

// New creates a new error with the provided message and captures the current call stack.
func New(msg string) error {
	return &wrapper{
		err:   errors.New(msg),
		stack: callers(),
	}
}

// WithStack wraps an existing error with a captured call stack.
// It attempts to reuse the existing stack trace if the error is already wrapped using this package (`errors.wrapper`).
func WithStack(err error) error {
	if err == nil {
		return nil
	}

	if w := new(wrapper); errors.As(err, &w) {
		return err
	}

	return &wrapper{
		err:   err,
		stack: callers(),
	}
}

// Errorf creates a new error using fmt.Errorf and attempts to capture the call stack from any argument that implements the Traceable interface.
// If no argument implements Traceable, it falls back to capturing the current call stack.
func Errorf(msg string, args ...any) error {
	e := fmt.Errorf(msg, args...)

	for _, arg := range args {
		if t, ok := arg.(Traceable); ok {
			return &wrapper{
				err:   e,
				stack: t.Stack(),
			}
		}
	}

	return &wrapper{
		err:   e,
		stack: callers(),
	}
}

// wrapper struct is used to wrap an underlying error and store its call stack information.
type wrapper struct {
	err   error
	stack *Stack
}

// Error implements the built-in error interface.
// It returns the string representation of the underlying error wrapped by this instance.
func (t *wrapper) Error() string {
	if t.err != nil {
		return t.err.Error()
	}

	return ""
}

// Unwrap retrieves the cause (original) error wrapped by this instance.
func (t *wrapper) Unwrap() error {
	return t.err
}

// Stack retrieves the call stack information associated with this error.
func (t *wrapper) Stack() *Stack {
	return t.stack
}

// Format implements the fmt.Formatter interface for the wrapper struct.
// It customizes the formatting of the wrapped error and stack trace based on the verb and flags.
func (t *wrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = io.WriteString(s, t.Error())
			_, _ = fmt.Fprintf(s, "%s", t.Stack())

			return
		}

		fallthrough
	case 's':
		_, _ = io.WriteString(s, t.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", t.Error())
	}
}
