package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// httpBackendServer owns one loopback listener for a Terraform backend.
type httpBackendServer struct {
	// listener provides the address exported to Terraform after startup.
	listener net.Listener
	// handler handles state and lock requests on the loopback listener.
	handler http.HandlerFunc
	// server owns active connections after startup.
	server *http.Server
	// serveErr receives the sole Serve result without delaying its goroutine.
	serveErr chan error
	// stopOnce serializes shutdown and consumption of the Serve result.
	stopOnce sync.Once
	// stopErr retains the result for repeated shutdown calls.
	stopErr error
}

// newServer prepares an HTTP server without opening a listener.
func newServer(handler http.HandlerFunc) *httpBackendServer {
	return &httpBackendServer{handler: handler, serveErr: make(chan error, 1)}
}

// stop drains requests once and returns the same outcome to subsequent callers.
func (s *httpBackendServer) stop(ctx context.Context) error {
	s.stopOnce.Do(func() { s.stopErr = s.shutdown(ctx) })
	if s.stopErr != nil {
		return fmt.Errorf("stop backend HTTP server: %w", s.stopErr)
	}
	return nil
}

// shutdown forces remaining connections closed when graceful shutdown times out.
func (s *httpBackendServer) shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	var resErr error
	if err := s.server.Shutdown(ctx); err != nil {
		resErr = fmt.Errorf("drain HTTP requests: %w", err)
		if err := s.server.Close(); err != nil {
			resErr = errors.Join(resErr, fmt.Errorf("force HTTP connections closed: %w", err))
		}
	}
	if err := <-s.serveErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		resErr = errors.Join(resErr, fmt.Errorf("serve error: %w", err))
	}
	return resErr
}

// start binds an ephemeral loopback port before launching the HTTP server.
func (s *httpBackendServer) start(ctx context.Context) error {
	var config net.ListenConfig
	listener, err := config.Listen(ctx, "tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("could not listen on a port: %w", err)
	}
	s.listener = listener
	handler := http.NewServeMux()
	handler.HandleFunc("/", s.handler)
	s.server = &http.Server{Handler: handler, ReadHeaderTimeout: 10 * time.Second}
	go func() {
		s.serveErr <- s.server.Serve(listener)
	}()
	return nil
}

// address returns the bound endpoint after successful startup.
func (s *httpBackendServer) address() (string, error) {
	if s.listener == nil {
		return "", fmt.Errorf("listener is not set")
	}
	return s.listener.Addr().String(), nil
}
