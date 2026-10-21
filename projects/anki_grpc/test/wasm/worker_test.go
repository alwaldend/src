package wasm

import (
	"bufio"
	"bytes"
	"context"
	"debug/elf"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/ankiclient"
	"github.com/ankitects/anki/proto/anki/backend"
	"github.com/ankitects/anki/proto/anki/generic"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// workerRunfile selects the cgo-disabled executable containing the guest fixture.
var workerRunfile = flag.String("worker", "", "WASM fixture worker runfile")

// nativeWorkerRunfile selects the native CLI used to verify shared Rust command behavior.
var nativeWorkerRunfile = flag.String("native-worker", "", "native CLI parity fixture")

// workerExecutable copies only the executable to prove that no companion guest is needed.
func workerExecutable(t *testing.T) string {
	t.Helper()
	path, err := runfiles.Rlocation(*workerRunfile)
	if err != nil {
		t.Fatalf("resolve fixture worker: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture worker: %v", err)
	}
	output := filepath.Join(t.TempDir(), "anki_grpc")
	if err := os.WriteFile(output, data, 0o700); err != nil {
		t.Fatalf("copy standalone worker: %v", err)
	}
	return output
}

// TestWorkerProcess exercises real gRPC over child-process pipes and native error trailers.
func TestWorkerProcess(t *testing.T) {
	path := workerExecutable(t)
	binary, err := elf.Open(path)
	if err != nil {
		t.Fatalf("open worker ELF: %v", err)
	}
	libraries, libraryErr := binary.ImportedLibraries()
	closeErr := binary.Close()
	if libraryErr != nil || closeErr != nil || len(libraries) != 0 {
		t.Fatalf("worker has native shared-library requirements: %v; %v; %v", libraries, libraryErr, closeErr)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: path, Arguments: []string{"serve"}})
	if err := worker.Start(ctx); err != nil {
		t.Fatalf("start WASM worker: %v", err)
	}
	t.Cleanup(func() {
		stopContext, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer stopCancel()
		if err := worker.Stop(stopContext); err != nil {
			t.Errorf("reap WASM worker: %v", err)
		}
	})
	connection, err := worker.Connection()
	if err != nil {
		t.Fatalf("get WASM connection: %v", err)
	}
	request := &generic.String{Val: "protobuf through an isolated Rust WASI guest"}
	response := &generic.String{}
	if err := connection.Invoke(ctx, "/fixture/Echo", request, response); err != nil || response.Val != request.Val {
		t.Fatalf("WASM worker protobuf exchange: %v; %q", err, response.Val)
	}
	details := &backend.BackendError{Message: "fixture native failure"}
	want, err := proto.Marshal(details)
	if err != nil {
		t.Fatalf("encode native error fixture: %v", err)
	}
	var trailers metadata.MD
	err = connection.Invoke(ctx, "/fixture/NativeError", details, &generic.Empty{}, grpc.Trailer(&trailers))
	if status.Code(err) != codes.Unknown || !strings.Contains(err.Error(), details.Message) {
		t.Fatalf("native error status/details were lost: %v", err)
	}
	values := trailers.Get("anki-error-bin")
	if len(values) != 1 || !bytes.Equal([]byte(values[0]), want) {
		t.Fatalf("native error trailer changed: %x", values)
	}
	if err := connection.Invoke(ctx, "/fixture/Missing", &generic.Empty{}, &generic.Empty{}); status.Code(err) != codes.Unimplemented {
		t.Fatalf("unknown method status: %v", err)
	}
	if err := worker.Stop(ctx); err != nil {
		t.Fatalf("graceful WASM worker exit: %v", err)
	}
	if _, err := worker.Connection(); err == nil {
		t.Fatal("stopped worker remained ready")
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		data, err := json.MarshalIndent(map[string]any{"protobuf_echo": response.Val, "native_error_preserved": true, "standalone_executable": true, "native_shared_libraries": libraries}, "", "  ")
		if err != nil {
			t.Fatalf("encode worker result: %v", err)
		}
		if err := os.WriteFile(filepath.Join(output, "wasm-worker.json"), data, 0o600); err != nil {
			t.Fatalf("publish worker artifact: %v", err)
		}
	}
}

// TestWorkerCompilationCache exercises reuse across processes and explicit cache failures.
func TestWorkerCompilationCache(t *testing.T) {
	path := workerExecutable(t)
	cache := filepath.Join(t.TempDir(), "compiled")
	t.Setenv("ANKI_WASM_CACHE_DIR", cache)
	for range 2 {
		ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
		defer cancel()
		worker := ankiclient.NewSTD(ankiclient.ClientSTDOptions{Binary: path, Arguments: []string{"serve"}})
		t.Cleanup(func() {
			stop, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer stopCancel()
			if err := worker.Stop(stop); err != nil {
				t.Errorf("reap cached worker: %v", err)
			}
		})
		if err := worker.Start(ctx); err != nil {
			cancel()
			t.Fatalf("start cached worker: %v", err)
		}
		connection, err := worker.Connection()
		if err != nil {
			t.Fatalf("get cached worker connection: %v", err)
		}
		response := &generic.String{}
		if err := connection.Invoke(ctx, "/fixture/Echo", &generic.String{Val: "cached guest"}, response); err != nil || response.Val != "cached guest" {
			t.Fatalf("cached worker exchange: %q; %v", response.Val, err)
		}
		if err := worker.Stop(ctx); err != nil {
			t.Fatalf("stop cached worker: %v", err)
		}
		cancel()
	}
	blocked := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(blocked, nil, 0o600); err != nil {
		t.Fatalf("prepare invalid cache: %v", err)
	}
	t.Setenv("ANKI_WASM_CACHE_DIR", blocked)
	var output, diagnostics bytes.Buffer
	command := exec.CommandContext(t.Context(), path, "serve")
	command.Stdout, command.Stderr = &output, &diagnostics
	if err := command.Run(); err == nil || output.Len() != 0 || !strings.Contains(diagnostics.String(), "compilation cache") {
		t.Fatalf("invalid cache was not rejected before serving: %v; %q; %q", err, output.String(), diagnostics.String())
	}
}

// TestWorkerCLI compares installed WASM help and validation with the native Rust CLI.
func TestWorkerCLI(t *testing.T) {
	path := workerExecutable(t)
	native, err := runfiles.Rlocation(*nativeWorkerRunfile)
	if err != nil {
		t.Fatalf("resolve native CLI: %v", err)
	}
	unavailable := filepath.Join(t.TempDir(), "unavailable")
	t.Setenv("TMPDIR", unavailable)
	t.Setenv("ANKI_WASM_CACHE_DIR", filepath.Join(unavailable, "cache"))
	t.Setenv("PATH", "")
	var report strings.Builder
	for _, arguments := range [][]string{
		nil,
		{"--help"},
		{"serve", "--help"},
		{"unknown"},
		{"serve", "--transport", "invalid"},
		{"serve", "--transport", "tcp"},
		{"serve", "--listen", "unused"},
		{"serve", "--transport", "stdio", "--transport", "stdio"},
		{"serve", "--listen", "first", "--listen", "second"},
		{"serve", "--transport=stdio"},
		{"serve", "--listen", "--transport", "tcp"},
	} {
		var expected, expectedDiagnostics bytes.Buffer
		reference := exec.CommandContext(t.Context(), native, arguments...)
		reference.Stdout, reference.Stderr = &expected, &expectedDiagnostics
		referenceErr := reference.Run()
		var output, diagnostics bytes.Buffer
		command := exec.CommandContext(t.Context(), path, arguments...)
		command.Stdout, command.Stderr = &output, &diagnostics
		err := command.Run()
		if (err == nil) != (referenceErr == nil) || !bytes.Equal(output.Bytes(), expected.Bytes()) {
			t.Fatalf("CLI parity %v: native=%v stdout=%q stderr=%q; WASM=%v stdout=%q stderr=%q", arguments, referenceErr, expected.String(), expectedDiagnostics.String(), err, output.String(), diagnostics.String())
		}
		if command.ProcessState.ExitCode() != reference.ProcessState.ExitCode() {
			t.Fatalf("CLI exit status %v: native=%d WASM=%d", arguments, reference.ProcessState.ExitCode(), command.ProcessState.ExitCode())
		}
		if referenceErr != nil {
			reason := strings.SplitN(expectedDiagnostics.String(), "\n", 2)[0]
			if reason == "" || !strings.Contains(diagnostics.String(), reason) {
				t.Fatalf("Rust CLI validation %v: want %q; stderr=%q", arguments, reason, diagnostics.String())
			}
		}
		fmt.Fprintf(&report, "%q: exit=%d\n%s%s", arguments, command.ProcessState.ExitCode(), output.String(), diagnostics.String())
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "rust-cli-parity.txt"), []byte(report.String()), 0o600); err != nil {
			t.Fatalf("write CLI parity artifact: %v", err)
		}
	}
}

// TestWorkerNetwork verifies that the same WASM dispatch is reachable over TCP and Unix sockets.
func TestWorkerNetwork(t *testing.T) {
	path := workerExecutable(t)
	for _, transport := range []string{"tcp", "unix"} {
		t.Run(transport, func(t *testing.T) {
			t.Chdir(t.TempDir())
			ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
			defer cancel()
			listen := "127.0.0.1:0"
			if transport == "unix" {
				listen = "worker.sock"
			}
			command := exec.CommandContext(ctx, path, "serve", "--transport", transport, "--listen", listen)
			command.Stdout = io.Discard
			stderr, err := command.StderrPipe()
			if err != nil {
				t.Fatalf("connect worker diagnostics: %v", err)
			}
			if err := command.Start(); err != nil {
				t.Fatalf("start network worker: %v", err)
			}
			t.Cleanup(func() { _ = command.Process.Kill(); _ = command.Wait() })
			ready := make(chan string, 1)
			go func() {
				scanner := bufio.NewScanner(stderr)
				if scanner.Scan() {
					ready <- scanner.Text()
				} else {
					ready <- ""
				}
				for scanner.Scan() {
				}
			}()
			var address string
			select {
			case line := <-ready:
				prefix := "Listening on " + transport + "://"
				if !strings.HasPrefix(line, prefix) {
					t.Fatalf("worker startup: %q", line)
				}
				address = strings.TrimPrefix(line, prefix)
			case <-ctx.Done():
				t.Fatal(ctx.Err())
			}
			connection, err := grpc.NewClient("passthrough:///fixture", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
				connection, err := (&net.Dialer{}).DialContext(ctx, transport, address)
				if err != nil {
					return nil, fmt.Errorf("dial fixture %s endpoint %q: %w", transport, address, err)
				}
				return connection, nil
			}))
			if err != nil {
				t.Fatalf("create network client: %v", err)
			}
			t.Cleanup(func() { _ = connection.Close() })
			response := &generic.String{}
			if err := connection.Invoke(ctx, "/fixture/Echo", &generic.String{Val: transport}, response); err != nil || response.Val != transport {
				t.Fatalf("network guest exchange: %v; %q", err, response.Val)
			}
		})
	}
}

// TestWorkerReconnect keeps one Rust backend alive across independent network clients.
func TestWorkerReconnect(t *testing.T) {
	path := workerExecutable(t)
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, path, "serve", "--transport", "tcp", "--listen", "127.0.0.1:0")
	diagnostics, err := command.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	command.Stdout = io.Discard
	if err := command.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = command.Process.Kill(); _ = command.Wait() })
	ready := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(diagnostics)
		if scanner.Scan() {
			ready <- scanner.Text()
		}
		for scanner.Scan() {
		}
	}()
	var address string
	select {
	case line := <-ready:
		address = strings.TrimPrefix(line, "Listening on tcp://")
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	for _, expected := range []string{"1", "2"} {
		connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			t.Fatal(err)
		}
		response := &generic.String{}
		invokeErr := connection.Invoke(ctx, "/fixture/Count", &generic.Empty{}, response)
		closeErr := connection.Close()
		if invokeErr != nil || closeErr != nil || response.Val != expected {
			t.Fatalf("backend state after reconnect: value=%q expected=%q call=%v close=%v", response.Val, expected, invokeErr, closeErr)
		}
	}
}

// TestWorkerPrematureEOF rejects a stream that never completes the HTTP/2 handshake.
func TestWorkerPrematureEOF(t *testing.T) {
	command := exec.CommandContext(t.Context(), workerExecutable(t), "serve")
	var diagnostics bytes.Buffer
	command.Stdout, command.Stderr = io.Discard, &diagnostics
	if err := command.Run(); err == nil {
		t.Fatalf("premature protocol EOF succeeded: %s", diagnostics.String())
	}
}

// TestWorkerProtocolFailure exits after rejecting input even while the peer retains stdin.
func TestWorkerProtocolFailure(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, workerExecutable(t), "serve")
	var diagnostics bytes.Buffer
	command.Stdout, command.Stderr = io.Discard, &diagnostics
	input, err := command.StdinPipe()
	if err != nil {
		t.Fatalf("prepare worker input: %v", err)
	}
	defer func() {
		if err := input.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
			t.Errorf("close worker input: %v", err)
		}
	}()
	if err := command.Start(); err != nil {
		t.Fatalf("start worker with malformed protocol input: %v", err)
	}
	if _, err := io.WriteString(input, "GET / HTTP/1.1\r\n\r\ninvalid"); err != nil {
		cancel()
		_ = command.Wait()
		t.Fatalf("send malformed protocol preface: %v", err)
	}
	err = command.Wait()
	if ctx.Err() != nil || err == nil || !strings.Contains(diagnostics.String(), "PROTOCOL_ERROR") {
		t.Fatalf("worker did not reject protocol while stdin remained open: exit=%v context=%v diagnostics=%s", err, ctx.Err(), diagnostics.String())
	}
	if output := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); output != "" {
		if err := os.WriteFile(filepath.Join(output, "rust-command-protocol-failure.txt"), []byte(diagnostics.String()), 0o600); err != nil {
			t.Fatalf("publish protocol rejection artifact: %v", err)
		}
	}
}
