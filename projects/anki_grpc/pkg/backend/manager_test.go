package backend_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
)

// TestManagerSTD verifies real process scopes, copied options and cleanup on every exit.
func TestManagerSTD(t *testing.T) {
	t.Setenv("ANKI_E2E_RECORD_PID", "1")
	arguments := []string{"-test.run=^TestWorkerFailure$", "-worker=" + *workerBinary}
	var manager backend.Manager = backend.NewManagerSTD(backend.BackendSTDOptions{Binary: os.Args[0], Arguments: arguments})
	arguments[0] = "-invalid-option"
	for _, scenario := range []string{"success", "failure", "cancel", "panic", "crash"} {
		t.Run(scenario, func(t *testing.T) {
			pidPath := filepath.Join(t.TempDir(), "worker.pid")
			t.Setenv("ANKI_E2E_PID_FILE", pidPath)
			ctx, cancel := context.WithCancelCause(t.Context())
			defer cancel(nil)
			failure := errors.New("callback failed")
			var retained backend.Backend
			var pid int
			var recovered any
			var scopeError error
			func() {
				defer func() { recovered = recover() }()
				scopeError = manager.WithBackend(ctx, func(prepared backend.Backend) error {
					retained = prepared
					connection, err := prepared.Connection()
					if err != nil {
						return fmt.Errorf("get managed connection: %w", err)
					}
					if _, err := collection.NewBackendCollectionServiceClient(connection).LatestProgress(ctx, &generic.Empty{}); err != nil {
						return fmt.Errorf("call managed native backend: %w", err)
					}
					data, err := os.ReadFile(pidPath)
					if err != nil {
						return fmt.Errorf("read managed worker PID: %w", err)
					}
					pid, err = strconv.Atoi(string(data))
					if err != nil {
						return fmt.Errorf("parse managed worker PID: %w", err)
					}
					switch scenario {
					case "failure":
						return failure
					case "cancel":
						cancel(failure)
					case "panic":
						panic(failure)
					case "crash":
						if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
							return fmt.Errorf("terminate managed worker: %w", err)
						}
						return failure
					}
					return nil
				})
			}()
			if retained == nil || pid == 0 {
				t.Fatalf("manager never reached its callback: %v", scopeError)
			}
			if _, err := retained.Connection(); err == nil {
				t.Fatal("scope retained a ready backend")
			}
			if err := syscall.Kill(pid, 0); !errors.Is(err, syscall.ESRCH) {
				t.Fatalf("scope retained worker PID %d: %v", pid, err)
			}
			switch scenario {
			case "success":
				if scopeError != nil || recovered != nil {
					t.Fatalf("successful scope failed: %v %v", scopeError, recovered)
				}
			case "cancel":
				if !errors.Is(scopeError, context.Canceled) || !errors.Is(scopeError, failure) {
					t.Fatalf("scope lost cancellation: %v", scopeError)
				}
			case "panic":
				if recovered != failure {
					t.Fatalf("scope replaced callback panic: %v", recovered)
				}
			default:
				if !errors.Is(scopeError, failure) {
					t.Fatalf("scope lost callback error: %v", scopeError)
				}
				if scenario == "crash" {
					var exit *exec.ExitError
					if !errors.As(scopeError, &exit) {
						t.Fatalf("scope lost worker crash: %v", scopeError)
					}
					assertReaped(t, exit)
				}
			}
		})
	}
	writeReport(t, "manager-stdio.json", map[string]any{"checks": []string{"copied-options", "fresh-scopes", "native-readiness", "callback-error", "cancelled-callback", "panic-unwinding", "joined-crash-error", "child-reaped"}})
}

// collectionProbeManager delegates real workers and checks collection cleanup before transport release.
type collectionProbeManager struct {
	// delegate creates and owns the actual native transport scopes.
	delegate backend.Manager
	// directory holds native probe databases opened after each collection callback.
	directory string
	// workers records transport instances for post-scope readiness checks.
	workers []backend.Backend
}

// WithBackend verifies that the collection is closed while its worker remains available.
func (manager *collectionProbeManager) WithBackend(ctx context.Context, operation func(backend.Backend) error) error {
	if err := manager.delegate.WithBackend(ctx, func(worker backend.Backend) error {
		manager.workers = append(manager.workers, worker)
		operationError := operation(worker)
		connection, err := worker.Connection()
		if err != nil {
			return errors.Join(operationError, fmt.Errorf("collection scope stopped borrowed transport: %w", err))
		}
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), 5*time.Second)
		defer cancel()
		path := filepath.Join(manager.directory, fmt.Sprintf("probe-%d.anki2", len(manager.workers)))
		client := collection.NewBackendCollectionServiceClient(connection)
		if _, err := client.OpenCollection(cleanup, &collection.OpenCollectionRequest{CollectionPath: path, MediaFolderPath: path + ".media", MediaDbPath: path + ".media.db2"}); err != nil {
			return errors.Join(operationError, fmt.Errorf("collection was not closed before transport release: %w", err))
		}
		if _, err := client.CloseCollection(cleanup, &collection.CloseCollectionRequest{}); err != nil {
			return errors.Join(operationError, fmt.Errorf("close native probe: %w", err))
		}
		if operationError != nil {
			return fmt.Errorf("run probed collection scope: %w", operationError)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("delegate native probe scope: %w", err)
	}
	return nil
}

// TestManagerCollection verifies native persistence, package cleanup and injected-manager ownership.
func TestManagerCollection(t *testing.T) {
	root, restores := t.TempDir(), t.TempDir()
	t.Setenv("TMPDIR", restores)
	transport := &collectionProbeManager{delegate: backend.NewManagerSTD(backend.BackendSTDOptions{Binary: workerPath(t)}), directory: root}
	manager := backend.NewManagerCollection(transport)
	database, archive := filepath.Join(root, "collection.anki2"), filepath.Join(root, "collection.colpkg")
	var noteID int64
	var retained backend.Backend
	if err := manager.WithCollection(t.Context(), database, func(prepared backend.Backend) error {
		retained = prepared
		connection, err := prepared.Connection()
		if err != nil {
			return fmt.Errorf("get collection connection: %w", err)
		}
		definition, err := notetypes.NewNotetypesServiceClient(connection).GetNotetypeIdByName(t.Context(), &generic.String{Val: "Basic"})
		if err != nil {
			return fmt.Errorf("get native note type: %w", err)
		}
		added, err := notes.NewNotesServiceClient(connection).AddNote(t.Context(), &notes.AddNoteRequest{DeckId: 1, Note: &notes.Note{NotetypeId: definition.Ntid, Fields: []string{"managed front", "managed back"}}})
		if err != nil {
			return fmt.Errorf("create managed note: %w", err)
		}
		noteID = added.NoteId
		if _, err := media.NewMediaServiceClient(connection).AddMediaFile(t.Context(), &media.AddMediaFileRequest{DesiredName: "managed.txt", Data: []byte("managed media")}); err != nil {
			return fmt.Errorf("create managed media: %w", err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := retained.Connection(); err == nil {
		t.Fatal("collection connection survived scope cleanup")
	}
	if err := manager.WithCollection(t.Context(), database, func(prepared backend.Backend) error {
		connection, err := prepared.Connection()
		if err != nil {
			return fmt.Errorf("get reopened collection: %w", err)
		}
		if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(t.Context(), &import_export.ExportCollectionPackageRequest{OutPath: archive, IncludeMedia: true}); err != nil {
			return fmt.Errorf("export managed archive: %w", err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(archive)
	if err != nil {
		t.Fatal(err)
	}
	for _, scenario := range []string{"success", "failure", "cancel"} {
		ctx, cancel := context.WithCancelCause(t.Context())
		failure := errors.New("collection callback failed")
		err := manager.WithCollection(ctx, archive, func(prepared backend.Backend) error {
			connection, err := prepared.Connection()
			if err != nil {
				return fmt.Errorf("get restored collection: %w", err)
			}
			actual, err := notes.NewNotesServiceClient(connection).GetNote(ctx, &notes.NoteId{Nid: noteID})
			if err != nil {
				return fmt.Errorf("read restored note: %w", err)
			}
			if actual.GetFields()[0] != "managed front" {
				return fmt.Errorf("native note content changed")
			}
			if scenario == "failure" {
				return failure
			}
			if scenario == "cancel" {
				cancel(failure)
			}
			return nil
		})
		cancel(nil)
		if scenario == "success" && err != nil || scenario == "failure" && !errors.Is(err, failure) || scenario == "cancel" && (!errors.Is(err, context.Canceled) || !errors.Is(err, failure)) {
			t.Fatalf("collection scope %s lost its result: %v", scenario, err)
		}
		files, err := os.ReadDir(restores)
		if err != nil || len(files) != 0 {
			t.Fatalf("scope retained completed native restores: %v %v", files, err)
		}
	}
	after, err := os.ReadFile(archive)
	if err != nil || string(original) != string(after) {
		t.Fatalf("manager changed source archive: %v", err)
	}
	for _, worker := range transport.workers {
		if _, err := worker.Connection(); err == nil {
			t.Fatal("collection scope retained its underlying transport")
		}
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		for _, name := range []string{"collection.anki2", "collection.colpkg"} {
			data, err := os.ReadFile(filepath.Join(root, name))
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(output, "managed-"+name), data, 0o600); err != nil {
				t.Fatal(err)
			}
		}
	}
	writeReport(t, "manager-collection.json", map[string]any{"note_id": noteID, "checks": []string{"injected-manager", "collection-before-transport-cleanup", "database-persistence", "package-restoration", "temporary-restore-cleanup", "source-archive-preserved", "callback-error", "callback-cancellation"}})
}

// TestManagerFailures rejects invalid scopes and preserves native startup failures.
func TestManagerFailures(t *testing.T) {
	cause := errors.New("scope cancelled before startup")
	ctx, cancel := context.WithCancelCause(t.Context())
	cancel(cause)
	manager := backend.NewManagerSTD(backend.BackendSTDOptions{Binary: workerPath(t)})
	called := false
	operation := func(backend.Backend) error { called = true; return nil }
	if err := manager.WithBackend(ctx, operation); (!errors.Is(err, context.Canceled) || !errors.Is(err, cause)) || called {
		t.Fatalf("cancelled scope acquired a backend: %v", err)
	}
	if err := manager.WithBackend(t.Context(), nil); err == nil {
		t.Fatal("transport scope accepted a nil callback")
	}
	missing := backend.NewManagerSTD(backend.BackendSTDOptions{Binary: filepath.Join(t.TempDir(), "missing-worker")})
	if err := missing.WithBackend(t.Context(), operation); !errors.Is(err, os.ErrNotExist) || called {
		t.Fatalf("scope lost executable startup failure: %v", err)
	}
	var unconfigured backend.ManagerCollection
	if err := unconfigured.WithCollection(t.Context(), "collection.anki2", operation); err == nil || called {
		t.Fatalf("unconfigured collection scope acquired a backend: %v", err)
	}
	collections := backend.NewManagerCollection(manager)
	if err := collections.WithCollection(ctx, "collection.anki2", operation); !errors.Is(err, context.Canceled) || !errors.Is(err, cause) || called {
		t.Fatalf("cancelled collection scope lost its cause: %v", err)
	}
	if err := collections.WithCollection(t.Context(), "", operation); err == nil || called {
		t.Fatalf("empty collection path reached a callback: %v", err)
	}
	if err := collections.WithCollection(t.Context(), "collection.anki2", nil); err == nil {
		t.Fatal("collection scope accepted a nil callback")
	}
	if err := collections.WithCollection(t.Context(), filepath.Join(t.TempDir(), "missing.colpkg"), operation); err == nil || called {
		t.Fatalf("missing native archive reached a callback: %v", err)
	}
	writeReport(t, "manager-failures.json", map[string]any{"checks": []string{"cancelled-before-startup", "nil-callback", "missing-executable", "missing-manager", "empty-path", "native-restore-failure"}})
}
