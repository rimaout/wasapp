package api

import (
	"encoding/json"
	"net/http"
)

// errorResponse defines a structural type for consistent error responses
type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// respondWithError sets the JSON headers, writes the HTTP status,
// and marshals the custom string code and error message.
func (rt *_router) respondWithError(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// We ignore the error here because encoding a basic hardcoded struct
	// to a client response stream will practically never fail.
	_ = json.NewEncoder(w).Encode(errorResponse{
		Code:    code,
		Message: message,
	})
}
