package unbrokenwing

import (
	"fmt"
	"regexp"
)

// FailureError are suitable to be
// returned from behaviour descriptions
// that failes during test.
type FailureError struct {
	reason string
}

func (e FailureError) Error() string {
	return e.reason
}

// Failure returns error that
// are suitable to be returned from
// behaviour descriptions that failes
// during test.
func Failure(reason string) error {
	return FailureError{reason: reason}
}

// PendingError are suitable to be
// returned from behaviour descriptions
// that are generated but not implemented yet.
type PendingError struct {
	reason string
}

func (e PendingError) Error() string {
	return e.reason
}

// Pending returns error that
// are suitable to be returned from
// behaviour descriptions that are
// generated but not implemented yet.
func Pending(reason string) error {
	return PendingError{reason: reason}
}

// NotImplError are suitable to be returned
// when unbrokenwings driver are not able
// to find matching behaviour implementation.
type NotImplError struct{ t Testable }

// Generates a behaviour snippet
// matching missing implementation.
func (e NotImplError) Snippet() string {
	t := e.t.(Step)

	r := regexp.MustCompile("\"([0-9+-]+)\"")
	newRe := r.ReplaceAllString(t.Description, "\\\"([0-9+-]+)\\\"")

	snippet := fmt.Sprintf(`
    %s("^%s$", func(args Args) error {
        return Pending("Not implemented")
    })`, t.Cmd, newRe)

	return snippet
}

func (e NotImplError) Error() string {
	intro := "You can implement step definition with following snippet:"
	return fmt.Sprintf("Not Implemented: %s\n%s\n%s", e.t, intro, e.Snippet())
}

// NotImplemented returns error that
// are suitable to be returned when
// unbrokenwings driver are not able
// to find matching behaviour implementation.
func NotImplemented(t Testable) error {
	return NotImplError{t: t}
}
