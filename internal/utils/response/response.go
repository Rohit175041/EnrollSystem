package response

import (
	"encoding/json"
	"net/http"
)

// Response is the base response structure
type Response struct {
	Status  string `json:"status"`            // "success", "error", "fail"
	Message string `json:"message,omitempty"` // optional message
	Error   string `json:"error,omitempty"`   // error message if any
}

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusError   = "error"
)

// WriteJSON writes JSON response with given HTTP status code
func WriteJSON(w http.ResponseWriter, status int, resp Response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(resp)
}

// Success returns a success response with optional data and message
func Success(message string) Response {
	return Response{
		Status:  StatusSuccess,
		Message: message,
	}
}

// Fail returns a failure response with message (e.g. validation errors)
func Fail(message string) Response {
	return Response{
		Status:  StatusFail,
		Message: message,
	}
}

// Error returns an error response with error message
func Error(err error) Response {
	return Response{
		Status: StatusError,
		Error:  err.Error(),
	}
}
