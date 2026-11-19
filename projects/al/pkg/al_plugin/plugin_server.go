package al_plugin

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"google.golang.org/grpc"
)

func grpcRecover(rec any) error {
	if rec == nil {
		return nil
	}
	return fmt.Errorf("plugin panic (%T)\n%s", rec, debug.Stack())
}

type PluginServer struct {
	server  *grpc.Server
	res     chan error
	ctx     *al.CmdCtx
	started bool
}

func streamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = grpcRecover(r)
		}
	}()
	return handler(srv, ss)
}

func unaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = grpcRecover(r)
		}
	}()
	return handler(ctx, req)
}

func NewPluginServer(ctx *al.CmdCtx, plugin al_proto.PluginServiceServer) *PluginServer {
	server := grpc.NewServer(
		grpc.WaitForHandlers(true),
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
	)
	al_proto.RegisterPluginServiceServer(server, plugin)
	res := &PluginServer{
		server: server,
		res:    make(chan error, 1),
		ctx:    ctx,
	}
	return res
}

func (self *PluginServer) Stop(ctx context.Context) error {
	self.ctx.Logger.Printf("stopping plugin server")
	if !self.started {
		return nil
	}
	stopCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	done := make(chan struct{})
	go func() { self.server.GracefulStop(); close(done) }()
	var stopErr error
	select {
	case <-done:
	case <-stopCtx.Done():
		stopErr = fmt.Errorf("plugin server drain deadline: %w", stopCtx.Err())
		self.server.Stop()
		<-done
	}
	if err := <-self.res; err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return errors.Join(stopErr, err)
	}
	return stopErr
}

func (self *PluginServer) Start(_ context.Context) error {
	self.ctx.Logger.Printf("starting plugin server")
	listener, err := NewIOListener(self.ctx.Stdin, self.ctx.Stdout)
	if err != nil {
		return fmt.Errorf("could not create a listener: %w", err)
	}
	listener.onDisconnect = self.ctx.RequestShutdown
	self.started = true
	go func() {
		err := self.server.Serve(listener)
		if err != nil {
			if self.ctx.RequestShutdown != nil {
				self.ctx.RequestShutdown()
			}
			err = fmt.Errorf("could not serve: %w", err)
		}
		self.res <- err
	}()
	return nil
}
