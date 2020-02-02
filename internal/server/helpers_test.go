package server

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRequestID(t *testing.T) {
	t.Run("returns empty id if no id exists", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "example.com", nil)
		require.Nil(t, err)

		requestID := RequestID(req)
		require.Equal(t, requestID, uuid.Nil)
	})

	t.Run("returns empty id if invalid id exists", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "example.com", nil)
		require.Nil(t, err)

		req = req.WithContext(context.WithValue(req.Context(), requestIDKey, "foobar"))
		requestID := RequestID(req)

		require.Equal(t, requestID, uuid.Nil)
	})

	t.Run("returns id if one exists", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "example.com", nil)
		require.Nil(t, err)

		expected, err := uuid.NewRandom()
		require.Nil(t, err)

		req = req.WithContext(context.WithValue(req.Context(), requestIDKey, expected))
		requestID := RequestID(req)

		require.Equal(t, requestID, expected)
	})
}
