package wasm

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	wasmadapter "git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/wasm"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/config"
	"github.com/ankitects/anki/proto/anki/generic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TestInprocessLifecycle checks the production guest's startup and terminal cleanup contract.
func TestInprocessLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	client := wasmadapter.New()
	if _, err := client.Connection(); err == nil {
		t.Fatal("connection available before Start")
	}
	startup, finishStartup := context.WithCancel(ctx)
	if err := client.Start(startup); err != nil {
		t.Fatalf("start in-process guest: %v", err)
	}
	finishStartup()
	t.Cleanup(func() { _ = client.Stop(context.Background()) })
	connection, err := client.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if err := connection.Invoke(ctx, "/unknown.Service/Method", &generic.Empty{}, &generic.Empty{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("startup cancellation stopped the ready guest or lost unknown-method status: %v", err)
	}
	var trailers metadata.MD
	_, err = config.NewConfigServiceClient(connection).GetConfigJson(ctx, &generic.String{Val: "absent"}, grpc.Trailer(&trailers))
	if err == nil || len(trailers.Get("anki-error-bin")) != 1 {
		t.Fatalf("native failure details were lost: %v; %v", err, trailers)
	}
	if err := client.Start(ctx); err == nil {
		t.Fatal("second Start succeeded")
	}
	stop, finishStop := context.WithTimeout(ctx, 10*time.Second)
	defer finishStop()
	if err := client.Stop(stop); err != nil {
		t.Fatalf("stop in-process guest: %v", err)
	}
	select {
	case <-client.Done():
	default:
		t.Fatal("Stop returned before final cleanup")
	}
	if err := client.Stop(ctx); err != nil {
		t.Fatalf("repeated Stop: %v", err)
	}
	if _, err := client.Connection(); err == nil {
		t.Fatal("connection available after Stop")
	}
	if err := connection.Invoke(ctx, "/unknown.Service/Method", &generic.Empty{}, &generic.Empty{}); err == nil {
		t.Fatal("retained connection accepted a call after Stop")
	}
}

// TestInprocessStartupFailures checks cancellation, early Stop and runtime setup failures.
func TestInprocessStartupFailures(t *testing.T) {
	t.Run("cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		client := wasmadapter.New()
		if err := client.Start(ctx); !errors.Is(err, context.Canceled) {
			t.Fatalf("cancelled startup: %v", err)
		}
		if err := client.Stop(context.Background()); err != nil {
			t.Fatalf("cleanup cancelled startup: %v", err)
		}
	})
	t.Run("stopped", func(t *testing.T) {
		client := wasmadapter.New()
		if err := client.Stop(t.Context()); err != nil {
			t.Fatal(err)
		}
		if err := client.Start(t.Context()); err == nil {
			t.Fatal("Start succeeded after Stop")
		}
	})
	t.Run("during-startup", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
		defer cancel()
		client := wasmadapter.New()
		started := make(chan error, 1)
		go func() { started <- client.Start(ctx) }()
		_ = client.Stop(ctx)
		select {
		case <-started:
		case <-ctx.Done():
			t.Fatal("Start did not finish after Stop")
		}
		select {
		case <-client.Done():
		case <-ctx.Done():
			t.Fatal("Stop did not release the starting instance")
		}
	})
	t.Run("cache-failure", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "not-a-directory")
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("ANKI_WASM_CACHE_DIR", path)
		client := wasmadapter.New()
		ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
		defer cancel()
		if err := client.Start(ctx); err == nil || !strings.Contains(err.Error(), "cache") {
			t.Fatalf("cache failure was not retained at startup: %v", err)
		}
		_ = client.Stop(ctx)
		select {
		case <-client.Done():
		case <-ctx.Done():
			t.Fatal("failed instance was not released")
		}
	})
}

// TestInprocessIndependentCollections checks fresh state, large RPCs and independent shutdown.
func TestInprocessIndependentCollections(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Minute)
	defer cancel()
	clients := []*wasmadapter.Client{wasmadapter.New(), wasmadapter.New()}
	connections := make([]grpc.ClientConnInterface, len(clients))
	paths := make([]string, len(clients))
	value := []byte(`"` + strings.Repeat("in-process WASM ", 1100000) + `"`)
	for index, client := range clients {
		if err := client.Start(ctx); err != nil {
			t.Fatalf("start instance %d: %v", index, err)
		}
		t.Cleanup(func() { _ = client.Stop(context.Background()) })
		connection, err := client.Connection()
		if err != nil {
			t.Fatal(err)
		}
		connections[index] = connection
		directory := t.TempDir()
		paths[index] = filepath.Join(directory, "collection.anki2")
		if _, err := collection.NewBackendCollectionServiceClient(connection).OpenCollection(ctx, &collection.OpenCollectionRequest{
			CollectionPath: paths[index], MediaFolderPath: filepath.Join(directory, "media"), MediaDbPath: filepath.Join(directory, "media.db2"),
		}); err != nil {
			t.Fatalf("open independent collection: %v", err)
		}
		data := value
		if index == 1 {
			data = []byte(`"second instance"`)
		}
		if _, err := config.NewConfigServiceClient(connection).SetConfigJson(ctx, &config.SetConfigJsonRequest{Key: "instance", ValueJson: data}); err != nil {
			t.Fatalf("write instance config: %v", err)
		}
	}
	for index, client := range clients {
		got, err := config.NewConfigServiceClient(connections[index]).GetConfigJson(ctx, &generic.String{Val: "instance"})
		expected := string(value)
		if index == 1 {
			expected = `"second instance"`
		}
		if err != nil || string(got.GetJson()) != expected {
			t.Fatalf("instance %d state changed or large response was truncated: %v", index, err)
		}
		if _, err := collection.NewBackendCollectionServiceClient(connections[index]).CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
			t.Fatal(err)
		}
		if err := client.Stop(ctx); err != nil {
			t.Fatal(err)
		}
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := os.ReadFile(paths[0])
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "inprocess-large.anki2"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(output, "inprocess-verification.txt"), []byte("independent production guests; >10 MiB protobuf round trip; final cleanup passed\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
