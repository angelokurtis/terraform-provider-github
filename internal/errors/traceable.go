package errors

// Traceable defines an interface for errors that can provide detailed
// stack trace information. Implementing this interface allows errors
// to carry stack trace data that can be retrieved and analyzed.
//
// Example usage:
//
//	if t, ok := err.(Traceable); ok {
//	    fmt.Println(t.Stack())
//	}
type Traceable interface {
	// Stack retrieves the call stack associated with this error.
	// It returns a pointer to a Stack, which contains the series of
	// function calls that led to the error.
	Stack() *Stack
}
