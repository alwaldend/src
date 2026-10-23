//go:build cgo && linux && amd64

package ankiclient_test

import (
	"context"
	"errors"
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
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/grpc"
)

// startupBarrierContext pauses a context check to exercise overlapping native lifecycle calls.
type startupBarrierContext struct {
	// Context provides the real test cancellation and deadline.
	context.Context
	// entered reports that startup reached the controlled scheduling point.
	entered chan struct{}
	// release resumes startup after the competing CloseCollection has finished.
	release chan struct{}
	// once pauses only the first check, allowing later gRPC context use.
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

// TestUnifiedClientCollections verifies collection changes retain one transport and invalidate old connections.
func TestUnifiedClientCollections(t *testing.T) {
	registerWorkerFactories(t)
	client, err := ankigrpc.NewClient(ankigrpc.SubprocessAddress(workerPath(t), "serve"))
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Start(t.Context()); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	}()
	for _, name := range []string{"first", "second"} {
		if err := client.OpenCollection(t.Context(), ankiclient.CollectionOptions{Path: filepath.Join(t.TempDir(), name+".anki2")}); err != nil {
			t.Fatal(err)
		}
		connection, err := client.Connection()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := notetypes.NewNotetypesServiceClient(connection).GetNotetypeIdByName(t.Context(), &generic.String{Val: "Basic"}); err != nil {
			t.Fatal(err)
		}
		if err := client.CloseCollection(t.Context()); err != nil {
			t.Fatal(err)
		}
		if _, err := notetypes.NewNotetypesServiceClient(connection).GetNotetypeIdByName(t.Context(), &generic.String{Val: "Basic"}); err == nil {
			t.Fatal("closed collection retained a usable connection")
		}
		connection, err = client.Connection()
		if err != nil {
			t.Fatal(err)
		}
		if _, err := collection.NewBackendCollectionServiceClient(connection).LatestProgress(t.Context(), &generic.Empty{}); err != nil {
			t.Fatalf("closing collection stopped its transport: %v", err)
		}
	}
	if err := client.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	<-client.Done()
	writeReport(t, "unified-client.json", map[string]any{"collections": 2, "transport_reused": true, "native_completion": true})
}

// TestOwnedCollection verifies worker ownership, native access, failed preparation and terminal cleanup.
func TestOwnedCollection(t *testing.T) {
	for _, scenario := range []string{"success", "invalid-source", "stop-before-start"} {
		t.Run(scenario, func(t *testing.T) {
			worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}})
			database := filepath.Join(t.TempDir(), "owned.anki2")
			if scenario == "invalid-source" {
				database += ".invalid"
			}
			client := ankiclient.NewClient(worker)
			defer func() {
				if err := client.Stop(context.Background()); err != nil {
					t.Error(err)
				}
			}()
			if _, err := client.Connection(); err == nil {
				t.Fatal("unstarted owner exposed a connection")
			}
			if scenario == "stop-before-start" {
				if err := client.Stop(t.Context()); err != nil {
					t.Fatal(err)
				}
			}
			startup, cancelStartup := context.WithCancel(t.Context())
			err := client.Start(startup)
			cancelStartup()
			if scenario != "stop-before-start" {
				if err != nil {
					t.Fatal(err)
				}
				err = client.OpenCollection(t.Context(), ankiclient.CollectionOptions{Path: database, CloseTimeout: 10 * time.Second})
			}
			if scenario == "success" {
				if err != nil {
					t.Fatal(err)
				}
				connection, err := client.Connection()
				if err != nil {
					t.Fatal(err)
				}
				if _, err := collection.NewBackendCollectionServiceClient(connection).CloseCollection(t.Context(), &collection.CloseCollectionRequest{}); err != nil {
					t.Fatal(err)
				}
			} else if err == nil {
				t.Fatal("invalid lifecycle operation succeeded")
			}
			if err := client.Stop(t.Context()); err != nil {
				t.Fatal(err)
			}
			select {
			case <-client.Done():
			default:
				t.Fatal("successful Stop returned before native completion")
			}
			if _, err := worker.Connection(); err == nil {
				t.Fatal("owning client left its worker available")
			}
			if _, err := client.Connection(); err == nil {
				t.Fatal("stopped owner exposed a connection")
			}
			if scenario == "success" {
				data, err := os.ReadFile(database)
				if err != nil {
					t.Fatal(err)
				}
				if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
					if err := os.WriteFile(filepath.Join(output, "owned-client.anki2"), data, 0o600); err != nil {
						t.Fatal(err)
					}
				}
			}
			writeReport(t, "owned-client-"+scenario+".json", map[string]any{"scenario": scenario, "worker_stopped": true, "native_completion": true})
		})
	}
}

// TestCollectionBackend verifies native persistence and package ownership while reusing one client.
func TestCollectionBackend(t *testing.T) {
	ctx := t.Context()
	worker := ankiclient.NewClient(ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}}))
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := worker.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	})
	root := t.TempDir()
	database := filepath.Join(root, "collection.anki2")
	options := ankiclient.CollectionOptions{Path: database}
	if err := worker.OpenCollection(ctx, options); err != nil {
		t.Fatal(err)
	}
	options.Path = filepath.Join(root, "must-not-be-used.anki2")
	if err := worker.OpenCollection(ctx, options); err == nil {
		t.Fatal("client accepted another collection before closing the current one")
	}
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	types := notetypes.NewNotetypesServiceClient(connection)
	definition, err := types.GetNotetypeIdByName(ctx, &generic.String{Val: "Basic"})
	if err != nil {
		t.Fatal(err)
	}
	items := notes.NewNotesServiceClient(connection)
	added, err := items.AddNote(ctx, &notes.AddNoteRequest{
		Note: &notes.Note{NotetypeId: definition.Ntid, Fields: []string{"front", "back"}}, DeckId: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := media.NewMediaServiceClient(connection).AddMediaFile(ctx, &media.AddMediaFileRequest{DesiredName: "example.txt", Data: []byte("collection media")}); err != nil {
		t.Fatal(err)
	}
	if err := worker.CloseCollection(ctx); err != nil {
		t.Fatal(err)
	}
	if err := worker.CloseCollection(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := items.GetNote(ctx, &notes.NoteId{Nid: added.NoteId}); err == nil {
		t.Fatal("retained collection connection survived CloseCollection")
	}
	if _, err := worker.Connection(); err != nil {
		t.Fatalf("collection close stopped its transport: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "collection.media", "example.txt")); err != nil {
		t.Fatalf("standard media layout: %v", err)
	}
	packagePath := filepath.Join(root, "collection.colpkg")
	if err := worker.OpenCollection(ctx, ankiclient.CollectionOptions{Path: database}); err != nil {
		t.Fatal(err)
	}
	connection, err = worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: packagePath, IncludeMedia: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := worker.Connection(); err == nil {
		t.Fatal("consuming export left collection access prepared")
	}
	if err := worker.CloseCollection(ctx); err != nil {
		t.Fatalf("consumed collection was closed twice: %v", err)
	}
	original, err := os.ReadFile(packagePath)
	if err != nil {
		t.Fatal(err)
	}
	for _, persistent := range []bool{false, true} {
		directory := ""
		if persistent {
			directory = filepath.Join(t.TempDir(), "restore")
		}
		if err := worker.OpenCollection(ctx, ankiclient.CollectionOptions{Path: packagePath, RestoreDirectory: directory}); err != nil {
			t.Fatal(err)
		}
		connection, err := worker.Connection()
		if err != nil {
			t.Fatal(err)
		}
		actual, err := notes.NewNotesServiceClient(connection).GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
		if err != nil || actual.GetFields()[0] != "front" {
			t.Fatalf("restored note: %v", err)
		}
		mediaPath, err := media.NewMediaServiceClient(connection).GetAbsoluteMediaPath(ctx, &generic.String{Val: "example.txt"})
		if err != nil {
			t.Fatal(err)
		}
		data, err := os.ReadFile(mediaPath.Val)
		if err != nil || string(data) != "collection media" {
			t.Fatalf("restored media: %v", err)
		}
		restored := filepath.Dir(filepath.Dir(mediaPath.Val))
		if restored == "" || persistent && restored != directory {
			t.Fatal("native media path did not use the selected restore destination")
		}
		if err := worker.CloseCollection(ctx); err != nil {
			t.Fatal(err)
		}
		_, err = os.Stat(restored)
		if persistent && err != nil || !persistent && !os.IsNotExist(err) {
			t.Fatalf("restore directory ownership: %v", err)
		}
		if _, err := worker.Connection(); err != nil {
			t.Fatal("restore cleanup stopped the transport")
		}
	}
	if err := worker.OpenCollection(ctx, ankiclient.CollectionOptions{Path: packagePath}); err != nil {
		t.Fatal(err)
	}
	connection, err = worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	restored := collectionRestoreDirectory(t, connection)
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: filepath.Join(root, "missing", "failed.colpkg")}); err == nil {
		t.Fatal("export to a missing directory succeeded")
	}
	if err := worker.CloseCollection(ctx); err != nil {
		t.Fatalf("completed failed export prevented collection cleanup: %v", err)
	}
	if _, err := os.Stat(restored); !os.IsNotExist(err) {
		t.Fatalf("completed failed export retained temporary restore: %v", err)
	}
	startup := &startupBarrierContext{Context: ctx, entered: make(chan struct{}), release: make(chan struct{})}
	raceDirectory := filepath.Join(root, "stopped-before-restore")
	started := make(chan error, 1)
	go func() {
		started <- worker.OpenCollection(startup, ankiclient.CollectionOptions{Path: packagePath, RestoreDirectory: raceDirectory})
	}()
	<-startup.entered
	stopError := worker.CloseCollection(ctx)
	close(startup.release)
	startError := <-started
	if stopError != nil || startError == nil {
		t.Fatalf("overlapping lifecycle result: start=%v stop=%v", startError, stopError)
	}
	if _, err := os.Stat(raceDirectory); !os.IsNotExist(err) {
		// Reap the real server before reporting any leaked native collection.
		if err := worker.Stop(ctx); err != nil {
			t.Fatal(err)
		}
		if err := os.RemoveAll(raceDirectory); err != nil {
			t.Fatal(err)
		}
		t.Fatal("startup prepared restore files after CloseCollection had completed")
	}
	if err := worker.OpenCollection(ctx, ankiclient.CollectionOptions{Path: database}); err != nil {
		t.Fatalf("overlapping lifecycle leaked an open collection: %v", err)
	}
	if err := worker.CloseCollection(ctx); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(packagePath)
	if err != nil || string(original) != string(after) {
		t.Fatal("package source was modified")
	}
	for _, name := range []string{"collection.anki2", "collection.colpkg"} {
		if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
			data, err := os.ReadFile(filepath.Join(root, name))
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(output, "backend-"+name), data, 0o600); err != nil {
				t.Fatal(err)
			}
		}
	}
	writeReport(t, "collection-ankiclient.json", map[string]any{"note_id": added.NoteId, "temporary_restore_removed": true, "persistent_restore_retained": true, "transport_retained": true})

	for _, scenario := range []string{"active-call", "cancelled-call", "concurrent-close"} {
		t.Run(scenario, func(t *testing.T) { testCollectionCallCleanup(t, packagePath, scenario) })
	}
}

// collectionRestoreDirectory locates a package restore independently through Anki's media API.
func collectionRestoreDirectory(t *testing.T, connection grpc.ClientConnInterface) string {
	t.Helper()
	path, err := media.NewMediaServiceClient(connection).GetAbsoluteMediaPath(t.Context(), &generic.String{Val: "example.txt"})
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Dir(filepath.Dir(path.Val))
}

// testCollectionCallCleanup verifies active or unconfirmed native calls retain restore files.
func testCollectionCallCleanup(t *testing.T, packagePath, scenario string) {
	t.Helper()
	ctx := t.Context()
	cancelCall := scenario == "cancelled-call"
	transport := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}})
	worker := ankiclient.NewClient(transport)
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	var restored string
	var closeCause error
	t.Cleanup(func() {
		if err := worker.Stop(context.Background()); closeCause == nil && err != nil || closeCause != nil && !errors.Is(err, closeCause) {
			t.Error(err)
		}
		<-worker.Done()
		if err := transport.Stop(context.Background()); err != nil {
			t.Errorf("transport cleanup: %v", err)
		}
		if err := os.RemoveAll(restored); err != nil {
			t.Error(err)
		}
	})
	if err := worker.OpenCollection(ctx, ankiclient.CollectionOptions{Path: packagePath, CloseTimeout: 100 * time.Millisecond}); err != nil {
		t.Fatal(err)
	}
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	restored = collectionRestoreDirectory(t, connection)
	entered, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		close(entered)
		<-release
		response.Header().Set("Content-Type", "text/plain")
		if _, err := response.Write([]byte("pending collection media")); err != nil {
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
		t.Fatalf("native download returned before reaching the local fixture: %v", err)
	case <-ctx.Done():
		t.Fatal("native media request did not reach the local fixture")
	}
	progress, cancelProgress := context.WithTimeout(ctx, time.Second)
	defer cancelProgress()
	if _, err := collection.NewBackendCollectionServiceClient(connection).LatestProgress(progress, &generic.Empty{}); err != nil {
		t.Fatalf("ordinary calls prevented concurrent progress access: %v", err)
	}
	if scenario == "concurrent-close" {
		closing := make(chan error, 1)
		go func() { closing <- worker.CloseCollection(context.Background()) }()
		waitForRejection := func(connection func() (grpc.ClientConnInterface, error), failure string) {
			t.Helper()
			deadline := time.NewTimer(2 * time.Second)
			defer deadline.Stop()
			poll := time.NewTicker(time.Millisecond)
			defer poll.Stop()
			for {
				if _, err := connection(); err != nil {
					return
				}
				select {
				case <-poll.C:
				case <-deadline.C:
					t.Fatal(failure)
				}
			}
		}
		waitForRejection(worker.Connection, "first collection close did not invalidate access")
		stopped := make(chan error, 1)
		go func() { stopped <- worker.Stop(context.Background()) }()
		waitForRejection(transport.Connection, "terminal cleanup waited beyond its collection budget before stopping the transport")
		releaseOnce.Do(func() { close(release) })
		completion, cancelCompletion := context.WithTimeout(ctx, 5*time.Second)
		defer cancelCompletion()
		select {
		case closeCause = <-stopped:
			for errors.Unwrap(closeCause) != nil {
				closeCause = errors.Unwrap(closeCause)
			}
			if !errors.Is(closeCause, context.DeadlineExceeded) {
				t.Fatalf("terminal cleanup lost its collection deadline: %v", closeCause)
			}
		case <-completion.Done():
			t.Fatal("terminal cleanup did not finish after native work was released")
		}
		select {
		case <-worker.Done():
		case <-completion.Done():
			t.Fatal("terminal cleanup did not observe final native completion")
		}
		select {
		case <-closing:
		case <-completion.Done():
			t.Fatal("first collection close did not finish after transport shutdown")
		}
		select {
		case <-downloaded:
		case <-completion.Done():
			t.Fatal("native media caller did not finish after transport shutdown")
		}
		if _, err := os.Stat(restored); err != nil {
			t.Fatalf("concurrent cleanup discarded an unconfirmed restore: %v", err)
		}
		writeReport(t, "concurrent-close.json", map[string]any{"collection_deadline_preserved": true, "transport_shutdown_before_native_release": true, "native_completion": true, "restore_retained": true})
		return
	}
	if cancelCall {
		cancelDownload()
		if err := <-downloaded; !errors.Is(err, context.Canceled) {
			t.Fatalf("cancelled native download: %v", err)
		}
		if _, err := worker.Connection(); err == nil {
			t.Fatal("unconfirmed native call left collection access prepared")
		}
		releaseOnce.Do(func() { close(release) })
	}
	shutdown, cancelShutdown := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancelShutdown()
	err = worker.CloseCollection(shutdown)
	closeCause = err
	for errors.Unwrap(closeCause) != nil {
		closeCause = errors.Unwrap(closeCause)
	}
	if cancelCall {
		if err == nil || errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("unconfirmed native call was treated as completed: %v", err)
		}
	} else {
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("shutdown did not wait for the active call: %v", err)
		}
		cancelDownload()
		if err := <-downloaded; !errors.Is(err, context.Canceled) {
			t.Fatalf("cancelled active native download: %v", err)
		}
		releaseOnce.Do(func() { close(release) })
	}
	if _, err := os.Stat(restored); err != nil {
		t.Fatalf("shutdown removed an unconfirmed native restore: %v", err)
	}
	termination, cancelTermination := context.WithCancel(ctx)
	cancelTermination()
	if err := worker.Stop(termination); !errors.Is(err, context.Canceled) && !errors.Is(err, closeCause) {
		t.Fatalf("stop client after unconfirmed work: %v", err)
	}
	<-worker.Done()
	if _, err := transport.Connection(); err == nil {
		t.Fatal("native completion left the transport available")
	}
	if err := os.RemoveAll(restored); err != nil {
		t.Fatal(err)
	}
	report := "active-call.json"
	if cancelCall {
		report = "cancelled-call.json"
	}
	writeReport(t, report, map[string]any{"active_restore_retained": true, "cleanup_after_reaping": true})
}

// TestCollectionBackendFailures verifies preparation failures retain a reusable transport and terminal cleanup.
func TestCollectionBackendFailures(t *testing.T) {
	ctx := t.Context()
	root := t.TempDir()
	path := filepath.Join(root, "collection.anki2")
	transport := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}})
	worker := ankiclient.NewClient(transport)
	if err := worker.OpenCollection(ctx, ankiclient.CollectionOptions{Path: path}); err == nil {
		t.Fatal("unstarted client prepared a collection")
	}
	if _, err := transport.Connection(); err == nil {
		t.Fatal("failed collection preparation launched its transport")
	}
	if err := worker.CloseCollection(ctx); err != nil {
		t.Fatal(err)
	}
	if err := worker.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	if err := worker.Start(ctx); err == nil {
		t.Fatal("stopped client accepted startup")
	}
	withoutTransport := ankiclient.NewClient(nil)
	if err := withoutTransport.Start(ctx); err == nil {
		t.Fatal("client without a transport started")
	}
	if err := withoutTransport.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	for _, scenario := range []string{"cancelled open", "nonempty restore", "invalid package"} {
		t.Run(scenario, func(t *testing.T) {
			worker := ankiclient.NewClient(ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}}))
			if err := worker.Start(t.Context()); err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := worker.Stop(context.Background()); err != nil {
					t.Error(err)
				}
				<-worker.Done()
			})
			options := ankiclient.CollectionOptions{Path: filepath.Join(t.TempDir(), "collection.anki2")}
			startup, cancel := context.WithCancel(t.Context())
			defer cancel()
			if scenario == "cancelled open" {
				cancel()
			}
			if scenario == "nonempty restore" || scenario == "invalid package" {
				options.Path = filepath.Join(t.TempDir(), "invalid.colpkg")
				if err := os.WriteFile(options.Path, []byte("not an Anki package"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if scenario == "nonempty restore" {
				options.RestoreDirectory = t.TempDir()
				if err := os.WriteFile(filepath.Join(options.RestoreDirectory, "keep.txt"), []byte("keep"), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			if err := worker.OpenCollection(startup, options); err == nil {
				t.Fatal("invalid collection preparation succeeded")
			}
			if _, err := worker.Connection(); err == nil {
				t.Fatal("failed collection exposed a connection")
			}
			if err := worker.CloseCollection(t.Context()); err != nil {
				t.Fatal(err)
			}
			connection, err := worker.Connection()
			if err != nil {
				t.Fatal(err)
			}
			if _, err := collection.NewBackendCollectionServiceClient(connection).LatestProgress(t.Context(), &generic.Empty{}); err != nil {
				t.Fatalf("failed collection preparation stopped its transport: %v", err)
			}
			if options.RestoreDirectory != "" {
				if _, err := os.Stat(filepath.Join(options.RestoreDirectory, "keep.txt")); err != nil {
					t.Fatal("failed restore removed caller files")
				}
			}
		})
	}
	transport = ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}})
	worker = ankiclient.NewClient(transport)
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	var closeCause error
	t.Cleanup(func() {
		if err := worker.Stop(context.Background()); closeCause == nil && err != nil || closeCause != nil && !errors.Is(err, closeCause) {
			t.Error(err)
		}
		<-worker.Done()
		if err := transport.Stop(context.Background()); err != nil {
			t.Errorf("transport cleanup: %v", err)
		}
	})
	if err := worker.OpenCollection(ctx, ankiclient.CollectionOptions{Path: path}); err != nil {
		t.Fatal(err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	closeCause = worker.CloseCollection(cancelled)
	for errors.Unwrap(closeCause) != nil {
		closeCause = errors.Unwrap(closeCause)
	}
	if !errors.Is(closeCause, context.Canceled) {
		t.Fatalf("cancelled collection close: %v", closeCause)
	}
	if _, err := worker.Connection(); err == nil {
		t.Fatal("cancelled close exposed the collection")
	}
	shutdown, finish := context.WithTimeout(context.Background(), 5*time.Second)
	defer finish()
	if err := worker.Stop(shutdown); !errors.Is(err, context.Canceled) {
		t.Fatalf("terminal cleanup lost the collection close error: %v", err)
	}
	select {
	case <-worker.Done():
	case <-shutdown.Done():
		t.Fatal("terminal cleanup did not finish its native transport")
	}
	if _, err := transport.Connection(); err == nil {
		t.Fatal("terminal cleanup retained the native transport")
	}
	writeReport(t, "collection-failures.json", map[string]any{"failed_preparation_retains_transport": true, "caller_files_retained": true, "cancelled_close_preserved": true, "native_completion": true})
}
