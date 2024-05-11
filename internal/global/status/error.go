package status

import (
	"strconv"
)

// Err contains status code for http.
type Error struct {
	StatusCode int
	Message    string
}

func NewErr(status int, msg string) Error {
	return Error{
		StatusCode: status,
		Message:    msg,
	}
}

func (e Error) Error() string {
	return strconv.Itoa(e.StatusCode) + ": " + e.Message
}
