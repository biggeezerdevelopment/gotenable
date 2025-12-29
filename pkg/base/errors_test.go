package base

import (
	"errors"
	"net/http"
	"testing"
)

func TestAPIError(t *testing.T) {
	tests := []struct {
		name       string
		err        *APIError
		wantMsg    string
		isNotFound bool
		isUnauth   bool
		isForbid   bool
		isRateL    bool
		isServer   bool
	}{
		{
			name: "not found error",
			err: &APIError{
				StatusCode: http.StatusNotFound,
				Message:    "resource not found",
			},
			wantMsg:    "tenable API error (status 404): resource not found",
			isNotFound: true,
		},
		{
			name: "unauthorized error",
			err: &APIError{
				StatusCode: http.StatusUnauthorized,
				Message:    "invalid credentials",
			},
			wantMsg:  "tenable API error (status 401): invalid credentials",
			isUnauth: true,
		},
		{
			name: "forbidden error",
			err: &APIError{
				StatusCode: http.StatusForbidden,
				Message:    "access denied",
			},
			wantMsg:  "tenable API error (status 403): access denied",
			isForbid: true,
		},
		{
			name: "rate limited error",
			err: &APIError{
				StatusCode: http.StatusTooManyRequests,
				Message:    "too many requests",
			},
			wantMsg: "tenable API error (status 429): too many requests",
			isRateL: true,
		},
		{
			name: "server error",
			err: &APIError{
				StatusCode: http.StatusInternalServerError,
				Message:    "internal error",
			},
			wantMsg:  "tenable API error (status 500): internal error",
			isServer: true,
		},
		{
			name: "error with request ID",
			err: &APIError{
				StatusCode: http.StatusBadRequest,
				Message:    "bad request",
				RequestID:  "req-123",
			},
			wantMsg: "tenable API error (status 400, request_id: req-123): bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.wantMsg {
				t.Errorf("Error() = %v, want %v", got, tt.wantMsg)
			}
			if got := tt.err.IsNotFound(); got != tt.isNotFound {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.isNotFound)
			}
			if got := tt.err.IsUnauthorized(); got != tt.isUnauth {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.isUnauth)
			}
			if got := tt.err.IsForbidden(); got != tt.isForbid {
				t.Errorf("IsForbidden() = %v, want %v", got, tt.isForbid)
			}
			if got := tt.err.IsRateLimited(); got != tt.isRateL {
				t.Errorf("IsRateLimited() = %v, want %v", got, tt.isRateL)
			}
			if got := tt.err.IsServerError(); got != tt.isServer {
				t.Errorf("IsServerError() = %v, want %v", got, tt.isServer)
			}
		})
	}
}

func TestConnectionError(t *testing.T) {
	underlying := errors.New("connection refused")
	err := &ConnectionError{
		URL:     "https://example.com",
		Message: "failed to connect",
		Err:     underlying,
	}

	want := "connection error to https://example.com: failed to connect: connection refused"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %v, want %v", got, want)
	}

	if unwrapped := err.Unwrap(); unwrapped != underlying {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, underlying)
	}
}

func TestValidationError(t *testing.T) {
	err := &ValidationError{
		Field:   "scan_id",
		Message: "must be a positive integer",
	}

	want := "validation error for scan_id: must be a positive integer"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %v, want %v", got, want)
	}
}

func TestExportError(t *testing.T) {
	err := &ExportError{
		ExportType: "assets",
		UUID:       "abc-123",
		Message:    "export failed",
	}

	want := "assets export abc-123 error: export failed"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %v, want %v", got, want)
	}
}

func TestExportTimeoutError(t *testing.T) {
	err := &ExportTimeoutError{
		ExportType: "vulns",
		UUID:       "def-456",
	}

	want := "vulns export def-456 has timed out"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %v, want %v", got, want)
	}
}

