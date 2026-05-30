package al_plugin

import (
	"context"
	"fmt"
	"io"
	"os/signal"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"google.golang.org/grpc"
)

func recoverError(rec any) error {
	switch recTyped := rec.(type) {
	case string:
		return fmt.Errorf("panicked: %s", recTyped)
	case error:
		return fmt.Errorf("panicked: %w", recTyped)
	default:
		return fmt.Errorf("panicked: %s", recTyped)
	}
}

func Serve(ctx context.Context, stdin io.Reader, stdout io.Writer, plugin al_proto.PluginServiceServer) error {
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	server := grpc.NewServer(
		grpc.StreamInterceptor(
			func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
				defer func() {
					if r := recover(); r != nil {
						err = recoverError(r)
					}
				}()
				return handler(ctx, ss)
			},
		),
		grpc.UnaryInterceptor(
			func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
				defer func() {
					if r := recover(); r != nil {
						err = recoverError(r)
					}
				}()
				return handler(ctx, req)
			},
		),
	)
	go func() {
		for range ctx.Done() {
		}
		server.GracefulStop()
	}()
	al_proto.RegisterPluginServiceServer(server, plugin)
	listener, err := NewIOListener(stdin, stdout)
	if err != nil {
		return fmt.Errorf("could not create a listener: %w", err)
	}
	err = server.Serve(listener)
	if err != nil {
		return fmt.Errorf("could not serve: %w", err)
	}
	return nil
}
