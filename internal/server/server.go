package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Addr string
}

// Server implements a graceful HTTP server.
type Server struct {
	server *http.Server
	done   chan error
}

func New(handler http.Handler, config Config) *Server {
	return &Server{
		server: &http.Server{
			Addr:         config.Addr,
			Handler:      handler,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 2 * time.Second,
			IdleTimeout:  2 * time.Second,
		},
		done: make(chan error, 1),
	}
}

// Serve starts the HTTP server and serves requests until SIGINT or SIGTERM.
func (s *Server) Serve() error {
	go s.shutdownSignal()

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return <-s.done
}

func (s *Server) shutdownSignal() {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, os.Interrupt, syscall.SIGTERM)

	<-notifier

	ctx, cancel := context.WithTimeout(context.Background(), s.server.WriteTimeout)
	defer cancel()

	s.done <- s.server.Shutdown(ctx)
}
