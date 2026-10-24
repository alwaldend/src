//go:build cgo && linux && amd64

package ankiclient_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// startupBarrierContext pauses a context check to exercise overlapping native lifecycle calls.
type startupBarrierContext struct {
	// Context provides the real test cancellation and deadline.
	context.Context
	// entered reports that startup reached the controlled scheduling point.
	entered chan struct{}
	// release resumes startup after the competing Stop has finished.
	release chan struct{}
	// once pauses only the first check, allowing later context use.
	once sync.Once
}

// Err pauses the first check and preserves context.Context's canonical cancellation sentinel.
func (ctx *startupBarrierContext) Err() error {
	ctx.once.Do(func() {
		close(ctx.entered)
		<-ctx.release
	})
	return ctx.Context.Err()
}

// TestCallerOwnedCollections reuses one transport while the caller opens and closes native storage.
func TestCallerOwnedCollections(t *testing.T) {
	registerWorkerFactories(t)
	root := t.TempDir()
	client, err := ankigrpc.NewClient(ankigrpc.SubprocessAddress(workerPath(t), "serve"))
	if err != nil {
		t.Fatal(err)
	}
	var collections collection.BackendCollectionServiceClient
	opened := false
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if opened {
			if _, err := collections.CloseCollection(cleanup, &collection.CloseCollectionRequest{}); err != nil {
				t.Error(err)
			}
		}
		if err := client.Stop(cleanup); err != nil {
			t.Error(err)
		}
	})
	if err := client.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	connection, err := client.Connection()
	if err != nil {
		t.Fatal(err)
	}
	collections = collection.NewBackendCollectionServiceClient(connection)
	for _, name := range []string{"first", "second"} {
		directory := filepath.Join(root, name)
		mediaDirectory := filepath.Join(directory, "collection.media")
		if err := os.MkdirAll(mediaDirectory, 0o700); err != nil {
			t.Fatal(err)
		}
		if _, err := collections.OpenCollection(t.Context(), &collection.OpenCollectionRequest{
			CollectionPath: filepath.Join(directory, "collection.anki2"), MediaFolderPath: mediaDirectory, MediaDbPath: filepath.Join(directory, "collection.media.db2"),
		}); err != nil {
			t.Fatal(err)
		}
		opened = true
		if _, err := notetypes.NewNotetypesServiceClient(connection).GetNotetypeIdByName(t.Context(), &generic.String{Val: "Basic"}); err != nil {
			t.Fatal(err)
		}
		if _, err := collections.CloseCollection(t.Context(), &collection.CloseCollectionRequest{}); err != nil {
			t.Fatal(err)
		}
		opened = false
		if _, err := collections.LatestProgress(t.Context(), &generic.Empty{}); err != nil {
			t.Fatalf("closing caller-owned storage stopped its transport: %v", err)
		}
	}
	if err := client.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Connection(); err == nil {
		t.Fatal("stopped transport retained connection access")
	}
	writeReport(t, "caller-owned-collections.json", map[string]any{"collections": 2, "transport_reused": true, "worker_reaped": true})
}

// TestCompletedMediaWorkerShutdown verifies successful media work permits graceful stdio worker exit.
func TestCompletedMediaWorkerShutdown(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	root := t.TempDir()
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}})
	var collections collection.BackendCollectionServiceClient
	opened := false
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if opened {
			if _, err := collections.CloseCollection(cleanup, &collection.CloseCollectionRequest{}); err != nil {
				t.Error(err)
			}
		}
		if err := worker.Stop(cleanup); err != nil {
			t.Error(err)
		}
	})
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	mediaDirectory := filepath.Join(root, "collection.media")
	if err := os.Mkdir(mediaDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	collections = collection.NewBackendCollectionServiceClient(connection)
	if _, err := collections.OpenCollection(ctx, &collection.OpenCollectionRequest{
		CollectionPath: filepath.Join(root, "collection.anki2"), MediaFolderPath: mediaDirectory, MediaDbPath: filepath.Join(root, "collection.media.db2"),
	}); err != nil {
		t.Fatal(err)
	}
	opened = true
	const content = "completed native media download\n"
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "text/plain")
		if _, err := response.Write([]byte(content)); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	downloaded, err := media.NewBackendMediaServiceClient(connection).AddMediaFromUrl(ctx, &media.AddMediaFromUrlRequest{Url: server.URL + "/completed.txt"})
	if err != nil {
		t.Fatal(err)
	}
	if downloaded.GetError() != "" || downloaded.GetFilename() == "" {
		t.Fatalf("media download did not complete successfully: %v", downloaded)
	}
	data, err := os.ReadFile(filepath.Join(mediaDirectory, downloaded.GetFilename()))
	if err != nil || string(data) != content {
		t.Fatalf("downloaded media mismatch: %q, %v", data, err)
	}
	if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	opened = false
	if err := worker.Stop(ctx); err != nil {
		t.Fatalf("stop worker after completed media: %v", err)
	}
	select {
	case <-worker.Done():
	default:
		t.Fatal("successful worker shutdown did not report native completion")
	}
	writeReport(t, "completed-media-shutdown.json", map[string]any{"filename": downloaded.GetFilename(), "content": string(data), "collection_closed": true, "worker_exit_success": true, "native_completion": true})
}
