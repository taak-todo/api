package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestRequestIDMiddleware(t *testing.T) {
	t.Run("sets id on request context", func(t *testing.T) {
		handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			val := req.Context().Value(requestIDKey)

			require.IsType(t, val, uuid.UUID{})
			require.NotNil(t, val)
		})

		server := httptest.NewServer(RequestIDMiddleware(handler))
		defer server.Close()

		res, err := http.Get(server.URL)
		require.Nil(t, err)
		defer res.Body.Close()
	})
}

func TestRecovererMiddleware(t *testing.T) {
	t.Run("responds with a 500 if a panic occurs", func(t *testing.T) {
		handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			panic(io.EOF)
		})

		server := httptest.NewServer(RecovererMiddleware(handler))
		defer server.Close()

		res, err := http.Get(server.URL)
		require.Nil(t, err)
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusInternalServerError)
	})
}
