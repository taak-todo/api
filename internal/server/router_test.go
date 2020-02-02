package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	t.Run("responds with a 405 for invalid methods for routes", func(t *testing.T) {
		handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte("test"))
		})

		route := Route{
			Method:  http.MethodPost,
			Path:    "/test",
			Handler: handler,
		}

		server := httptest.NewServer(NewRouter(route))
		defer server.Close()

		res, err := http.Get(server.URL + "/test")
		require.Nil(t, err)
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusMethodNotAllowed)
	})

	t.Run("responds with a 404 for invalid routes", func(t *testing.T) {
		server := httptest.NewServer(NewRouter())
		defer server.Close()

		res, err := http.Get(server.URL + "/test")
		require.Nil(t, err)
		defer res.Body.Close()

		require.Equal(t, res.StatusCode, http.StatusNotFound)
	})

	t.Run("handles the given routes", func(t *testing.T) {
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusAccepted)

				next.ServeHTTP(rw, req)
			})
		}

		handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte("test"))
		})

		route := Route{
			Method: http.MethodGet,
			Path:   "/test",
			Middlewares: []func(http.Handler) http.Handler{
				middleware,
			},
			Handler: handler,
		}

		server := httptest.NewServer(NewRouter(route))
		defer server.Close()

		res, err := http.Get(server.URL + "/test")
		require.Nil(t, err)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		require.Nil(t, err)

		require.Equal(t, res.StatusCode, http.StatusAccepted)
		require.Equal(t, string(body), "test")
	})

	t.Run("for invalid route methods it defaults to GET", func(t *testing.T) {
		handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte("test"))
		})

		route := Route{
			Method:  "FOO",
			Path:    "/test",
			Handler: handler,
		}

		server := httptest.NewServer(NewRouter(route))
		defer server.Close()

		res, err := http.Get(server.URL + "/test")
		require.Nil(t, err)
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		require.Nil(t, err)

		require.Equal(t, string(body), "test")
	})
}
