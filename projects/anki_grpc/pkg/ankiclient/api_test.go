//go:build cgo && linux && amd64

package ankiclient_test

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc"
	_ "git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankigrpc/native"
	ankibackend "github.com/ankitects/anki/proto/anki/backend"
	"github.com/ankitects/anki/proto/anki/collection"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/ankitects/anki/proto/anki/notes"
	"github.com/ankitects/anki/proto/anki/notetypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// workerBinary identifies the Bazel-built worker in the test runfiles.
var workerBinary = flag.String("worker", "", "Anki gRPC worker runfile")

// workerPath resolves the declared native executable for the E2E.
func workerPath(t *testing.T) string {
	t.Helper()
	path, err := runfiles.Rlocation(*workerBinary)
	if err != nil {
		t.Fatal(err)
	}
	return path
}

// workerFactoriesOnce installs worker schemes for this test executable's owned fixture.
var workerFactoriesOnce sync.Once

// registerWorkerFactories explicitly registers workers packaged by this E2E executable.
func registerWorkerFactories(t *testing.T) {
	t.Helper()
	workerFactoriesOnce.Do(func() {
		factory := ankigrpc.ClientFactoryFunc[*ankiclient.ClientSTD](ankigrpc.NewSubprocessClient)
		for _, scheme := range []string{"subprocess", "subprocess+cgo"} {
			if _, err := ankigrpc.NewClient(scheme + "://unconfigured-worker"); err == nil {
				t.Fatalf("library import registered an unowned worker: %s", scheme)
			}
			ankigrpc.RegisterClientFactory(scheme, factory)
		}
	})
}

// TestConfiguredFactories exercises private registration and defaults through real worker RPCs.
func TestConfiguredFactories(t *testing.T) {
	binary := workerPath(t)
	first := ankigrpc.NewFactory("private-worker://configured")
	second := ankigrpc.NewFactory("private-worker://configured")
	calls := []int{0, 0}
	for index, factory := range []*ankigrpc.Factory{first, second} {
		factory.RegisterClientFactory("private-worker", ankigrpc.ClientFactoryFunc[*ankiclient.ClientSTD](func(address string) (*ankiclient.ClientSTD, error) {
			if address != "private-worker://configured" {
				return nil, fmt.Errorf("private constructor received unexpected address %q", address)
			}
			calls[index]++
			client, err := ankigrpc.NewSubprocessClient(ankigrpc.SubprocessAddress(binary, "serve"))
			if err != nil {
				return nil, fmt.Errorf("configure private worker: %w", err)
			}
			return client, nil
		}))
	}
	for _, factory := range []*ankigrpc.Factory{first, second, first} {
		for _, address := range []string{"", "inprocess://explicit", ""} {
			if err := ankigrpc.WithClient(t.Context(), factory, address, func(connection grpc.ClientConnInterface) error {
				if _, err := collection.NewBackendCollectionServiceClient(connection).LatestProgress(t.Context(), &generic.Empty{}); err != nil {
					return fmt.Errorf("call configured factory worker: %w", err)
				}
				return nil
			}); err != nil {
				t.Fatal(err)
			}
		}
	}
	if calls[0] != 4 || calls[1] != 2 {
		t.Fatalf("factory defaults or registrations crossed scopes: %v", calls)
	}
	if _, err := ankigrpc.NewClient("private-worker://configured"); err == nil {
		t.Fatal("private registration changed the global factory")
	}
	if _, err := ankigrpc.NewFactory("").NewClient("private-worker://configured"); err == nil {
		t.Fatal("private registration changed a later factory")
	}
	writeReport(t, "configured-factories.json", map[string]any{"worker_calls": calls, "independent_registrations": true, "explicit_override": true, "default_restored": true})
}

// TestSubprocessArguments verifies owner registration, explicit executable selection and exact argv.
func TestSubprocessArguments(t *testing.T) {
	registerWorkerFactories(t)
	binary := workerPath(t)
	root := t.TempDir()
	if err := os.Symlink(binary, filepath.Join(root, "anki_grpc")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", root)
	query := url.Values{"binary": {binary}}
	for _, scheme := range []string{"subprocess", "subprocess+cgo"} {
		t.Run(scheme+"-zero-arguments", func(t *testing.T) {
			worker, err := ankigrpc.NewClient(scheme + "://?" + query.Encode())
			if err != nil {
				t.Fatal(err)
			}
			startupError := worker.Start(t.Context())
			if err := worker.Stop(context.Background()); err != nil {
				var exit *exec.ExitError
				if !errors.As(err, &exit) {
					t.Error(err)
				}
			}
			if startupError == nil {
				t.Fatal("zero URL arguments were replaced with a server command")
			}
		})
	}
	if _, err := ankigrpc.NewClient(""); err == nil {
		t.Fatal("empty library address selected an unconfigured worker")
	}
	if _, err := ankigrpc.NewSubprocessClient("subprocess://?arg=serve"); err == nil {
		t.Fatal("subprocess constructor guessed an executable")
	}
	worker, err := ankigrpc.NewClient(ankigrpc.SubprocessAddress(binary, "serve"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := worker.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	})
	if err := worker.Start(t.Context()); err != nil {
		t.Fatalf("explicit worker did not start its server command: %v", err)
	}
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := collection.NewBackendCollectionServiceClient(connection).LatestProgress(t.Context(), &generic.Empty{}); err != nil {
		t.Fatalf("explicit server command did not accept generated calls: %v", err)
	}
	writeReport(t, "subprocess-arguments.json", map[string]any{"zero_url_arguments": true, "qualified_zero_url_arguments": true, "explicit_worker_command": true, "owner_registration_required": true, "empty_library_address_rejected": true})
}

// TestNativeCollection verifies native persistence, card generation, and transport failures.
func TestNativeCollection(t *testing.T) {
	registerWorkerFactories(t)
	for _, transport := range []string{"stdio", "inprocess", "inprocess+cgo", "subprocess+cgo"} {
		t.Run(transport, func(t *testing.T) { nativeCollection(t, transport) })
	}
}

// nativeCollection exercises the same generated-client contract through each owned backend.
func nativeCollection(t *testing.T, transport string) {
	ctx := t.Context()
	empty, err := ankigrpc.NewClient(ankigrpc.SubprocessAddress(workerPath(t)))
	if err != nil {
		t.Fatal(err)
	}
	startupError := empty.Start(ctx)
	if err := empty.Stop(context.Background()); err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			t.Error(err)
		}
	}
	if startupError == nil {
		t.Fatal("an explicit empty argv was replaced with the serve default")
	}
	binary := filepath.Join(t.TempDir(), "worker # literal? name")
	if err := os.Symlink(workerPath(t), binary); err != nil {
		t.Fatal(err)
	}
	arguments := []string{"serve", "--transport", "stdio"}
	address := ankigrpc.SubprocessAddress(binary, arguments...)
	if strings.HasPrefix(transport, "inprocess") {
		address = transport + "://native-e2e"
	} else if transport == "subprocess+cgo" {
		address = strings.Replace(address, "subprocess://", "subprocess+cgo://", 1)
	}
	worker, err := ankigrpc.NewClient(address)
	if err != nil {
		t.Fatal(err)
	}
	arguments[1] = "invalid"
	for _, address := range []string{"subprocess+missing://worker", "subprocess://?shutdown_timeout=invalid", "subprocess://?unknown=value", "subprocess:///%zz", "subprocess+cgo:///%zz"} {
		if _, err := ankigrpc.NewClient(address); err == nil {
			t.Fatalf("unsupported or malformed address accepted: %s", address)
		}
	}
	if connection, err := worker.Connection(); err == nil || connection != nil {
		t.Fatal("constructor exposed a connection without an availability error")
	}
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := worker.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	})
	if err := worker.Start(ctx); err == nil {
		t.Fatal("second start accepted")
	}
	path := filepath.Join(t.TempDir(), "collection.anki2")
	open := &collection.OpenCollectionRequest{
		CollectionPath:  path,
		MediaFolderPath: path + ".media",
		MediaDbPath:     path + ".media.db",
	}
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	collections := collection.NewBackendCollectionServiceClient(connection)
	collectionOperations := collection.NewCollectionServiceClient(connection)
	models := notetypes.NewNotetypesServiceClient(connection)
	items := notes.NewNotesServiceClient(connection)
	if _, err := collections.OpenCollection(ctx, open); err != nil {
		t.Fatal(err)
	}
	if _, err := collectionOperations.LatestProgress(ctx, &generic.Empty{}); err != nil {
		t.Fatalf("overlapping collection/backend method: %v", err)
	}
	nt, err := models.GetNotetypeIdByName(ctx, &generic.String{Val: "Basic"})
	if err != nil {
		t.Fatal(err)
	}
	fields := []string{"<b>Hello</b>", "World"}
	added, err := items.AddNote(ctx, &notes.AddNoteRequest{
		Note: &notes.Note{
			NotetypeId: nt.Ntid,
			Fields:     fields,
		},
		DeckId: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if added.NoteId == 0 {
		t.Fatal("Anki did not allocate a note ID")
	}
	generated, err := items.CardsOfNote(ctx, &notes.NoteId{Nid: added.NoteId})
	if err != nil {
		t.Fatal(err)
	}
	if len(generated.Cids) != 1 {
		t.Fatalf("generated cards: %v", generated.Cids)
	}
	var trailers metadata.MD
	_, err = items.GetNote(ctx, &notes.NoteId{Nid: -1}, grpc.Trailer(&trailers))
	if status.Code(err) != codes.Unknown {
		t.Fatalf("missing note status: %v", err)
	}
	failure := &ankibackend.BackendError{}
	data := trailers.Get("anki-error-bin")
	if len(data) != 1 {
		t.Fatalf("caller error trailer: %v", trailers)
	}
	if decodeError := proto.Unmarshal([]byte(data[0]), failure); decodeError != nil || failure.Message == "" {
		t.Fatalf("native error trailer: %v, %v", failure, decodeError)
	}
	if !strings.Contains(err.Error(), failure.Message) {
		t.Fatalf("native failure detail missing from error: %v", err)
	}
	if err := connection.Invoke(ctx, "/anki.invalid.BackendInvalidService/Invalid", &generic.Empty{}, &generic.Empty{}); status.Code(err) != codes.Unimplemented {
		t.Fatal("unknown command accepted")
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := items.GetNote(cancelled, &notes.NoteId{Nid: added.NoteId}); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled call: %v", err)
	}
	if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	if _, err := collections.OpenCollection(ctx, open); err != nil {
		t.Fatal(err)
	}
	stored, err := items.GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(stored.Fields, fields) {
		t.Fatalf("persisted fields: %v", stored.Fields)
	}
	results := make(chan error, 8)
	for range 8 {
		go func() {
			note, err := items.GetNote(ctx, &notes.NoteId{Nid: added.NoteId})
			if err == nil && !reflect.DeepEqual(note.Fields, fields) {
				err = fmt.Errorf("concurrent note fields: %v", note.Fields)
			}
			results <- err
		}()
	}
	for range 8 {
		if err := <-results; err != nil {
			t.Fatal(err)
		}
	}
	if _, err := collections.CloseCollection(ctx, &collection.CloseCollectionRequest{}); err != nil {
		t.Fatal(err)
	}
	if err := worker.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := worker.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	if connection, err := worker.Connection(); err == nil || connection != nil {
		t.Fatal("stopped backend exposed a connection without an availability error")
	}
	if err := worker.Start(ctx); err == nil {
		t.Fatal("stopped backend restarted")
	}
	if _, err := items.GetNote(ctx, &notes.NoteId{Nid: added.NoteId}); err == nil {
		t.Fatal("closed backend accepted a command")
	}
	output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if output == "" {
		output = t.TempDir()
	}
	output = filepath.Join(output, transport)
	if err := os.MkdirAll(output, 0o700); err != nil {
		t.Fatal(err)
	}
	database, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "collection.anki2"), database, 0o600); err != nil {
		t.Fatal(err)
	}
	report, err := json.MarshalIndent(map[string]any{
		"note_id":   added.NoteId,
		"card_ids":  generated.Cids,
		"fields":    stored.Fields,
		"transport": transport,
		"checks": []string{
			"persisted-note", "generated-card", "native-error",
			"unknown-command", "cancelled-context", "closed-backend", "concurrent-reads",
			"generated-service-clients", "original-service-paths", "overlapping-service-method", "caller-error-trailer",
			"backend-interface", "explicit-start-stop", "repeated-start-rejected",
		},
	}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, "verification.json"), append(report, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestWorkerFailure replaces a PID-recording child with Anki to verify crash reaping.
func TestWorkerFailure(t *testing.T) {
	if os.Getenv("ANKI_E2E_RECORD_PID") == "1" {
		path := os.Getenv("ANKI_E2E_PID_FILE")
		if err := os.WriteFile(path, []byte(strconv.Itoa(os.Getpid())), 0o600); err != nil {
			t.Fatal(err)
		}
		executable := workerPath(t)
		if err := syscall.Exec(executable, []string{executable, "serve"}, os.Environ()); err != nil {
			t.Fatal(err)
		}
		return
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	worker, pid := startRecordedWorker(t, ctx)
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	items := notes.NewNotesServiceClient(connection)
	if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	if _, err := items.GetNote(ctx, &notes.NoteId{Nid: 1}); err == nil {
		t.Fatal("crashed worker accepted a call")
	}
	var exit *exec.ExitError
	if err := worker.Stop(t.Context()); !errors.As(err, &exit) {
		t.Fatalf("worker crash: %v", err)
	}
	assertReaped(t, exit)
	if exit.Pid() != pid {
		t.Fatalf("reaped PID %d, expected %d", exit.Pid(), pid)
	}
	if err := worker.Stop(t.Context()); !errors.As(err, &exit) {
		t.Fatalf("repeated stop lost crash: %v", err)
	}
	writeReport(t, "worker-failure.json", map[string]any{"checks": []string{"crash-call-error", "child-reaped", "repeated-stop-error"}, "exit": exit.ProcessState.String()})
}

// TestWorkerStartupDeadline runs an unresponsive child to verify startup cleanup.
func TestWorkerStartupDeadline(t *testing.T) {
	if os.Getenv("ANKI_E2E_UNRESPONSIVE_WORKER") == "1" {
		if err := os.WriteFile(os.Getenv("ANKI_E2E_STARTUP_MARKER"), []byte("started"), 0o600); err != nil {
			t.Fatal(err)
		}
		fmt.Fprintln(os.Stderr, "synthetic worker: no gRPC handshake")
		time.Sleep(time.Minute)
		return
	}
	ctx, cancel := context.WithTimeout(t.Context(), 250*time.Millisecond)
	defer cancel()
	t.Setenv("ANKI_E2E_UNRESPONSIVE_WORKER", "1")
	marker := filepath.Join(t.TempDir(), "startup-marker")
	t.Setenv("ANKI_E2E_STARTUP_MARKER", marker)
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{
		Binary:    os.Args[0],
		Arguments: []string{"-test.run=^TestWorkerStartupDeadline$"},
	})
	starting := make(chan error, 1)
	go func() { starting <- worker.Start(ctx) }()
	for {
		if _, err := os.Stat(marker); err == nil {
			break
		}
		if ctx.Err() != nil {
			t.Fatal("unresponsive worker did not start before its deadline")
		}
		time.Sleep(time.Millisecond)
	}
	if connection, err := worker.Connection(); err == nil || connection != nil {
		t.Fatal("Connection exposed an unready worker during startup")
	}
	err := <-starting
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("startup deadline: %v", err)
	}
	var exit *exec.ExitError
	if !errors.As(err, &exit) {
		t.Fatalf("startup cleanup lost worker exit: %v", err)
	}
	assertReaped(t, exit)
	if err := worker.Start(t.Context()); err == nil {
		t.Fatal("failed backend restarted")
	}
	writeReport(t, "startup-deadline.json", map[string]any{"checks": []string{"startup-deadline", "connection-unavailable-during-start", "child-reaped", "failed-start-not-retried"}, "exit": exit.ProcessState.String()})
}

// TestBackendLifecycle verifies no-op stop, invalid commands, and startup context isolation.
func TestBackendLifecycle(t *testing.T) {
	var unstarted ankiclient.Client = ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: workerPath(t), Arguments: []string{"serve"}})
	if err := unstarted.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := unstarted.Start(t.Context()); err == nil {
		t.Fatal("stopped configuration started a worker")
	}
	missing := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: filepath.Join(t.TempDir(), "missing-worker")})
	if err := missing.Start(t.Context()); err == nil {
		t.Fatal("missing executable accepted")
	}
	if err := missing.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	negative := ankiclient.NewSTD(ankiclient.ClientSTDOptions{
		Binary:          workerPath(t),
		ShutdownTimeout: -time.Second,
	})
	t.Cleanup(func() {
		if err := negative.Stop(context.Background()); err != nil {
			t.Error(err)
		}
	})
	if err := negative.Start(ctx); err == nil {
		t.Fatal("negative shutdown timeout accepted")
	}
	invalid := ankiclient.NewSTD(ankiclient.ClientSTDOptions{
		Binary:    workerPath(t),
		Arguments: []string{"--invalid-argument"},
	})
	var exit *exec.ExitError
	if err := invalid.Start(ctx); !errors.As(err, &exit) {
		t.Fatalf("invalid command arguments: %v", err)
	}
	assertReaped(t, exit)
	worker, pid := startRecordedWorker(t, ctx)
	connection, err := worker.Connection()
	if err != nil {
		t.Fatal(err)
	}
	collections := collection.NewBackendCollectionServiceClient(connection)
	cancel()
	if _, err := collections.LatestProgress(t.Context(), &generic.Empty{}); err != nil {
		t.Fatalf("startup context cancellation killed the ready worker: %v", err)
	}
	if err := syscall.Kill(pid, 0); err != nil {
		t.Fatalf("ready worker disappeared: %v", err)
	}
	if err := worker.Stop(t.Context()); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Kill(pid, 0); !errors.Is(err, syscall.ESRCH) {
		t.Fatalf("stopped worker PID still exists: %v", err)
	}
	writeReport(t, "lifecycle.json", map[string]any{"checks": []string{"stop-before-start", "missing-executable", "invalid-arguments", "invalid-shutdown-timeout", "startup-context-detached-after-readiness", "child-reaped"}})
}

// TestWorkerStopDeadline verifies forced shutdown after deadline, cancellation, and grace expiry.
func TestWorkerStopDeadline(t *testing.T) {
	if os.Getenv("ANKI_E2E_SLOW_STOP") == "1" {
		command := exec.CommandContext(t.Context(), workerPath(t), "serve")
		command.Stdin, command.Stdout, command.Stderr = os.Stdin, os.Stdout, os.Stderr
		if err := command.Run(); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(os.Getenv("ANKI_E2E_INNER_REAPED"), []byte("reaped"), 0o600); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Minute)
		return
	}
	t.Setenv("ANKI_E2E_SLOW_STOP", "1")
	for _, scenario := range []string{"deadline", "cancellation", "grace", "configured"} {
		t.Run(scenario, func(t *testing.T) {
			marker := filepath.Join(t.TempDir(), "inner-reaped")
			t.Setenv("ANKI_E2E_INNER_REAPED", marker)
			options := ankiclient.ClientSTDOptions{
				Binary:    os.Args[0],
				Arguments: []string{"-test.run=^TestWorkerStopDeadline$", "-worker=" + *workerBinary},
			}
			if scenario == "configured" {
				options.ShutdownTimeout = 500 * time.Millisecond
			}
			worker := ankiclient.NewSTD(options)
			if err := worker.Start(t.Context()); err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()
			var cause error
			if scenario == "deadline" {
				deadline, finish := context.WithTimeout(ctx, 300*time.Millisecond)
				defer finish()
				ctx, cause = deadline, context.DeadlineExceeded
			} else if scenario == "cancellation" {
				cause = context.Canceled
			}
			stopping := make(chan error, 1)
			go func() { stopping <- worker.Stop(ctx) }()
			observation, finish := context.WithTimeout(t.Context(), time.Second)
			defer finish()
			for {
				if observation.Err() != nil {
					t.Fatal("Connection failed to report shutdown promptly")
				}
				if connection, err := worker.Connection(); err != nil {
					if connection != nil {
						t.Fatal("shutdown returned a connection alongside its error")
					}
					break
				}
				time.Sleep(time.Millisecond)
			}
			if observation.Err() != nil || scenario == "deadline" && ctx.Err() != nil {
				t.Fatal("Connection waited for shutdown instead of returning an availability error")
			}
			for {
				if _, err := os.Stat(marker); err == nil {
					break
				}
				if observation.Err() != nil {
					t.Fatal("forwarding worker did not reap Anki after EOF")
				}
				time.Sleep(time.Millisecond)
			}
			if scenario == "cancellation" {
				cancel()
			}
			var err error
			if scenario == "configured" {
				select {
				case err = <-stopping:
				case <-observation.Done():
					t.Fatal("configured shutdown timeout was not applied")
				}
			} else {
				err = <-stopping
			}
			if cause != nil && !errors.Is(err, cause) {
				t.Fatalf("shutdown lost its context cause: %v", err)
			}
			if cause == nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
				t.Fatalf("grace expiry acquired a caller cancellation cause: %v", err)
			}
			var exit *exec.ExitError
			if !errors.As(err, &exit) {
				t.Fatalf("shutdown lost the reaped worker exit: %v", err)
			}
			assertReaped(t, exit)
			if repeat := worker.Stop(t.Context()); !errors.As(repeat, &exit) || cause != nil && !errors.Is(repeat, cause) {
				t.Fatalf("repeated Stop lost its cause: %v", repeat)
			}
			writeReport(t, "stop-"+scenario+".json", map[string]any{"checks": []string{"stop-context-" + scenario, "connection-unavailable-during-stop", "forced-termination", "inner-child-reaped", "child-reaped", "repeated-stop-cause"}})
		})
	}
}

// startRecordedWorker runs Anki through an exec replacement that records its PID.
func startRecordedWorker(t *testing.T, ctx context.Context) (*ankiclient.ClientSTD, int) {
	t.Helper()
	pidFile := filepath.Join(t.TempDir(), "worker.pid")
	t.Setenv("ANKI_E2E_RECORD_PID", "1")
	t.Setenv("ANKI_E2E_PID_FILE", pidFile)
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{
		Binary:    os.Args[0],
		Arguments: []string{"-test.run=^TestWorkerFailure$", "-worker=" + *workerBinary},
	})
	if _, err := os.Stat(pidFile); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("constructor started the child: %v", err)
	}
	if err := worker.Start(ctx); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatal(err)
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		t.Fatal(err)
	}
	return worker, pid
}

// assertReaped checks that the operating system no longer retains the child PID.
func assertReaped(t *testing.T, exit *exec.ExitError) {
	t.Helper()
	if err := syscall.Kill(exit.Pid(), 0); !errors.Is(err, syscall.ESRCH) {
		t.Fatalf("child PID %d still exists: %v", exit.Pid(), err)
	}
}

// writeReport emits repeatable E2E evidence alongside the persisted collection.
func writeReport(t *testing.T, name string, report any) {
	t.Helper()
	output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR")
	if output == "" {
		output = t.TempDir()
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(output, name), append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}
