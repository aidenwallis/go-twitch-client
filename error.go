package twitch

import "fmt"

// Error wraps a Twitch API error
type Error struct {
	Message string
	Status  int
}

// NewError creates a new instance of error
func NewError(message string, status int) Error {
	return Error{
		Message: message,
		Status:  status,
	}
}

// Error returns a stringified error
func (e Error) Error() string {
	message := e.Message
	if message == "" {
		message = "<unknown>"
	}

	return fmt.Sprintf("[%d] %s", e.Status, message)
}
