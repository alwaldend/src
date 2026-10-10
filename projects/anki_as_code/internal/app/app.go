// Package app owns Anki collection operations invoked by the CLI.
package app

import (
	"context"
	"errors"
)

// App coordinates collection operations. This scaffold retains placeholders;
// collection processing is implemented in the separate business-logic change.
type App struct{}

// Export writes an editable text representation of a collection.
func (application *App) Export(ctx context.Context, input, output string) error {
	return errors.New("export is not implemented")
}

// Plan describes changes needed to reconcile a collection with its text.
func (application *App) Plan(ctx context.Context, input, text string) (any, error) {
	return nil, errors.New("plan is not implemented")
}

// Apply reconciles an offline collection with its text representation.
func (application *App) Apply(ctx context.Context, input, text string) error {
	return errors.New("apply is not implemented")
}

// Build creates an archive from a base collection and desired text.
func (application *App) Build(ctx context.Context, base, text, output string) error {
	return errors.New("build is not implemented")
}

// GenerateIDsWithConfig assigns identities to notes using a collection config.
func (application *App) GenerateIDsWithConfig(ctx context.Context, paths []string, config string) (int, error) {
	return 0, errors.New("generate-id is not implemented")
}
