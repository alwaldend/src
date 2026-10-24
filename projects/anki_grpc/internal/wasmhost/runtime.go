package wasmhost

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Run executes the complete embedded Rust CLI in the current supervised worker.
func Run(arguments []string, guest []byte) error {
	if err := RunCommand(context.Background(), arguments, guest, os.Stdin, os.Stdout, os.Stderr); err != nil {
		return fmt.Errorf("run embedded Rust Anki CLI: %w", err)
	}
	return nil
}

// RunCommand owns one WASI command and its host resources without replay or replacement.
// Serving takes ownership of the supplied protocol streams; callers also close
// them after execution to cover failures before the guest opens its transport.
func RunCommand(ctx context.Context, arguments []string, guest []byte, stdin io.ReadCloser, stdout io.WriteCloser, stderr io.Writer) (err error) {
	if len(arguments) == 0 {
		return errors.New("Anki CLI requires an argument vector including its program name")
	}
	for index, argument := range arguments {
		if strings.ContainsRune(argument, 0) {
			return fmt.Errorf("Anki CLI argument %d contains a null byte", index)
		}
	}
	configuration := wazero.NewRuntimeConfig().WithCloseOnContextDone(true)
	var cache wazero.CompilationCache
	var cacheError error
	if directory := os.Getenv("ANKI_WASM_CACHE_DIR"); directory != "" {
		if createErr := os.MkdirAll(directory, 0o700); createErr != nil {
			cacheError = fmt.Errorf("create Anki compilation cache %q: %w", directory, createErr)
		} else if cache, cacheError = wazero.NewCompilationCacheWithDir(directory); cacheError != nil {
			cacheError = fmt.Errorf("open Anki compilation cache %q: %w", directory, cacheError)
		} else {
			configuration = configuration.WithCompilationCache(cache)
		}
	}
	runtime := wazero.NewRuntimeWithConfig(ctx, configuration)
	defer func() {
		if closeErr := runtime.Close(context.Background()); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close Anki WASI runtime: %w", closeErr))
		}
		if cache != nil {
			if closeErr := cache.Close(context.Background()); closeErr != nil {
				err = errors.Join(err, fmt.Errorf("close Anki compilation cache: %w", closeErr))
			}
		}
	}()
	host := newHostIO(ctx, stdin, stdout, stderr, cacheError)
	defer func() {
		if closeErr := host.closeAll(); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("close Anki host streams: %w", closeErr))
		}
	}()
	if err := host.instantiate(runtime); err != nil {
		return fmt.Errorf("initialize Anki host I/O: %w", err)
	}
	if _, err := runtime.NewHostModuleBuilder("anki_timezone").NewFunctionBuilder().WithFunc(offsetForUTC).Export("offset_for_utc").NewFunctionBuilder().WithFunc(offsetsForLocal).Export("offsets_for_local").Instantiate(ctx); err != nil {
		return fmt.Errorf("initialize Anki timezone host: %w", err)
	}
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, runtime); err != nil {
		return fmt.Errorf("initialize WASI host: %w", err)
	}
	compiled, err := runtime.CompileModule(ctx, guest)
	if err != nil {
		return fmt.Errorf("compile Rust Anki command: %w", errors.Join(err, context.Cause(ctx)))
	}
	if compiled.ExportedFunctions()["_start"] == nil {
		return errors.New("Anki guest has no WASI command entry point")
	}
	temporaryDirectory, err := filepath.Abs(os.TempDir())
	if err != nil {
		return fmt.Errorf("resolve Anki temporary directory: %w", err)
	}
	config := wazero.NewModuleConfig().WithArgs(arguments...).WithStdin(stdin).WithStdout(stdout).WithStderr(stderr).
		WithFSConfig(wazero.NewFSConfig().WithDirMount("/", "/")).WithEnv("TMPDIR", temporaryDirectory).
		WithRandSource(rand.Reader).WithSysWalltime().WithSysNanotime().WithSysNanosleep()
	if _, err := runtime.InstantiateModule(ctx, compiled, config); err != nil {
		return fmt.Errorf("execute Rust Anki command: %w", errors.Join(err, context.Cause(ctx), host.failure))
	}
	if host.failure != nil {
		return fmt.Errorf("Rust Anki host I/O: %w", host.failure)
	}
	return nil
}
