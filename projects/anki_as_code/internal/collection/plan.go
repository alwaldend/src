package collection

import (
	"context"
	"errors"
	"fmt"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/planning"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/textstore"
)

// Plan returns the changes required to reconcile a collection with its declarations.
func Plan(ctx context.Context, input, dir string) (plan *collectionpb.CollectionChangePlan, err error) {
	input, err = textstore.ExpandPath(input)
	if err != nil {
		return &collectionpb.CollectionChangePlan{}, fmt.Errorf("resolve input path: %w", err)
	}
	dir, err = textstore.ExpandPath(dir)
	if err != nil {
		return &collectionpb.CollectionChangePlan{}, fmt.Errorf("resolve dir path: %w", err)
	}
	d, err := textstore.Load(dir)
	if err != nil {
		return plan, fmt.Errorf("load planned state: %w", err)
	}
	p, store, err := openCollection(ctx, input, "", false)
	if err != nil {
		return plan, fmt.Errorf("open planned collection: %w", err)
	}
	defer func() {
		err = errors.Join(err, store.Close(), p.Close())
	}()
	s, err := store.ValidateDesired(ctx, &d)
	if err != nil {
		return plan, fmt.Errorf("read current state: %w", err)
	}
	return planning.CalculatePlan(d, s), nil
}
