package domain

import "fmt"

type ErrorCode string

const (
    CodeValidation ErrorCode = "VALIDATION_ERROR"
    CodeNotFound   ErrorCode = "NOT_FOUND"
    CodeConflict   ErrorCode = "CONFLICT"
    CodeInternal   ErrorCode = "INTERNAL"
)

// AppError is a typed error that can be safely returned to HTTP clients.
// Fields is optional (used mainly for validation errors).
type AppError struct {
    Code    ErrorCode
    Message string
    Fields  map[string]string
}

func (e *AppError) Error() string {
    if e == nil {
        return ""
    }
    if len(e.Fields) > 0 {
        return fmt.Sprintf("%s: %s", e.Code, e.Message)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewValidationError(message string, fields map[string]string) *AppError {
    return &AppError{Code: CodeValidation, Message: message, Fields: fields}
}

func NewNotFound(message string) *AppError {
    return &AppError{Code: CodeNotFound, Message: message}
}

func NewConflict(message string) *AppError {
    return &AppError{Code: CodeConflict, Message: message}
}

func NewInternal(message string) *AppError {
    return &AppError{Code: CodeInternal, Message: message}
}
