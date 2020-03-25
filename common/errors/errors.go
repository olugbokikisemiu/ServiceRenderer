package errors

import (
	"fmt"
)

// ErrorCode type error code for all
// possible errors
type ErrorCode string

const (
	ErrUnknown            = ErrorCode("unknown_error")
	ErrInvalidPinFormat   = ErrorCode("invalid_pin_format")
	ErrNoDocumentFound    = ErrorCode("no_document_found")
	ErrBadActor           = ErrorCode("badactor_error")
	ErrJailedLogin        = ErrorCode("login_jailed")
	ErrInvalidCredentials = ErrorCode("invalid_credentials")
	ErrUnableToGetSession = ErrorCode("unable_to_get_session")
	ErrNotFound           = ErrorCode("not_found")
	ErrAuthenticated      = ErrorCode("user_authenticated")
	ErrUnauthenticated    = ErrorCode("permission_denied")
)

// ErrorMessages all error messages to return to user
var ErrorMessages = map[ErrorCode]*Error{
	ErrInvalidPinFormat: {
		ErrMessage: "Only digits are allowed for pin",
	},
	ErrNoDocumentFound: {
		ErrMessage: "mongo: no documents in result",
	},
	ErrBadActor: {
		ErrMessage: "Badactor Error: ",
	},
	ErrJailedLogin: {
		ErrMessage: "You have reach the maximum number of invalid login attempts",
	},
	ErrInvalidCredentials: {
		ErrMessage: "Incorrect pin",
	},
	ErrUnableToGetSession: {
		ErrMessage: "Unable to get user session",
	},
	ErrNotFound: {
		ErrMessage: "Record not found",
	},
	ErrAuthenticated: {
		ErrMessage: "User already authenticated",
	},
	ErrUnauthenticated: {
		ErrMessage: "Permission denied",
	},
	ErrUnknown: {
		ErrMessage: "An unknown error occured",
	},
}

type Error struct {
	ErrMessage string
}

func (self *Error) Error() string {
	return self.ErrMessage
}

func newError(text string, args ...interface{}) *Error {
	return &Error{ErrMessage: fmt.Sprintf(text, args...)}
}

// ErrorLog finds error message in
// ErrorMessage map by key and construct
// a message with newError()
func ErrorLog(key ErrorCode, args ...interface{}) *Error {
	errorMessage, ok := ErrorMessages[key]
	if !ok {
		return newError("Unable to find error with code %s", key)
	}
	return newError(errorMessage.Error(), args)
}
