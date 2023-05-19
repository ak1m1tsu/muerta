// Package errors provides a way to wrap errors with additional context.
//
// Example usage:
//
//	err := errors.New("an error occurred")
//	wrappedErr := err.With(fmt.Errorf("additional context"))
//	fmt.Println(wrappedErr.Error())
//
// Output: an error occurred: additional context
package errors

import (
	"fmt"
	"strings"
)

// The ErrorWrapper interface represents an error that can have additional context added to it with the With method. The Error method returns a string representation of the error.
type ErrorWrapper interface {
	With(err error) ErrorWrapper
	Error() string
}

type errorWrapper struct {
	err string
}

// The With method of the errorWrapper struct returns a new ErrorWrapper that combines the original error message with the given error. The Error method returns the combined error message.
func (e *errorWrapper) With(err error) ErrorWrapper {
	return &errorWrapper{err: fmt.Sprintf("%s: %s", e.err, err.Error())}
}

func (e *errorWrapper) Error() string {
	return e.err
}

// Contains returns true if the error message contains the given error message.
func (e *errorWrapper) Contains(err ErrorWrapper) bool {
	return strings.Contains(e.err, err.Error())
}

// The New function creates a new ErrorWrapper with the given error message.
func New(err string) ErrorWrapper {
	return &errorWrapper{err: err}
}
