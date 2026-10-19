package wasmhost

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// guestPath selects the actual Rust WASI command in Bazel runfiles.
var guestPath = flag.String("guest", "", "Rust WASI command fixture runfile")

// guestFixture reads the declared command artifact rather than a host mock.
func guestFixture(t *testing.T) []byte {
	t.Helper()
	path, err := runfiles.Rlocation(*guestPath)
	if err != nil {
		t.Fatalf("resolve guest fixture: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read guest fixture: %v", err)
	}
	return data
}

// commandConnection executes the Rust CLI and connects a real gRPC client to its pipes.
func commandConnection(t *testing.T) (*grpc.ClientConn, context.CancelCauseFunc, <-chan error) {
	t.Helper()
	ctx, cancel := context.WithCancelCause(t.Context())
	client, server := net.Pipe()
	done := make(chan error, 1)
	guest := guestFixture(t)
	go func() { done <- RunCommand(ctx, []string{"anki_grpc", "serve"}, guest, server, server, io.Discard) }()
	var mutex sync.Mutex
	available := true
	connection, err := grpc.NewClient("passthrough:///rust-cli", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		mutex.Lock()
		defer mutex.Unlock()
		if !available {
			return nil, errors.New("Rust command connection already consumed")
		}
		available = false
		return client, nil
	}))
	if err != nil {
		cancel(err)
		t.Fatalf("create command gRPC client: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close(); cancel(context.Canceled); _ = client.Close(); _ = server.Close() })
	return connection, cancel, done
}

// TestCommandExchange checks large protobuf streaming through Rust's gRPC implementation.
func TestCommandExchange(t *testing.T) {
	connection, _, _ := commandConnection(t)
	request := wrapperspb.String(strings.Repeat("Rust CLI protobuf exchange ", 80000))
	response := &wrapperspb.StringValue{}
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	if err := connection.Invoke(ctx, "/fixture/Echo", request, response); err != nil || response.Value != request.Value {
		t.Fatalf("large Rust command exchange: %v; bytes=%d", err, len(response.Value))
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		if err := os.WriteFile(filepath.Join(directory, "rust-command-exchange.txt"), []byte(fmt.Sprintf("opaque protobuf payload bytes=%d; Rust gRPC round trip passed\n", len(response.Value))), 0o600); err != nil {
			t.Fatalf("retain command artifact: %v", err)
		}
	}
}

// TestCommandCancellation interrupts CPU-bound Rust and confirms that it is not restarted.
func TestCommandCancellation(t *testing.T) {
	connection, cancel, done := commandConnection(t)
	ctx, stop := context.WithTimeout(t.Context(), 20*time.Second)
	defer stop()
	if err := connection.Invoke(ctx, "/fixture/Echo", wrapperspb.String("ready"), &wrapperspb.StringValue{}); err != nil {
		t.Fatalf("initialize command connection: %v", err)
	}
	result := make(chan error, 1)
	go func() {
		result <- connection.Invoke(ctx, "/fixture/Loop", wrapperspb.String("loop"), &wrapperspb.StringValue{})
	}()
	timer := time.AfterFunc(100*time.Millisecond, func() { cancel(context.DeadlineExceeded) })
	defer timer.Stop()
	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("command cancellation cause: %v", err)
		}
	case <-ctx.Done():
		t.Fatal("CPU-bound command did not stop")
	}
	if err := <-result; err == nil {
		t.Fatal("cancelled command RPC succeeded")
	}
	retry, retryCancel := context.WithTimeout(t.Context(), time.Second)
	defer retryCancel()
	if err := connection.Invoke(retry, "/fixture/Echo", wrapperspb.String("no replay"), &wrapperspb.StringValue{}); err == nil {
		t.Fatal("cancelled command was restarted")
	}
}

// TestCommandFailures checks actual malformed modules and Rust traps at the command boundary.
func TestCommandFailures(t *testing.T) {
	for _, data := range [][]byte{[]byte("not wasm"), {0, 97, 115, 109, 1, 0, 0, 0}} {
		if err := RunCommand(t.Context(), []string{"anki_grpc", "serve"}, data, io.NopCloser(bytes.NewReader(nil)), discardCloser{}, io.Discard); err == nil {
			t.Fatal("accepted a guest without a WASI command")
		}
	}
	connection, _, done := commandConnection(t)
	ctx, cancel := context.WithTimeout(t.Context(), 20*time.Second)
	defer cancel()
	if err := connection.Invoke(ctx, "/fixture/Trap", wrapperspb.String("trap"), &wrapperspb.StringValue{}); err == nil {
		t.Fatal("trapped command RPC succeeded")
	}
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Rust trap was accepted")
		}
	case <-ctx.Done():
		t.Fatal("trapped command did not exit")
	}
}

// discardCloser supplies disposable command output for malformed-module cases.
type discardCloser struct{}

// Write accepts bytes without retaining malformed-command output.
func (discardCloser) Write(data []byte) (int, error) { return len(data), nil }

// Close releases the output without external resources.
func (discardCloser) Close() error { return nil }

// TestWASICommandCLI executes help as a WASI command without backend or transport inputs.
func TestWASICommandCLI(t *testing.T) {
	ctx := t.Context()
	runtime := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfigInterpreter())
	defer func() {
		if err := runtime.Close(ctx); err != nil {
			t.Errorf("close command runtime: %v", err)
		}
	}()
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		t.Fatal(err)
	}
	if _, err := runtime.NewHostModuleBuilder("anki_timezone").NewFunctionBuilder().WithFunc(offsetForUTC).Export("offset_for_utc").NewFunctionBuilder().WithFunc(offsetsForLocal).Export("offsets_for_local").Instantiate(ctx); err != nil {
		t.Fatal(err)
	}
	compiled, err := runtime.CompileModule(ctx, guestFixture(t))
	if err != nil {
		t.Fatalf("compile WASI command: %v", err)
	}
	builder := runtime.NewHostModuleBuilder("anki_io")
	for _, definition := range compiled.ImportedFunctions() {
		module, name, _ := definition.Import()
		if module == "anki_io" {
			builder.NewFunctionBuilder().WithGoModuleFunction(api.GoModuleFunc(func(context.Context, api.Module, []uint64) { panic("CLI help attempted host stream I/O") }), definition.ParamTypes(), definition.ResultTypes()).Export(name)
		}
	}
	if _, err := builder.Instantiate(ctx); err != nil {
		t.Fatalf("instantiate command I/O guards: %v", err)
	}
	var output, diagnostics bytes.Buffer
	if _, err := runtime.InstantiateModule(ctx, compiled, wazero.NewModuleConfig().WithArgs("anki_grpc", "--help").WithStdout(&output).WithStderr(&diagnostics)); err != nil {
		t.Fatalf("execute WASI CLI: %v; %s", err, diagnostics.String())
	}
	if !bytes.Contains(output.Bytes(), []byte("Usage: anki_grpc <COMMAND>")) {
		t.Fatalf("WASI command did not run the Rust CLI: stdout=%q stderr=%q", output.String(), diagnostics.String())
	}
	if directory := os.Getenv("TEST_UNDECLARED_OUTPUTS_DIR"); directory != "" {
		if err := os.WriteFile(filepath.Join(directory, "rust-wasi-cli-help.txt"), output.Bytes(), 0o600); err != nil {
			t.Fatalf("retain direct CLI artifact: %v", err)
		}
	}
}
