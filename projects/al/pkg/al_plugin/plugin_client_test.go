package al_plugin

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/al"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/lifecycle"
	"github.com/bazelbuild/rules_go/go/runfiles"
)

type clientTestPlugin struct {
	al_proto.UnimplementedPluginServiceServer
	fail bool
}

func (p *clientTestPlugin) PluginStart(context.Context, *al_proto.PluginStartRequest) (*al_proto.PluginStartResponse, error) {
	if p.fail {
		return nil, errors.New("synthetic startup failure")
	}
	return &al_proto.PluginStartResponse{}, nil
}

func TestPluginClientHelper(t *testing.T) {
	mode := os.Getenv("AL_PLUGIN_CLIENT_TEST")
	if mode == "" {
		return
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()
	cmdCtx := al.NewCmdCtx(ctx, "")
	cmdCtx.Logger = log.New(io.Discard, "", 0)
	server := NewPluginServer(cmdCtx, &clientTestPlugin{fail: mode == "start-failure"})
	if err := server.Start(ctx); err != nil {
		os.Exit(30)
	}
	<-cmdCtx.Ctx.Done()
	cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), time.Second)
	defer cleanupCancel()
	if err := server.Stop(cleanupCtx); err != nil {
		os.Exit(31)
	}
	if err := os.WriteFile(os.Getenv("AL_PLUGIN_CLIENT_MARKER"), []byte("cleaned"), 0o600); err != nil {
		os.Exit(32)
	}
	if mode == "exit-failure" {
		os.Exit(23)
	}
	os.Exit(0)
}

func TestPluginClientLifecycle(t *testing.T) {
	run, err := runfiles.New()
	if err != nil {
		t.Fatal(err)
	}
	bin, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"clean", "start-failure", "exit-failure"} {
		t.Run(mode, func(t *testing.T) {
			marker := filepath.Join(t.TempDir(), "cleaned")
			plugin := &al_proto.PluginConfig{Name: "test", Bin: bin, Args: []string{"-test.run=^TestPluginClientHelper$"}, Env: map[string]string{"AL_PLUGIN_CLIENT_TEST": mode, "AL_PLUGIN_CLIENT_MARKER": marker}}
			cmdCtx := al.NewCmdCtx(context.Background(), "")
			defer cmdCtx.RequestShutdown()
			cmdCtx.Logger = log.New(io.Discard, "", 0)
			client, err := NewPluginClient(cmdCtx, run, &al_proto.Config{}, plugin)
			if err != nil {
				t.Fatal(err)
			}
			var manager lifecycle.Manager
			manager.Add(client)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			startErr := manager.Start(ctx)
			if (startErr != nil) != (mode == "start-failure") {
				t.Fatalf("start error %v", startErr)
			}
			stopErr := manager.Stop(ctx)
			if (stopErr != nil) != (mode == "exit-failure") {
				t.Fatalf("stop error %v", stopErr)
			}
			if client.cmd.ProcessState == nil {
				t.Fatal("plugin not reaped")
			}
			if got, err := os.ReadFile(marker); err != nil || string(got) != "cleaned" {
				t.Fatalf("cleanup not completed: %q %v", got, err)
			}
			if err := manager.Stop(ctx); err != nil {
				t.Fatalf("repeat manager stop: %v", err)
			}
		})
	}
}

func TestPluginPanicDoesNotExposeValue(t *testing.T) {
	err := grpcRecover(errors.New("synthetic-secret-value"))
	if err == nil || strings.Contains(err.Error(), "synthetic-secret-value") {
		t.Fatalf("panic was not redacted: %v", err)
	}
}

func TestPluginClientStopBeforeStartClosesPipes(t *testing.T) {
	run, err := runfiles.New()
	if err != nil {
		t.Fatal(err)
	}
	bin, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmdCtx := al.NewCmdCtx(context.Background(), "")
	defer cmdCtx.RequestShutdown()
	client, err := NewPluginClient(cmdCtx, run, &al_proto.Config{}, &al_proto.PluginConfig{Name: "unstarted", Bin: bin})
	if err != nil {
		t.Fatal(err)
	}
	childStdin, ok := client.cmd.Stdin.(*os.File)
	if !ok {
		t.Fatal("missing child stdin file")
	}
	childStdout, ok := client.cmd.Stdout.(*os.File)
	if !ok {
		t.Fatal("missing child stdout file")
	}
	var manager lifecycle.Manager
	manager.Add(client)
	if err := manager.Stop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if !client.ioConn.closed.Load() {
		t.Fatal("unstarted plugin transport remains open")
	}
	for _, endpoint := range []*os.File{childStdin, childStdout} {
		if _, err := endpoint.Stat(); !errors.Is(err, os.ErrClosed) {
			t.Fatalf("child pipe descriptor remains open: %v", err)
		}
	}
	if _, err := os.Stderr.Stat(); err != nil {
		t.Fatalf("shared stderr closed: %v", err)
	}
	if client.cmd.Process != nil {
		t.Fatal("cleanup started plugin")
	}
}

func TestPluginManagerStopsIndependentClientsConcurrently(t *testing.T) {
	cmdCtx := al.NewCmdCtx(context.Background(), "")
	defer cmdCtx.RequestShutdown()
	manager, err := NewManager(cmdCtx, &al_proto.Config{})
	if err != nil {
		t.Fatal(err)
	}
	entered := make(chan struct{}, 2)
	release := make(chan struct{})
	firstErr, secondErr := errors.New("first cleanup failure"), errors.New("second cleanup failure")
	for _, failure := range []error{firstErr, secondErr} {
		manager.Lifecycle().Add(lifecycle.StoppableFunc(func(ctx context.Context) error {
			entered <- struct{}{}
			select {
			case <-release:
			case <-ctx.Done():
				return ctx.Err()
			}
			return failure
		}))
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stopped := make(chan error, 1)
	go func() { stopped <- manager.Lifecycle().Stop(ctx) }()
	for range 2 {
		select {
		case <-entered:
		case <-ctx.Done():
			t.Fatal("one plugin prevented sibling shutdown")
		}
	}
	select {
	case <-stopped:
		t.Fatal("stop abandoned cleanup callbacks")
	default:
	}
	close(release)
	select {
	case err := <-stopped:
		if !errors.Is(err, firstErr) || !errors.Is(err, secondErr) {
			t.Fatalf("lost cleanup errors: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("shutdown did not finish")
	}
	if err := manager.Lifecycle().Stop(ctx); err != nil {
		t.Fatalf("cleanup repeated: %v", err)
	}
}
