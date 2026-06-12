package al_plugin

import (
	"context"
	"fmt"
	"runtime/debug"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"google.golang.org/grpc"
)

func grpcRecover(rec any) error {
	var err error
	switch recTyped := rec.(type) {
	case error:
		err = recTyped
	case nil:
		return nil
	}
	return fmt.Errorf("panic: %s\n%s: %w", rec, string(debug.Stack()), err)
}

type PluginServer struct {
	server *grpc.Server
	res    chan error
	ctx    *al.CmdCtx
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
	self.server.GracefulStop()
	if err := <-self.res; err != nil {
		return fmt.Errorf("serving error error: %w", err)
	}
	return nil
}

func (self *PluginServer) Start(_ context.Context) error {
	self.ctx.Logger.Printf("starting plugin server")
	listener, err := NewIOListener(self.ctx.Stdin, self.ctx.Stdout)
	if err != nil {
		return fmt.Errorf("could not create a listener: %w", err)
	}
	go func() {
		err := self.server.Serve(listener)
		if err != nil {
			err = fmt.Errorf("could not serve: %w", err)
		}
		self.res <- err
	}()
	return nil
}
