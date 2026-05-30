package al_plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/signal"
	"runtime/debug"
	"syscall"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func RecoverError(rec any) error {
	if rec == nil {
		return nil
	}
	return fmt.Errorf("panic: %s\n%s", rec, string(debug.Stack()))
}

func parseConfigItem(itemRaw *al_proto.PluginConfigItem) fp.Monad[any] {
	switch item := itemRaw.GetValue().(type) {
	case *al_proto.PluginConfigItem_ValueString:
		return fp.Right[any](item.ValueString)
	case *al_proto.PluginConfigItem_ValueBool:
		return fp.Right[any](item.ValueBool)
	case *al_proto.PluginConfigItem_ValueList_:
		return fp.Pipe3(
			fp.Just(item.ValueList.Items),
			fp.CollectSlice(parseConfigItem),
			fp.ToAny,
		)
	case *al_proto.PluginConfigItem_ValueMap_:
		return fp.Pipe3(
			fp.Just(item.ValueMap.Items),
			fp.CollectMap(func(k string, v *al_proto.PluginConfigItem) fp.Monad[any] { return parseConfigItem(v) }),
			fp.ToAny,
		)
	default:
		return fp.Left[any](fmt.Errorf("invalid item %T: %s", item, itemRaw))
	}
}

func ParseConfig(config *al_proto.PluginConfig, target proto.Message) fp.EmptyMonad {
	return fp.Pipe5(
		fp.Just(config.Config),
		fp.CollectMapV[string](parseConfigItem),
		fp.ToAny,
		fp.EitherFunc2(json.Marshal),
		fp.EitherFunc(func(b []byte) error { return protojson.Unmarshal(b, target) }),
	)
}

func Serve(ctx context.Context, stdin io.Reader, stdout io.Writer, plugin al_proto.PluginServiceServer) fp.EmptyMonad {
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
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
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	al_proto.RegisterPluginServiceServer(server, plugin)
	listener, err := NewIOListener(stdin, stdout)
	if err != nil {
		return fp.EmptyLeft(fmt.Errorf("could not create a listener: %w", err))
	}
	err = server.Serve(listener)
	if err != nil {
		return fp.EmptyLeft(fmt.Errorf("could not serve: %w", err))
	}
	return nil
}
