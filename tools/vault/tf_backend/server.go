package main

import (
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"net"
	"net/http"
	"time"
)

type Server struct {
	listener net.Listener
	wg       *fp.WaitGroupE
	ctx      context.Context
	handler  http.HandlerFunc
}

func NewServer(ctx context.Context, handler http.HandlerFunc, wg *fp.WaitGroupE) *Server {
	return &Server{ctx: ctx, handler: handler, wg: wg}
}

func (self *Server) Start() error {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("could not listen on a port: %w", err)
	}
	self.listener = listener
	handler := http.NewServeMux()
	handler.HandleFunc("/", self.handler)
	server := &http.Server{Handler: handler}
	self.wg.Go(func() error {
		<-self.ctx.Done()
		logger.Printf("Shutting down server %s", listener.Addr().String())
		shutdownCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
		defer cancel()
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			err = fmt.Errorf("could not shutdown: %w", err)
		}
		return err
	})
	self.wg.Go(func() error {
		err := server.Serve(listener)
		if err != nil {
			err = fmt.Errorf("could not serve: %w", err)
		}
		logger.Printf("Finished serving %s", listener.Addr().String())
		return err
	})
	return nil
}

func (self *Server) Address() (string, error) {
	if self.listener == nil {
		return "", fmt.Errorf("listener is not set")
	}
	return self.listener.Addr().String(), nil
}
