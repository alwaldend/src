//go:build cgo && linux && amd64

package ankiclient_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/native"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// TestDirectBackendOwnership verifies independent named instances and generated-client options.
func TestDirectBackendOwnership(t *testing.T) {
	if _, err := ankigrpc.NewClient("inprocess+unavailable://missing"); err == nil {
		t.Fatal("in-process selection accepted an unavailable engine")
	}
	first, err := ankigrpc.NewClient("inprocess://same-name")
	if err != nil {
		t.Fatal(err)
	}
	second, err := ankigrpc.NewClient("inprocess://same-name")
	if err != nil {
		t.Fatal(err)
	}
	for index, backend := range []ankiclient.Client{first, second} {
		if err := backend.Start(t.Context()); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if err := backend.Stop(context.Background()); err != nil {
				t.Error(err)
			}
		})
		connection, err := backend.Connection()
		if err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(t.TempDir(), "independent.anki2")
		collections := collection.NewBackendCollectionServiceClient(connection)
		if _, err := collections.OpenCollection(t.Context(), &collection.OpenCollectionRequest{
			CollectionPath: path, MediaFolderPath: path + ".media", MediaDbPath: path + ".media.db2",
		}); err != nil {
			t.Fatalf("open independent backend %d: %v", index, err)
		}
		var firstTrailer, secondTrailer metadata.MD
		_, err = notes.NewNotesServiceClient(connection).GetNote(t.Context(), &notes.NoteId{Nid: -1}, grpc.Trailer(&firstTrailer), grpc.Trailer(&secondTrailer))
		if status.Code(err) != codes.Unknown || len(firstTrailer.Get("anki-error-bin")) != 1 || len(secondTrailer.Get("anki-error-bin")) != 1 {
			t.Fatalf("native completion metadata was not shared with callers: %v %v %v", err, firstTrailer, secondTrailer)
		}
		_, err = notes.NewNotesServiceClient(connection).GetNote(t.Context(), &notes.NoteId{Nid: 1}, grpc.MaxCallSendMsgSize(0))
		if status.Code(err) != codes.ResourceExhausted {
			t.Fatalf("direct request size limit: %v", err)
		}
		_, err = connection.NewStream(t.Context(), &grpc.StreamDesc{ServerStreams: true}, "/unsupported")
		if status.Code(err) != codes.Unimplemented {
			t.Fatalf("direct streaming was not rejected: %v", err)
		}
		if err := connection.Invoke(t.Context(), "/unsupported", &generic.Empty{}, &generic.Empty{}, grpc.UseCompressor("gzip")); status.Code(err) != codes.Unimplemented {
			t.Fatalf("unsupported direct call option: %v", err)
		}
	}
	firstConnection, err := first.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := collection.NewBackendCollectionServiceClient(firstConnection).CloseCollection(t.Context(), &collection.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	if err := first.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	connection, err := second.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := collection.NewBackendCollectionServiceClient(connection).LatestProgress(t.Context(), &generic.Empty{}); err != nil {
		t.Fatalf("stopping another named backend affected this instance: %v", err)
	}
	if _, err := collection.NewBackendCollectionServiceClient(connection).CloseCollection(t.Context(), &collection.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	writeReport(t, "direct-ownership.json", map[string]any{"independent_instances": true, "native_error_trailers": true, "explicit_call_options": true})
}

// TestDirectBackendCancellation retains native state until a cancelled media call finishes.
func TestDirectBackendCancellation(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 15*time.Second)
	defer cancel()
	backend := native.New()
	if err := backend.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := backend.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	})
	connection, err := backend.Connection()
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "cancelled.anki2")
	if _, err := collection.NewBackendCollectionServiceClient(connection).OpenCollection(ctx, &collection.OpenCollectionRequest{
		CollectionPath: path, MediaFolderPath: path + ".media", MediaDbPath: path + ".media.db2",
	}); err != nil {
		t.Fatal(err)
	}
	entered, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		close(entered)
		<-release
		if _, err := response.Write([]byte("direct native media")); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	defer releaseOnce.Do(func() { close(release) })
	download, cancelDownload := context.WithCancel(ctx)
	defer cancelDownload()
	downloaded := make(chan error, 1)
	go func() {
		_, err := media.NewBackendMediaServiceClient(connection).AddMediaFromUrl(download, &media.AddMediaFromUrlRequest{Url: server.URL + "/pending.txt"})
		downloaded <- err
	}()
	select {
	case <-entered:
	case err := <-downloaded:
		t.Fatalf("native call did not reach media fixture: %v", err)
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if _, err := collection.NewBackendCollectionServiceClient(connection).LatestProgress(ctx, &generic.Empty{}); err != nil {
		t.Fatalf("direct native work blocked progress access: %v", err)
	}
	cancelDownload()
	select {
	case err := <-downloaded:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancelled direct call: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("cancelled caller kept waiting for native work")
	}
	shutdown, cancelShutdown := context.WithTimeout(ctx, 25*time.Millisecond)
	defer cancelShutdown()
	if err := backend.Stop(shutdown); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("direct stop did not retain active native work: %v", err)
	}
	if _, err := backend.Connection(); err == nil {
		t.Fatal("stopping direct backend exposed a connection")
	}
	if _, err := collection.NewBackendCollectionServiceClient(connection).LatestProgress(ctx, &generic.Empty{}); err == nil {
		t.Fatal("retained direct client accepted work after Stop")
	}
	releaseOnce.Do(func() { close(release) })
	if err := backend.Stop(ctx); err != nil {
		t.Fatalf("direct cleanup did not finish after native completion: %v", err)
	}
	if err := backend.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	writeReport(t, "direct-cancellation.json", map[string]any{"cancelled_caller": true, "active_native_work_drained": true, "repeatable_stop": true})
}

// TestDirectBackendStartup verifies terminal shutdown and cancelled startup with native ownership.
func TestDirectBackendStartup(t *testing.T) {
	t.Run("stop-during-startup", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()
		backend := native.New()
		done := backend.Done()
		barrier := &startupBarrierContext{Context: ctx, entered: make(chan struct{}), release: make(chan struct{})}
		var releaseOnce sync.Once
		release := func() { releaseOnce.Do(func() { close(barrier.release) }) }
		defer release()
		started := make(chan error, 1)
		go func() { started <- backend.Start(barrier) }()
		select {
		case <-barrier.entered:
		case <-ctx.Done():
			t.Fatal("startup did not reach its controlled scheduling point")
		}
		shutdown, cancelShutdown := context.WithTimeout(ctx, 25*time.Millisecond)
		defer cancelShutdown()
		stopped := make(chan error, 1)
		go func() { stopped <- backend.Stop(shutdown) }()
		select {
		case err := <-stopped:
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("Stop did not preserve its startup-overlap deadline: %v", err)
			}
		case <-time.After(500 * time.Millisecond):
			release()
			<-started
			<-stopped
			t.Fatal("Stop blocked behind synchronous startup instead of observing its deadline")
		}
		select {
		case <-done:
			t.Fatal("Done closed before startup drained")
		default:
		}
		if _, err := backend.Connection(); err == nil {
			t.Fatal("stopped startup exposed a connection")
		}
		release()
		if err := <-started; err == nil {
			t.Fatal("startup succeeded after terminal shutdown")
		}
		if err := backend.Stop(ctx); err != nil {
			t.Fatalf("startup cleanup did not finish: %v", err)
		}
		select {
		case <-done:
		default:
			t.Fatal("Done did not observe completed startup cleanup")
		}
		writeReport(t, "direct-startup-stop.json", map[string]any{"shutdown_deadline_bounded": true, "startup_drained": true, "stopped_connection_rejected": true})
	})
	for _, beforeStart := range []bool{false, true} {
		backend := native.New()
		if beforeStart {
			if err := backend.Stop(t.Context()); err != nil {
				t.Fatal(err)
			}
		}
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		err := backend.Start(ctx)
		if beforeStart && err == nil || !beforeStart && !errors.Is(err, context.Canceled) {
			t.Fatalf("invalid direct startup: stopped=%v error=%v", beforeStart, err)
		}
		if err := backend.Stop(t.Context()); err != nil {
			t.Fatal(err)
		}
		if err := backend.Start(t.Context()); err == nil {
			t.Fatal("terminal direct backend restarted")
		}
	}
}
