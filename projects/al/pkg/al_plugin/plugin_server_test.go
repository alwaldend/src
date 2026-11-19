package al_plugin

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type blockingPlugin struct {
	al_proto.UnimplementedPluginServiceServer
	entered chan struct{}
	exited  chan struct{}
}

func (p *blockingPlugin) PluginStart(ctx context.Context, _ *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	close(p.entered)
	<-ctx.Done()
	close(p.exited)
	return nil, ctx.Err()
}

func TestServerShutdownCancelsAndDrainsHandler(t *testing.T) {
	serverPipe, clientPipe := net.Pipe()
	defer clientPipe.Close()
	plugin := &blockingPlugin{entered: make(chan struct{}), exited: make(chan struct{})}
	server := NewPluginServer(&al.CmdCtx{Stdin: serverPipe, Stdout: serverPipe, Logger: log.New(io.Discard, "", 0)}, plugin)
	if err := server.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	conn, err := grpc.NewClient("passthrough:///test", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return clientPipe, nil }))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	callDone := make(chan error, 1)
	go func() {
		_, err := al_proto.NewPluginServiceClient(conn).PluginStart(context.Background(), &al_proto.PluginStartRequest{})
		callDone <- err
	}()
	select {
	case <-plugin.entered:
	case <-time.After(time.Second):
		t.Fatal("handler did not start")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := server.Stop(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("stop error %v", err)
	}
	select {
	case <-plugin.exited:
	default:
		t.Fatal("shutdown returned before handler cleanup")
	}
	select {
	case <-callDone:
	case <-time.After(time.Second):
		t.Fatal("call did not stop")
	}
}

func TestParentDisconnectRequestsShutdown(t *testing.T) {
	serverPipe, clientPipe := net.Pipe()
	called := make(chan struct{}, 1)
	server := NewPluginServer(&al.CmdCtx{Stdin: serverPipe, Stdout: serverPipe, Logger: log.New(io.Discard, "", 0), RequestShutdown: func() {
		select {
		case called <- struct{}{}:
		default:
		}
	}}, &blockingPlugin{})
	if err := server.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	clientPipe.Close()
	select {
	case <-called:
	case <-time.After(time.Second):
		t.Fatal("disconnect did not request shutdown")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.Stop(ctx); err != nil {
		t.Fatal(err)
	}
}
