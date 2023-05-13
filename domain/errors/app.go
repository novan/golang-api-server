package errors

import (
	"errors"
	"net/http"
)

const (
	// BadRequest error indicates a bad request from client
	BadRequest = "BadRequest"

	// Forbidden error indicates a forbidden access from client
	Forbidden = "Forbidden"

	// Forbidden error indicates a forbidden access from client
	MethodNotAllowed = "MethodNotAllowed"

	// NotFound error indicates a missing / not found record
	NotFound = "NotFound"

	// ValidationError indicates an error in input validation
	ValidationError = "ValidationError"

	// ResourceAlreadyExists indicates a duplicate / already existing record
	ResourceAlreadyExists = "ResourceAlreadyExists"

	// RepositoryError indicates a repository (e.g database) error
	RepositoryError = "RepositoryError"

	// NotAuthenticated indicates an authentication error
	NotAuthenticated = "NotAuthenticated"

	// TokenGeneratorError indicates an token generation error
	TokenGeneratorError = "TokenGeneratorError"

	// NotAuthorized indicates an authorization error
	NotAuthorized = "NotAuthorized"

	// InternalError indicates an error that the app cannot find the cause for
	InternalError = "InternalError"
)

var errorDescription = map[string]string{
	BadRequest:            "Bad Request",
	Forbidden:             "Forbidden",
	MethodNotAllowed:      "Method not allowed",
	NotFound:              "Resource not found",
	ValidationError:       "Validation error",
	ResourceAlreadyExists: "Resource already exists",
	RepositoryError:       "Error in repository operation",
	NotAuthenticated:      "Not Authenticated",
	TokenGeneratorError:   "Error in token generation",
	NotAuthorized:         "Not Authorized",
	InternalError:         "Something went wrong",
}

// AppError defines an application (domain) error
type AppError struct {
	Err  error
	Type string
}

// NewAppError initializes a new domain error using an error and its type.
func NewAppError(err error, errType string) *AppError {
	return &AppError{
		Err:  err,
		Type: errType,
	}
}

// NewAppErrorWithType initializes a new default error for a given type.
func NewAppErrorWithType(errType string) *AppError {
	var err error

	if msg, ok := errorDescription[errType]; ok {
		return &AppError{
			Err:  errors.New(msg),
			Type: errType,
		}
	}

	err = errors.New(errorDescription[InternalError])

	return &AppError{
		Err:  err,
		Type: errType,
	}
}

// String converts the app error to a human-readable string.
func (appErr *AppError) Error() string {
	return appErr.Err.Error()
}

// InternalError error
func InternalServerError(message string) error {
	return NewAppError(errors.New(message), InternalError)
}

func NewAppHttpError(err error, status int) *AppError {
	var errType string

	switch status {
	case http.StatusBadRequest:
		errType = BadRequest
	case http.StatusNotFound:
		errType = NotFound
	case http.StatusForbidden:
		errType = Forbidden
	case http.StatusMethodNotAllowed:
		errType = MethodNotAllowed
	default:
		errType = InternalError
	}

	return NewAppError(err, errType)
}