package collection

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	_ "git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/native"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/export"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/ankitects/anki/proto/anki/media"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

// TestDirectStoreReservation retains source ownership while cancelled native work drains.
func TestDirectStoreReservation(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	_, archive := nativeSource(t)
	t.Setenv("TMPDIR", t.TempDir())
	store := New(Options{CloseTimeout: 25 * time.Millisecond})
	options := CollectionOptions{Path: archive, Address: "inprocess://reservation"}
	entered, release := make(chan struct{}), make(chan struct{})
	var releaseOnce sync.Once
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		close(entered)
		<-release
		if _, err := response.Write([]byte("reserved native media")); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	defer releaseOnce.Do(func() { close(release) })
	completed := make(chan error, 1)
	go func() {
		completed <- store.WithCollection(ctx, options, func(connection grpc.ClientConnInterface) error {
			download, cancelDownload := context.WithCancel(ctx)
			defer cancelDownload()
			result := make(chan error, 1)
			go func() {
				_, err := media.NewBackendMediaServiceClient(connection).AddMediaFromUrl(download, &media.AddMediaFromUrlRequest{Url: server.URL + "/reserved.txt"})
				result <- err
			}()
			select {
			case <-entered:
			case err := <-result:
				return fmt.Errorf("native download exited before fixture: %w", err)
			case <-ctx.Done():
				return fmt.Errorf("wait for native fixture: %w", ctx.Err())
			}
			cancelDownload()
			if err := <-result; err != nil {
				return fmt.Errorf("cancel native download: %w", err)
			}
			return nil
		})
	}()
	select {
	case err := <-completed:
		if !errors.Is(err, context.Canceled) || !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("scope did not retain cancellation and cleanup timeout: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("direct collection scope did not return after its cleanup budget")
	}
	err := store.WithCollection(ctx, options, func(grpc.ClientConnInterface) error { return nil })
	if err == nil || !strings.Contains(err.Error(), "reserve collection:") || !strings.Contains(err.Error(), "already open") {
		t.Fatalf("native work lost its source reservation before completion: %v", err)
	}
	releaseOnce.Do(func() { close(release) })
	for {
		err := store.WithCollection(ctx, options, func(grpc.ClientConnInterface) error { return nil })
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "already open") {
			t.Fatalf("reuse after native cleanup: %v", err)
		}
		select {
		case <-ctx.Done():
			t.Fatal("source remained reserved after native completion")
		case <-time.After(10 * time.Millisecond):
		}
	}
	emitArtifact(t, archive, "direct-reservation.colpkg")
}

// TestStoreRequiresExplicitWorker verifies missing configuration cannot discover a server or open a collection.
func TestStoreRequiresExplicitWorker(t *testing.T) {
	database, _ := nativeSource(t)
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", filepath.Dir(binary))
	original, err := os.ReadFile(database)
	if err != nil {
		t.Fatal(err)
	}
	store := &Store{}
	called := false
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	err = store.WithCollection(ctx, CollectionOptions{Path: database}, func(grpc.ClientConnInterface) error {
		called = true
		return nil
	})
	if err == nil || called {
		t.Fatalf("unconfigured store discovered a worker: error=%v callback=%t", err, called)
	}
	actual, err := os.ReadFile(database)
	if err != nil || string(actual) != string(original) {
		t.Fatalf("missing worker configuration changed the collection: %v", err)
	}
	store.closeTimeout = 10 * time.Minute
	store.newClient = fixtureFactory(ankigrpc.SubprocessAddress(binary, "serve") + "&shutdown_timeout=-1s")
	if _, err := readState(ctx, store, database, ""); err == nil {
		t.Fatalf("explicit shutdown policy was ignored: %v", err)
	}
	store.newClient = fixtureFactory(ankigrpc.SubprocessAddress(binary, "serve"))
	if _, err := readState(ctx, store, database, ""); err != nil {
		t.Fatalf("explicit worker could not reuse the source after failed startup: %v", err)
	}
	emitArtifact(t, database, "explicit-worker.anki2")
}

// fixtureFactory supplies the declared default while preserving explicit operation addresses.
func fixtureFactory(defaultAddress string) *ankigrpc.Factory {
	factory := ankigrpc.NewFactory(defaultAddress)
	factory.RegisterClientFactory("subprocess", ankigrpc.ClientFactoryFunc[*ankiclient.ClientSTD](ankigrpc.NewSubprocessClient))
	return factory
}

// fixtureStore supplies the explicitly declared native executable to test scopes.
func fixtureStore(t *testing.T) *Store {
	t.Helper()
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	store := New(Options{NewClient: fixtureFactory(ankigrpc.SubprocessAddress(binary, "serve"))})
	return &store
}

// TestStoreBackendLifecycle verifies native process ownership and source release on every scope exit.
func TestStoreBackendLifecycle(t *testing.T) {
	if os.Getenv("ANKI_STORE_RECORD_WORKER") == "1" {
		if err := os.WriteFile(os.Getenv("ANKI_STORE_PID_FILE"), []byte(strconv.Itoa(os.Getpid())), 0o600); err != nil {
			t.Fatal(err)
		}
		binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
		if err != nil {
			t.Fatal(err)
		}
		if err := syscall.Exec(binary, []string{binary, "serve"}, os.Environ()); err != nil {
			t.Fatal(err)
		}
		return
	}
	_, archive := nativeSource(t)
	pidFile := filepath.Join(t.TempDir(), "worker.pid")
	t.Setenv("TMPDIR", t.TempDir())
	t.Setenv("ANKI_STORE_RECORD_WORKER", "1")
	t.Setenv("ANKI_STORE_PID_FILE", pidFile)
	store := New(Options{NewClient: fixtureFactory(ankigrpc.SubprocessAddress(os.Args[0], "-test.run=^TestStoreBackendLifecycle$"))})
	for _, scenario := range []string{"success", "slow-publication", "error", "cancel", "panic", "crash"} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithCancelCause(t.Context())
			defer cancel(nil)
			failure := errors.New("collection callback failed")
			var recovered any
			var scopeError error
			var retained grpc.ClientConnInterface
			var pid int
			func() {
				defer func() { recovered = recover() }()
				scopeError = store.WithCollection(ctx, CollectionOptions{Path: archive, WriteBack: true}, func(connection grpc.ClientConnInterface) error {
					retained = connection
					pid = recordedWorkerPID(t, pidFile)
					if _, err := state.New(connection, "").Dump(ctx); err != nil {
						return fmt.Errorf("read prepared native collection: %w", err)
					}
					switch scenario {
					case "slow-publication":
						value := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"slow-publication": {Id: 1}}}}
						if _, err := state.New(connection, "").Update(ctx, value); err != nil {
							return fmt.Errorf("update delayed archive: %w", err)
						}
						if err := syscall.Kill(pid, syscall.SIGSTOP); err != nil {
							return fmt.Errorf("pause native publication: %w", err)
						}
						resume := time.AfterFunc(6*time.Second, func() { _ = syscall.Kill(pid, syscall.SIGCONT) })
						t.Cleanup(func() {
							resume.Stop()
						})
					case "error":
						return failure
					case "cancel":
						cancel(failure)
					case "panic":
						panic(failure)
					case "crash":
						if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
							return fmt.Errorf("terminate fixture worker: %w", errors.Join(failure, err))
						}
						return failure
					}
					return nil
				})
			}()
			switch scenario {
			case "success", "slow-publication":
				if scopeError != nil || recovered != nil {
					t.Fatalf("native collection failed: %v %v", scopeError, recovered)
				}
			case "panic":
				if recovered != failure {
					t.Fatalf("collection scope lost its panic: %v", recovered)
				}
			default:
				if !errors.Is(scopeError, failure) || scenario == "cancel" && !errors.Is(scopeError, context.Canceled) {
					t.Fatalf("scope lost callback/cancellation errors: %v", scopeError)
				}
				if scenario == "crash" {
					var exit *exec.ExitError
					if !errors.As(scopeError, &exit) || exit.Pid() != pid {
						t.Fatalf("scope lost worker crash: %v", scopeError)
					}
				}
			}
			assertWorkerReaped(t, pid)
			if _, err := state.New(retained, "").Dump(t.Context()); err == nil {
				t.Fatal("completed scope retained a prepared connection")
			}
			if _, err := readState(t.Context(), &store, archive, ""); err != nil {
				t.Fatalf("scope retained its source reservation: %v", err)
			}
			if scenario == "slow-publication" {
				value, err := readState(t.Context(), &store, archive, "")
				if err != nil || value.GetResources().GetPresets()["slow-publication"].GetId() != 1 {
					t.Fatalf("delayed publication lost native state: %v %v", value, err)
				}
			}
			assertWorkerReaped(t, recordedWorkerPID(t, pidFile))
		})
	}
	store.newClient = fixtureFactory(ankigrpc.SubprocessAddress(filepath.Join(t.TempDir(), "missing-worker"), "-test.run=^TestStoreBackendLifecycle$"))
	if _, err := readState(t.Context(), &store, archive, ""); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("startup failure lost its cause: %v", err)
	}
	store.newClient = fixtureFactory(ankigrpc.SubprocessAddress(os.Args[0], "-test.run=^TestStoreBackendLifecycle$"))
	if _, err := readState(t.Context(), &store, archive, ""); err != nil {
		t.Fatalf("startup failure retained its source reservation: %v", err)
	}
	assertWorkerReaped(t, recordedWorkerPID(t, pidFile))
	emitArtifact(t, archive, "direct-ankiclient.colpkg")
}

// recordedWorkerPID reads the subprocess fixture's native worker identity after readiness.
func recordedWorkerPID(t *testing.T, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil || pid <= 0 {
		t.Fatalf("invalid worker PID: %q %v", data, err)
	}
	return pid
}

// assertWorkerReaped verifies the scope waited for its owned child and released its process.
func assertWorkerReaped(t *testing.T, pid int) {
	t.Helper()
	if err := syscall.Kill(pid, 0); !errors.Is(err, syscall.ESRCH) {
		t.Fatalf("scope retained worker PID %d: %v", pid, err)
	}
	var status syscall.WaitStatus
	if _, err := syscall.Wait4(pid, &status, syscall.WNOHANG, nil); !errors.Is(err, syscall.ECHILD) {
		t.Fatalf("scope did not reap worker PID %d: %v", pid, err)
	}
}

// nativeSource creates one representative note and its archive entirely through Anki.
func nativeSource(t *testing.T) (string, string) {
	t.Helper()
	database := filepath.Join(t.TempDir(), "collection.anki2")
	archive := filepath.Join(filepath.Dir(database), "source.colpkg")
	if err := os.MkdirAll(filepath.Join(filepath.Dir(database), "collection.media"), 0o755); err != nil {
		t.Fatal(err)
	}
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		t.Fatal(err)
	}
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: binary, Arguments: []string{"serve"}})
	prepared := ankiclient.NewCollection(ankiclient.ClientCollectionOptions{Backend: worker, Path: database})
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := prepared.Stop(cleanup); err != nil {
			t.Error(err)
		}
		if err := worker.Stop(cleanup); err != nil {
			t.Error(err)
		}
	}()
	ctx := t.Context()
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	if err := prepared.Start(ctx); err != nil {
		t.Fatal(err)
	}
	connection, err := prepared.Connection()
	if err != nil {
		t.Fatal(err)
	}
	definition, err := notetypes.NewNotetypesServiceClient(connection).AddNotetype(ctx, &notetypes.Notetype{
		Name: "Example", Config: &notetypes.Notetype_Config{},
		Fields: []*notetypes.Notetype_Field{
			{Name: "Front", Config: &notetypes.Notetype_Field_Config{}},
			{Name: "Back", Config: &notetypes.Notetype_Field_Config{}},
		},
		Templates: []*notetypes.Notetype_Template{{Name: "Card", Config: &notetypes.Notetype_Template_Config{QFormat: "{{Front}}", AFormat: "{{Back}}"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := notes.NewNotesServiceClient(connection).AddNote(ctx, &notes.AddNoteRequest{
		DeckId: 1, Note: &notes.Note{Guid: "store-note", NotetypeId: definition.Id, Fields: []string{"Question", "Answer"}},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: archive, IncludeMedia: true}); err != nil {
		t.Fatal(err)
	}
	return database, archive
}

// readState observes persisted state through a native collection scope.
func readState(ctx context.Context, store *Store, path, key string) (*ankiascodepb.State, error) {
	var value *ankiascodepb.State
	if err := store.WithCollection(ctx, CollectionOptions{Path: path, WriteBack: true}, func(connection grpc.ClientConnInterface) error {
		var err error
		value, err = state.New(connection, key).Dump(ctx)
		if err != nil {
			return fmt.Errorf("read native state: %w", err)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("open collection for state: %w", err)
	}
	return value, nil
}

// emitArtifact retains a repeatable native collection or exported config for independent inspection.
func emitArtifact(t *testing.T, path, name string) {
	t.Helper()
	root := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if root == "" {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, name), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestStoreServices verifies local service reuse, namespace isolation, config export and persistence.
func TestStoreServices(t *testing.T) {
	for _, format := range []string{"database", "package"} {
		t.Run(format, func(t *testing.T) {
			database, archive := nativeSource(t)
			input, artifact := archive, "scoped-package.colpkg"
			if format == "database" {
				input, artifact = database, "scoped-database.anki2"
			}
			original, err := os.ReadFile(input)
			if err != nil {
				t.Fatal(err)
			}
			scratch, output := t.TempDir(), t.TempDir()
			t.Setenv("TMPDIR", scratch)
			store := fixtureStore(t)
			ctx := t.Context()
			key := "scoped-state"
			expected := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"first": {Id: 1}}}}
			var retained grpc.ClientConnInterface
			if err := store.WithCollection(ctx, CollectionOptions{Path: input, WriteBack: true}, func(connection grpc.ClientConnInterface) error {
				retained = connection
				resources, states := resource.New(connection), state.New(connection, key)
				exporter := export.New(resources, states)
				if _, err := states.Replace(ctx, expected); err != nil {
					return fmt.Errorf("replace scoped state: %w", err)
				}
				update := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"second": {Id: 1}}}}
				if _, err := states.Update(ctx, update); err != nil {
					return fmt.Errorf("update shared state: %w", err)
				}
				if err := exporter.Export(ctx, output, ""); err != nil {
					return fmt.Errorf("export collection services: %w", err)
				}
				expected.Resources.Presets["second"] = &ankiascodepb.PresetState{Id: 1}
				actual, err := states.Dump(ctx)
				if err != nil {
					return fmt.Errorf("read shared changes: %w", err)
				}
				if !proto.Equal(actual, expected) {
					return fmt.Errorf("shared state lost earlier operations")
				}
				other, err := state.New(connection, "independent-state").Dump(ctx)
				if err != nil {
					return fmt.Errorf("read independent state: %w", err)
				}
				if len(other.GetResources().GetPresets()) != 0 {
					return fmt.Errorf("namespaces share stored bindings")
				}
				if format == "package" {
					current, err := os.ReadFile(input)
					if err != nil {
						return fmt.Errorf("read pending archive: %w", err)
					}
					if string(current) != string(original) {
						return fmt.Errorf("archive was published before scope success")
					}
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
			if _, err := state.New(retained, key).Dump(ctx); err == nil {
				t.Fatal("closed collection retained state access")
			}
			actual, err := readState(ctx, store, input, key)
			if err != nil || !proto.Equal(actual, expected) {
				t.Fatalf("state did not survive scope cleanup: value=%v error=%v", actual, err)
			}
			if _, err := config.New(output).Load(); err != nil {
				t.Fatalf("exporter produced invalid configs: %v", err)
			}
			files, err := os.ReadDir(scratch)
			if err != nil || len(files) != 0 {
				t.Fatalf("scope retained temporary files: %v %v", files, err)
			}
			emitArtifact(t, input, artifact)
			emitArtifact(t, filepath.Join(output, "collection.anki.toml"), "exported-"+format+".anki.toml")
		})
	}
}

// TestStoreReservations verifies duplicate source aliases fail and completed scopes release reservations.
func TestStoreReservations(t *testing.T) {
	_, archive := nativeSource(t)
	store := fixtureStore(t)
	ctx := t.Context()
	symlink, hardlink := archive+".symlink.colpkg", archive+".hardlink.colpkg"
	if err := os.Symlink(archive, symlink); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(archive, hardlink); err != nil {
		t.Fatal(err)
	}
	if err := store.WithCollection(ctx, CollectionOptions{Path: archive, WriteBack: true}, func(connection grpc.ClientConnInterface) error {
		for _, path := range []string{archive, symlink, hardlink} {
			if _, err := readState(ctx, store, path, ""); err == nil {
				return fmt.Errorf("reserved source alias %q was reopened", path)
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := readState(ctx, store, hardlink, ""); err != nil {
		t.Fatalf("completed scope retained its source reservation: %v", err)
	}
}

// TestStoreFailurePublication verifies request errors and cancellation do not suppress write-back.
func TestStoreFailurePublication(t *testing.T) {
	for _, cancelled := range []bool{false, true} {
		_, archive := nativeSource(t)
		store := fixtureStore(t)
		ctx, cancel := context.WithCancel(t.Context())
		failure := errors.New("scope failed after native mutation")
		err := store.WithCollection(ctx, CollectionOptions{Path: archive, WriteBack: true}, func(connection grpc.ClientConnInterface) error {
			value := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"discarded": {Id: 1}}}}
			if _, err := state.New(connection, "").Replace(ctx, value); err != nil {
				return fmt.Errorf("mutate restored collection: %w", err)
			}
			if cancelled {
				cancel()
				return nil
			}
			return failure
		})
		cancel()
		cause := failure
		if cancelled {
			cause = context.Canceled
		}
		if !errors.Is(err, cause) {
			t.Fatalf("scope lost failure cause: %v", err)
		}
		actual, err := readState(t.Context(), store, archive, "")
		if err != nil || actual.GetResources().GetPresets()["discarded"].GetId() != 1 {
			t.Fatalf("failed scope lost committed data or retained its reservation: %v", err)
		}
	}
}

// TestStoreExplicitExport verifies an explicit export and source write-back preserve the same changes.
func TestStoreExplicitExport(t *testing.T) {
	_, archive := nativeSource(t)
	store := fixtureStore(t)
	ctx := t.Context()
	output := filepath.Join(t.TempDir(), "result.colpkg")
	if err := store.WithCollection(ctx, CollectionOptions{Path: archive, WriteBack: true}, func(connection grpc.ClientConnInterface) error {
		value := &ankiascodepb.State{Version: 1, Resources: &ankiascodepb.Resources{Presets: map[string]*ankiascodepb.PresetState{"output": {Id: 1}}}}
		if _, err := state.New(connection, "").Replace(ctx, value); err != nil {
			return fmt.Errorf("write explicit output state: %w", err)
		}
		if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{OutPath: output, IncludeMedia: true}); err != nil {
			return fmt.Errorf("export caller-selected package: %w", err)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	actual, err := readState(ctx, store, archive, "")
	if err != nil || actual.GetResources().GetPresets()["output"].GetId() != 1 {
		t.Fatalf("explicit output lost source write-back: %v", err)
	}
	value, err := readState(ctx, store, output, "")
	if err != nil || value.GetResources().GetPresets()["output"].GetId() != 1 {
		t.Fatalf("explicit output lost native state: value=%v error=%v", value, err)
	}
	emitArtifact(t, output, "explicit-output.colpkg")
}
