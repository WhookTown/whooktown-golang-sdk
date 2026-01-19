package whooktown

import (
	"errors"
	"fmt"
)

// ErrorCode represents the type of error
type ErrorCode string

const (
	ErrUnauthorized   ErrorCode = "unauthorized"
	ErrForbidden      ErrorCode = "forbidden"
	ErrNotFound       ErrorCode = "not_found"
	ErrBadRequest     ErrorCode = "bad_request"
	ErrQuotaExceeded  ErrorCode = "quota_exceeded"
	ErrInternalServer ErrorCode = "internal_server"
	ErrNetworkError   ErrorCode = "network_error"
	ErrValidation     ErrorCode = "validation_error"
	ErrTimeout        ErrorCode = "timeout"
)

// Error is the SDK error type
type Error struct {
	Code       ErrorCode
	Message    string
	StatusCode int
	Details    map[string]interface{}
	Cause      error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Cause
}

// QuotaError represents a quota exceeded error with details
type QuotaError struct {
	Code       ErrorCode
	Message    string
	StatusCode int
	Plan       string
	Current    int
	Limit      int
	QuotaType  string // "assets" or "layouts"
}

func (e *QuotaError) Error() string {
	return fmt.Sprintf("%s: %s (plan: %s, current: %d, limit: %d)", e.Code, e.Message, e.Plan, e.Current, e.Limit)
}

// NewError creates a new SDK error
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewErrorWithCause creates a new SDK error with a cause
func NewErrorWithCause(code ErrorCode, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// IsUnauthorized checks if the error is an authorization error (401)
func IsUnauthorized(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == ErrUnauthorized
	}
	return false
}

// IsForbidden checks if the error is a forbidden error (403)
func IsForbidden(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == ErrForbidden
	}
	return false
}

// IsNotFound checks if the error is a not found error (404)
func IsNotFound(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == ErrNotFound
	}
	return false
}

// IsBadRequest checks if the error is a bad request error (400)
func IsBadRequest(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == ErrBadRequest
	}
	return false
}

// IsQuotaExceeded checks if the error is a quota exceeded error (402)
func IsQuotaExceeded(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == ErrQuotaExceeded
	}
	var qe *QuotaError
	if errors.As(err, &qe) {
		return true
	}
	return false
}

// IsNetworkError checks if the error is a network error
func IsNetworkError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == ErrNetworkError
	}
	return false
}

// IsTimeout checks if the error is a timeout error
func IsTimeout(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == ErrTimeout
	}
	return false
}

// GetErrorCode extracts the error code from an error
func GetErrorCode(err error) (ErrorCode, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e.Code, true
	}
	return "", false
}

// GetStatusCode extracts the HTTP status code from an error
func GetStatusCode(err error) (int, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e.StatusCode, true
	}
	return 0, false
}
