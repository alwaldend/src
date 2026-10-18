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
	// release resumes startup after the competing Stop has finished.
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

// TestCollectionBackend verifies native persistence, package ownership and a reusable borrowed ankiclient.
func TestCollectionBackend(t *testing.T) {
	ctx := t.Context()
	worker := ankiclient.New(ankiclient.BackendSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}})
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
	options := ankiclient.BackendCollectionOptions{Backend: worker, Path: database}
	managed := ankiclient.NewCollection(options)
	var instance ankiclient.Backend = managed
	options.Path = filepath.Join(root, "must-not-be-used.anki2")
	if _, err := instance.Connection(); err == nil {
		t.Fatal("unprepared collection exposed a connection")
	}
	if err := instance.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if err := instance.Start(ctx); err == nil {
		t.Fatal("collection accepted another start")
	}
	connection, err := instance.Connection()
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
	if err := instance.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	if err := instance.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := instance.Connection(); err == nil {
		t.Fatal("stopped collection exposed a connection")
	}
	if _, err := items.GetNote(ctx, &notes.NoteId{Nid: added.NoteId}); err == nil {
		t.Fatal("retained collection client survived Stop")
	}
	if _, err := worker.Connection(); err != nil {
		t.Fatalf("collection stopped its borrowed backend: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "collection.media", "example.txt")); err != nil {
		t.Fatalf("standard media layout: %v", err)
	}
	packagePath := filepath.Join(root, "collection.colpkg")
	managed = ankiclient.NewCollection(ankiclient.BackendCollectionOptions{Backend: worker, Path: database})
	if err := managed.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err = managed.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: packagePath, IncludeMedia: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := managed.Connection(); err == nil {
		t.Fatal("consuming export left collection access prepared")
	}
	if err := managed.Stop(ctx); err != nil {
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
		managed = ankiclient.NewCollection(ankiclient.BackendCollectionOptions{Backend: worker, Path: packagePath, RestoreDirectory: directory})
		if err := managed.Start(ctx); err != nil {
			t.Fatal(err)
		}
		connection, err := managed.Connection()
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
		if err := managed.Stop(ctx); err != nil {
			t.Fatal(err)
		}
		_, err = os.Stat(restored)
		if persistent && err != nil || !persistent && !os.IsNotExist(err) {
			t.Fatalf("restore directory ownership: %v", err)
		}
		if _, err := worker.Connection(); err != nil {
			t.Fatal("restore cleanup stopped the borrowed backend")
		}
	}
	managed = ankiclient.NewCollection(ankiclient.BackendCollectionOptions{Backend: worker, Path: packagePath})
	if err := managed.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err = managed.Connection()
	if err != nil {
		t.Fatal(err)
	}
	restored := collectionRestoreDirectory(t, connection)
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: filepath.Join(root, "missing", "failed.colpkg")}); err == nil {
		t.Fatal("export to a missing directory succeeded")
	}
	if err := managed.Stop(ctx); err != nil {
		t.Fatalf("completed failed export prevented collection cleanup: %v", err)
	}
	if _, err := os.Stat(restored); !os.IsNotExist(err) {
		t.Fatalf("completed failed export retained temporary restore: %v", err)
	}
	startup := &startupBarrierContext{Context: ctx, entered: make(chan struct{}), release: make(chan struct{})}
	raceDirectory := filepath.Join(root, "stopped-before-restore")
	managed = ankiclient.NewCollection(ankiclient.BackendCollectionOptions{Backend: worker, Path: packagePath, RestoreDirectory: raceDirectory})
	started := make(chan error, 1)
	go func() { started <- managed.Start(startup) }()
	<-startup.entered
	stopError := managed.Stop(ctx)
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
		t.Fatal("startup prepared restore files after Stop had completed")
	}
	managed = ankiclient.NewCollection(ankiclient.BackendCollectionOptions{Backend: worker, Path: database})
	if err := managed.Start(ctx); err != nil {
		t.Fatalf("overlapping lifecycle leaked an open collection: %v", err)
	}
	if err := managed.Stop(ctx); err != nil {
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
	writeReport(t, "collection-ankiclient.json", map[string]any{"note_id": added.NoteId, "temporary_restore_removed": true, "persistent_restore_retained": true, "borrowed_backend_retained": true})

	for _, cancelCall := range []bool{false, true} {
		scenario := "active-call"
		if cancelCall {
			scenario = "cancelled-call"
		}
		t.Run(scenario, func(t *testing.T) { testCollectionCallCleanup(t, packagePath, cancelCall) })
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
func testCollectionCallCleanup(t *testing.T, packagePath string, cancelCall bool) {
	t.Helper()
	ctx := t.Context()
	worker := ankiclient.New(ankiclient.BackendSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}})
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	managed := ankiclient.NewCollection(ankiclient.BackendCollectionOptions{Backend: worker, Path: packagePath})
	var restored string
	t.Cleanup(func() {
		if err := worker.Stop(context.Background()); err != nil && !errors.Is(err, context.Canceled) {
			t.Error(err)
		}
		if err := os.RemoveAll(restored); err != nil {
			t.Error(err)
		}
	})
	if err := managed.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err := managed.Connection()
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
	if cancelCall {
		cancelDownload()
		if err := <-downloaded; !errors.Is(err, context.Canceled) {
			t.Fatalf("cancelled native download: %v", err)
		}
		if _, err := managed.Connection(); err == nil {
			t.Fatal("unconfirmed native call left collection access prepared")
		}
		releaseOnce.Do(func() { close(release) })
	}
	shutdown, cancelShutdown := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancelShutdown()
	err = managed.Stop(shutdown)
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
	if err := worker.Stop(termination); !errors.Is(err, context.Canceled) {
		t.Fatalf("reap backend after unconfirmed work: %v", err)
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

// TestCollectionBackendFailures verifies preparation and shutdown errors never take over a borrowed process.
func TestCollectionBackendFailures(t *testing.T) {
	ctx := t.Context()
	worker := ankiclient.New(ankiclient.BackendSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}})
	root := t.TempDir()
	path := filepath.Join(root, "collection.anki2")
	unprepared := ankiclient.NewCollection(ankiclient.BackendCollectionOptions{Backend: worker, Path: path})
	if err := unprepared.Start(ctx); err == nil {
		t.Fatal("collection started an unprepared backend")
	}
	if _, err := worker.Connection(); err == nil {
		t.Fatal("failed collection preparation launched its backend")
	}
	if err := unprepared.Stop(ctx); err != nil {
		t.Fatal(err)
	}
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := worker.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	})
	for _, scenario := range []string{"nil backend", "cancelled start", "stop before start", "nonempty restore", "invalid package"} {
		t.Run(scenario, func(t *testing.T) {
			options := ankiclient.BackendCollectionOptions{Backend: worker, Path: path}
			startup, cancel := context.WithCancel(ctx)
			defer cancel()
			if scenario == "nil backend" {
				options.Backend = nil
			}
			if scenario == "cancelled start" {
				cancel()
			}
			if scenario == "nonempty restore" || scenario == "invalid package" {
				options.Path = filepath.Join(root, "invalid.colpkg")
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
			managed := ankiclient.NewCollection(options)
			if scenario == "stop before start" {
				if err := managed.Stop(ctx); err != nil {
					t.Fatal(err)
				}
			}
			if err := managed.Start(startup); err == nil {
				t.Fatal("invalid collection start succeeded")
			}
			if _, err := managed.Connection(); err == nil {
				t.Fatal("failed collection exposed a connection")
			}
			if err := managed.Stop(ctx); err != nil {
				t.Fatal(err)
			}
			if _, err := worker.Connection(); err != nil {
				t.Fatal("failed collection stopped its borrowed backend")
			}
			if options.RestoreDirectory != "" {
				if _, err := os.Stat(filepath.Join(options.RestoreDirectory, "keep.txt")); err != nil {
					t.Fatal("failed restore removed caller files")
				}
			}
		})
	}
	managed := ankiclient.NewCollection(ankiclient.BackendCollectionOptions{Backend: worker, Path: path})
	if err := managed.Start(ctx); err != nil {
		t.Fatal(err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if err := managed.Stop(cancelled); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled collection stop: %v", err)
	}
	if _, err := managed.Connection(); err == nil {
		t.Fatal("cancelled stop exposed the collection")
	}
	shutdown, finish := context.WithTimeout(context.Background(), 5*time.Second)
	defer finish()
	if err := worker.Stop(shutdown); err != nil {
		t.Fatal(err)
	}
}
