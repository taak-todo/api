package server

import (
	"net/http"

	"github.com/google/uuid"
)

// RequestID retrieves the request id from the context.
func RequestID(req *http.Request) uuid.UUID {
	requestID, ok := req.Context().Value(requestIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return requestID
}
