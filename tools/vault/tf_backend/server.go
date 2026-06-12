package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
)

type Server struct {
	listener net.Listener
	ctx      *al.CmdCtx
	handler  http.HandlerFunc
	server   *http.Server
	serveErr chan error
}

func NewServer(ctx *al.CmdCtx, handler http.HandlerFunc) *Server {
	return &Server{ctx: ctx, handler: handler, serveErr: make(chan error, 1)}
}

func (self *Server) Stop(ctx context.Context) error {
	if self.server == nil {
		return fmt.Errorf("missing server")
	}
	var resErr error
	if err := self.server.Shutdown(ctx); err != nil {
		resErr = fmt.Errorf("server shutdown error: %w", err)
	}
	if err := <-self.serveErr; err != nil {
		resErr = errors.Join(resErr, fmt.Errorf("serve error: %w", err))
	}
	return resErr
}

func (self *Server) Start(_ context.Context) error {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("could not listen on a port: %w", err)
	}
	self.listener = listener
	handler := http.NewServeMux()
	handler.HandleFunc("/", self.handler)
	self.server = &http.Server{Handler: handler}
	go func() {
		self.serveErr <- self.server.Serve(listener)
	}()
	return nil
}

func (self *Server) Address() (string, error) {
	if self.listener == nil {
		return "", fmt.Errorf("listener is not set")
	}
	return self.listener.Addr().String(), nil
}
