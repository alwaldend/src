package app

import (
	"context"
	"errors"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/config"
)

// App exposes configuration operations to the CLI and retains collection placeholders.
// Anki collection processing belongs to the separate reconciliation change.
type App struct{}

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

// Dump prints parsed configurations through the pure text adapter, without Anki access.
func (application *App) Dump(ctx context.Context, input, format string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("dump configuration: %w", err)
	}
	data, err := config.Dump(input, format)
	if err != nil {
		return nil, fmt.Errorf("dump configuration: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("finish configuration dump: %w", err)
	}
	return data, nil
}
