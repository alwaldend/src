package al_plugin

import (
	"context"
	"fmt"
	"io"
	"runtime/debug"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"google.golang.org/grpc"
)

func RecoverError(rec any) error {
	if rec == nil {
		return nil
	}
	return fmt.Errorf("panic: %s\n%s", rec, string(debug.Stack()))
}

func Serve[T al_proto.PluginServiceServer](ctx context.Context, stdin io.Reader, stdout io.Writer, plugin T) fp.EmptyEither {
	server := grpc.NewServer(
		grpc.StreamInterceptor(
			func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
				defer func() {
					if r := recover(); r != nil {
						err = RecoverError(r)
					}
				}()
				return handler(ctx, ss)
			},
		),
		grpc.UnaryInterceptor(
			func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
				defer func() {
					if r := recover(); r != nil {
						err = RecoverError(r)
					}
				}()
				return handler(ctx, req)
			},
		),
	)
	al_proto.RegisterPluginServiceServer(server, plugin)
	listener, err := NewIOListener(stdin, stdout)
	if err != nil {
		return fp.EmptyLeft(fmt.Errorf("could not create a listener: %w", err))
	}
	err = server.Serve(listener)
	if err != nil {
		return fp.EmptyLeft(fmt.Errorf("could not serve: %w", err))
	}
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	return fp.EmptyRight()
}
