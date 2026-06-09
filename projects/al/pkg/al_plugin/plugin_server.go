package al_plugin

import (
	"context"
	"fmt"
	"io"
	"log"
	"runtime/debug"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
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
	logger *log.Logger
	server *grpc.Server
	stdin  io.Reader
	stdout io.Writer
	res    chan error
}

func NewPluginServer(stdin io.Reader, stdout io.Writer, stderr io.Writer, plugin al_proto.PluginServiceServer) (*PluginServer, error) {
	server := grpc.NewServer(
		grpc.StreamInterceptor(
			func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
				defer func() {
					if r := recover(); r != nil {
						err = grpcRecover(r)
					}
				}()
				return handler(srv, ss)
			},
		),
		grpc.UnaryInterceptor(
			func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
				defer func() {
					if r := recover(); r != nil {
						err = grpcRecover(r)
					}
				}()
				return handler(ctx, req)
			},
		),
	)
	al_proto.RegisterPluginServiceServer(server, plugin)

	return &PluginServer{
		logger: log.New(stderr, "com.alwaldend.src.projects.al.pkg.al_plugin.PluginServer ", log.Flags()),
		server: server,
		stdin:  stdin,
		stdout: stdout,
		res:    make(chan error, 1),
	}, nil
}

func (self *PluginServer) Shutdown(ctx context.Context) error {
	self.logger.Printf("Shutting down the plugin server")
	var wg fp.WaitGroupE
	wg.Go(func() error {
		self.server.GracefulStop()
		return nil
	})
	if err := wg.WaitCtx(ctx); err != nil {
		return fmt.Errorf("shut down with an error: %w", err)
	}
	return nil
}

func (self *PluginServer) Start() error {
	listener, err := NewIOListener(self.stdin, self.stdout)
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
