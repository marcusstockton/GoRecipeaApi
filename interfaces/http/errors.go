package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Details any    `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func NewAPIError(status int, code, message string, details any) *APIError {
	return &APIError{Status: status, Code: code, Message: message, Details: details}
}

func NewBadRequest(message string, details any) *APIError {
	return NewAPIError(http.StatusBadRequest, "bad_request", message, details)
}

func NewUnauthorized(message string) *APIError {
	return NewAPIError(http.StatusUnauthorized, "unauthorized", message, nil)
}

func NewForbidden(message string) *APIError {
	return NewAPIError(http.StatusForbidden, "forbidden", message, nil)
}

func NewNotFound(message string) *APIError {
	return NewAPIError(http.StatusNotFound, "not_found", message, nil)
}

func NewConflict(message string) *APIError {
	return NewAPIError(http.StatusConflict, "conflict", message, nil)
}

func NewInternal(message string) *APIError {
	return NewAPIError(http.StatusInternalServerError, "internal_error", message, nil)
}

func WriteAPIError(c *gin.Context, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		c.JSON(apiErr.Status, gin.H{"error": gin.H{"code": apiErr.Code, "message": apiErr.Message, "details": apiErr.Details}})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "internal_error", "message": "internal server error"}})
}
