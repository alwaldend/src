package app

import (
	"context"
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
	"google.golang.org/grpc"
)

// App configures command workflows using the shared collection lifecycle.
// Options are set by the command layer; operations borrow a prepared connection.
type App struct {
	// Options selects the collection opened by WithCollection.
	Options Options
	// StateOptions selects the namespace shared by collection workflows and CLI state operations.
	StateOptions StateOptions
}

// Options selects the collection used by collection workflows.
type Options struct {
	// Collection names the .anki2 database or .colpkg archive opened by WithCollection.
	Collection string
}

// StateOptions selects the namespace shared by state and collection workflows.
// The command layer binds flags directly to these application options.
type StateOptions struct {
	// Key selects the metadata namespace; an empty key uses state.DefaultKey.
	Key string
}

// Export writes an editable text representation of a collection.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Export(ctx context.Context, input, output string) error {
	return errors.New("export is not implemented")
}

// Plan describes changes needed to reconcile a collection with its text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Plan(ctx context.Context, input, text string) (any, error) {
	return nil, errors.New("plan is not implemented")
}

// Apply reconciles an offline collection with its text representation.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Apply(ctx context.Context, input, text string) error {
	return errors.New("apply is not implemented")
}

// Build creates an archive from a base collection and desired text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Build(ctx context.Context, base, text, output string) error {
	return errors.New("build is not implemented")
}

// Dump loads and encodes config files without opening Anki.
func (application *App) Dump(ctx context.Context, input, format string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("dump configuration: %w", err)
	}
	data, err := config.New(input).Dump(format)
	if err != nil {
		return nil, fmt.Errorf("dump configuration: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("finish configuration dump: %w", err)
	}
	return data, nil
}

// WithCollection lends a prepared connection and owns native startup, persistence and cleanup.
// The callback's changed flag controls archive publication; errors discard restored archive changes.
func (application *App) WithCollection(ctx context.Context, operation func(grpc.ClientConnInterface) (bool, error)) error {
	input := application.Options.Collection
	if err := runCollection(ctx, input, input, operation); err != nil {
		return fmt.Errorf("run collection workflow: %w", err)
	}
	return nil
}
