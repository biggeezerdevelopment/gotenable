// Package base provides the foundational types and utilities for the goTenable SDK.
package base

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the Tenable API.
type APIError struct {
	StatusCode int
	Message    string
	RequestID  string
	Response   []byte
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("tenable API error (status %d, request_id: %s): %s", e.StatusCode, e.RequestID, e.Message)
	}
	return fmt.Sprintf("tenable API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error is a 404 Not Found error.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

// IsUnauthorized returns true if the error is a 401 Unauthorized error.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized
}

// IsForbidden returns true if the error is a 403 Forbidden error.
func (e *APIError) IsForbidden() bool {
	return e.StatusCode == http.StatusForbidden
}

// IsRateLimited returns true if the error is a 429 Too Many Requests error.
func (e *APIError) IsRateLimited() bool {
	return e.StatusCode == http.StatusTooManyRequests
}

// IsServerError returns true if the error is a 5xx server error.
func (e *APIError) IsServerError() bool {
	return e.StatusCode >= 500 && e.StatusCode < 600
}

// AuthenticationError represents an authentication failure.
type AuthenticationError struct {
	Message string
}

// Error implements the error interface.
func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("authentication error: %s", e.Message)
}

// ConnectionError represents a connection failure.
type ConnectionError struct {
	URL     string
	Message string
	Err     error
}

// Error implements the error interface.
func (e *ConnectionError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("connection error to %s: %s: %v", e.URL, e.Message, e.Err)
	}
	return fmt.Sprintf("connection error to %s: %s", e.URL, e.Message)
}

// Unwrap returns the underlying error.
func (e *ConnectionError) Unwrap() error {
	return e.Err
}

// ValidationError represents a validation failure for input parameters.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

// FileDownloadError represents a failure to download a file.
type FileDownloadError struct {
	Resource   string
	ResourceID string
	Filename   string
	Err        error
}

// Error implements the error interface.
func (e *FileDownloadError) Error() string {
	return fmt.Sprintf("file download error: resource %s:%s requested file %s and failed: %v",
		e.Resource, e.ResourceID, e.Filename, e.Err)
}

// Unwrap returns the underlying error.
func (e *FileDownloadError) Unwrap() error {
	return e.Err
}

// ExportError represents an error during export operations.
type ExportError struct {
	ExportType string
	UUID       string
	Message    string
}

// Error implements the error interface.
func (e *ExportError) Error() string {
	return fmt.Sprintf("%s export %s error: %s", e.ExportType, e.UUID, e.Message)
}

// ExportTimeoutError represents a timeout during export operations.
type ExportTimeoutError struct {
	ExportType string
	UUID       string
}

// Error implements the error interface.
func (e *ExportTimeoutError) Error() string {
	return fmt.Sprintf("%s export %s has timed out", e.ExportType, e.UUID)
}
