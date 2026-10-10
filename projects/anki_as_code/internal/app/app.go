package app

import (
	"context"
	"fmt"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/collection"
)

// App coordinates offline collection operations invoked by the command tree.
type App struct{}

// Export writes an editable text representation of a collection.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Export(ctx context.Context, input, output string) error {
	if err := collection.Export(ctx, input, output); err != nil {
		return fmt.Errorf("export collection: %w", err)
	}
	return nil
}

// Plan describes changes needed to reconcile a collection with its text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Plan(ctx context.Context, input, text string) (*collectionpb.CollectionChangePlan, error) {
	plan, err := collection.Plan(ctx, input, text)
	if err != nil {
		return nil, fmt.Errorf("plan collection: %w", err)
	}
	return plan, nil
}

// Apply reconciles an offline collection with its text representation.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Apply(ctx context.Context, input, text string) error {
	if err := collection.Apply(ctx, input, text); err != nil {
		return fmt.Errorf("apply collection: %w", err)
	}
	return nil
}

// Build creates an archive from a base collection and desired text.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) Build(ctx context.Context, base, text, output string) error {
	if err := collection.Build(ctx, base, text, output); err != nil {
		return fmt.Errorf("build collection: %w", err)
	}
	return nil
}

// GenerateIDsWithConfig assigns identities to notes using a collection config.
// The scaffold returns a not-implemented error without processing collections.
func (application *App) GenerateIDsWithConfig(ctx context.Context, paths []string, config string) (int, error) {
	count, err := collection.GenerateIDsWithConfig(ctx, paths, config)
	if err != nil {
		return 0, fmt.Errorf("generate collection identities: %w", err)
	}
	return count, nil
}
