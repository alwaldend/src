package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/anki_grpc/pkg/backend"
	"github.com/ankitects/anki/proto/anki/import_export"
	"github.com/bazelbuild/rules_go/go/runfiles"
	"google.golang.org/grpc"
)

// runCollection lends a prepared native connection to one operation and stops both backends.
// Database writes commit in place. Changed package operations use Anki's own
// atomic export; an empty output or unchanged operation skips package export.
func runCollection(ctx context.Context, input, output string, operation func(grpc.ClientConnInterface) (bool, error)) (err error) {
	input, err = expandPath(input)
	if err != nil {
		return fmt.Errorf("resolve collection input: %w", err)
	}
	if output != "" {
		output, err = expandPath(output)
		if err != nil {
			return fmt.Errorf("resolve collection output: %w", err)
		}
	}
	info, err := os.Lstat(input)
	if err != nil {
		return fmt.Errorf("inspect collection source: %w", err)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("collection source must be a regular file")
	}
	database := strings.EqualFold(filepath.Ext(input), ".anki2")
	if database && output != "" && output != input {
		return fmt.Errorf("build requires a collection package")
	}
	binary, err := runfiles.Rlocation("_main/projects/anki_grpc/cmd/anki_grpc/anki_grpc")
	if err != nil {
		binary, err = exec.LookPath("anki_grpc")
	}
	if err != nil {
		return fmt.Errorf("locate Anki server: %w", err)
	}
	worker := backend.New(backend.BackendSTDOptions{Binary: binary})
	collection := backend.NewCollection(backend.BackendCollectionOptions{Backend: worker, Path: input})
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if closeErr := collection.Stop(cleanup); closeErr != nil {
			err = errors.Join(err, fmt.Errorf("stop collection backend: %w", closeErr))
		}
		if stopErr := worker.Stop(cleanup); stopErr != nil {
			err = errors.Join(err, fmt.Errorf("stop Anki worker: %w", stopErr))
		}
	}()
	if err := worker.Start(ctx); err != nil {
		return fmt.Errorf("start Anki worker: %w", err)
	}
	if err := collection.Start(ctx); err != nil {
		return fmt.Errorf("prepare collection: %w", err)
	}
	connection, err := collection.Connection()
	if err != nil {
		return fmt.Errorf("get prepared collection: %w", err)
	}
	changed, err := operation(connection)
	if err != nil {
		return fmt.Errorf("run collection operation: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("finish collection operation: %w", err)
	}
	if database || !changed || output == "" {
		return nil
	}
	if _, err := import_export.NewBackendImportExportServiceClient(connection).ExportCollectionPackage(ctx, &import_export.ExportCollectionPackageRequest{
		OutPath: output, IncludeMedia: true,
	}); err != nil {
		return fmt.Errorf("export collection package: %w", err)
	}
	return nil
}
