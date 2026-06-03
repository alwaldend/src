package al_plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
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

func parseConfigItem(itemRaw *al_proto.PluginConfigItem) fp.Result[any] {
	switch item := itemRaw.GetValue().(type) {
	case *al_proto.PluginConfigItem_ValueString:
		return fp.Right[any](item.ValueString)
	case *al_proto.PluginConfigItem_ValueBool:
		return fp.Right[any](item.ValueBool)
	case *al_proto.PluginConfigItem_ValueList_:
		res := []any{}
		for _, v := range item.ValueList.Items {
			v2, err := fp.Get(parseConfigItem(v))
			if err != nil {
				return fp.Left[any](err)
			}
			res = append(res, v2)
		}
		return fp.Right[any](res)
	case *al_proto.PluginConfigItem_ValueMap_:
		res := map[string]any{}
		for k, v := range item.ValueMap.Items {
			v2, err := fp.Get(parseConfigItem(v))
			if err != nil {
				return fp.Left[any](err)
			}
			res[k] = v2
		}
		return fp.Right[any](res)
	default:
		return fp.Left[any](fmt.Errorf("invalid item %T: %s", item, itemRaw))
	}
}

func ParseConfig(config *al_proto.PluginConfig, target proto.Message) fp.EmptyEither {
	res := map[string]any{}
	for k, v := range config.Config {
		v2, err := fp.Get(parseConfigItem(v))
		if err != nil {
			return fp.EmptyLeft(err)
		}
		res[k] = v2
	}
	body, err := json.Marshal(res)
	if err != nil {
		return fp.EmptyLeft(fmt.Errorf("could not marshal plugin config: %w", err))
	}
	err = protojson.Unmarshal(body, target)
	if err != nil {
		return fp.EmptyLeft(fmt.Errorf("could not unmarshal plugin config: %w", err))
	}
	return fp.EmptyRight()
}

func ServeDefault[T al_proto.PluginServiceServer](plugin T) fp.EmptyEither {
	return Serve(context.Background(), os.Stdin, os.Stdout, plugin)
}

func Serve[T al_proto.PluginServiceServer](ctx context.Context, stdin io.Reader, stdout io.Writer, plugin T) fp.EmptyEither {
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
	return fp.EmptyRight()
}
