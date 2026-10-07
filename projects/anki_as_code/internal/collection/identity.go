package collection

import (
	"context"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
)

// GenerateIDs assigns fresh note, card, and synchronization identities to a batch.
func GenerateIDs(ctx context.Context, paths []string) (int, error) {
	count, err := GenerateIDsWithConfig(ctx, paths, "")
	if err != nil {
		return 0, fmt.Errorf("generate note identities: %w", err)
	}
	return count, nil
}

// GenerateIDsWithConfig assigns fresh identities while reserving those referenced by the collection config.
func GenerateIDsWithConfig(ctx context.Context, paths []string, config string) (int, error) {
	count, err := textstore.GenerateIDsWithConfig(ctx, paths, config)
	if err != nil {
		return 0, fmt.Errorf("generate note identities from text: %w", err)
	}
	return count, nil
}
