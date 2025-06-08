package utils

import (
	"encoding/json"
	"net/http"
	"os"
)

type CustomError struct {
	Message        string `json:"message"`
	ErrorCode      int    `json:"errorCode,omitempty"`
	ErrorType      string `json:"error"` // Add this!
	HTTPStatusCode int    `json:"httpStatusCode,omitempty"`
	Service        string `json:"service,omitempty"`
	Success        bool   `json:"success,omitempty"`
}

func (e *CustomError) Error() string {
	return e.Message
}

var serviceName = os.Getenv("SERVICE_NAME")

func NewUnauthorizedError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      401,
		ErrorType:      "VALIDATION_ERROR",
		HTTPStatusCode: http.StatusUnauthorized,
		Service:        serviceName,
		Success:        false,
	}
}

func NewBadRequestError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      400,
		ErrorType:      "VALIDATION_ERROR",
		HTTPStatusCode: http.StatusBadRequest,
		Service:        serviceName,
		Success:        false,
	}
}

func NewConflictError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      409,
		ErrorType:      "CONFLICT_ERROR",
		HTTPStatusCode: http.StatusConflict,
		Service:        serviceName,
		Success:        false,
	}
}

func NewInternalServerError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      500,
		ErrorType:      "VALIDATION_ERROR",
		HTTPStatusCode: http.StatusInternalServerError,
		Service:        serviceName,
		Success:        false,
	}
}

func NewUnauthenticatedError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      401,
		ErrorType:      "VALIDATION_ERROR",
		HTTPStatusCode: http.StatusUnauthorized,
		Service:        serviceName,
		Success:        false,
	}
}

func NewNotFoundError(message string) *CustomError {
	return &CustomError{
		Message:        message,
		ErrorCode:      404,
		ErrorType:      "VALIDATION_ERROR",
		HTTPStatusCode: http.StatusNotFound,
		Service:        serviceName,
		Success:        false,
	}
}

func (e *CustomError) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
