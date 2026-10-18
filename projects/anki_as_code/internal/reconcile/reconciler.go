package reconcile

import (
	"context"
	"fmt"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"google.golang.org/protobuf/proto"
)

// Reconciler executes prepared plans using caller-owned native resources and state.
type Reconciler struct {
	// resources supplies generated Anki clients on the prepared connection.
	resources *resource.Store
	// states persists bindings in the namespace captured by the caller.
	states *state.Store
}

// New borrows stores for the same prepared collection and state namespace.
// The caller owns connection, process and collection lifetime.
func New(resources *resource.Store, states *state.Store) *Reconciler {
	return &Reconciler{resources: resources, states: states}
}

// Apply executes a supplied plan without loading configs or discovering resources.
// It rejects stale stored state before writes and leaves the supplied plan unchanged.
// Callers own snapshot freshness and exclusive offline access. Database writes may
// persist after a later failure; archive publication belongs to the collection owner.
func (reconciler *Reconciler) Apply(ctx context.Context, planned *ankiascodepb.CollectionPlan) (bool, error) {
	if planned == nil || planned.ExpectedState == nil || planned.ProposedBindings == nil {
		return false, fmt.Errorf("prepared collection plan is required")
	}
	stored, err := reconciler.states.Dump(ctx)
	if err != nil {
		return false, fmt.Errorf("read execution state: %w", err)
	}
	if !proto.Equal(stored, planned.ExpectedState) {
		return false, fmt.Errorf("stored bindings changed since planning; build a fresh plan")
	}
	run := applier{resources: reconciler.resources, plan: planned, bindings: proto.CloneOf(planned.ProposedBindings)}
	if err := run.execute(ctx); err != nil {
		return false, fmt.Errorf("execute collection plan: %w", err)
	}
	recorded := &ankiascodepb.State{Version: stored.Version, Resources: run.bindings}
	if stored.Resources == nil && proto.Equal(recorded.Resources, &ankiascodepb.Resources{}) {
		recorded.Resources = nil
	}
	stateChanged := !proto.Equal(stored, recorded)
	if stateChanged {
		if _, err := reconciler.states.Replace(ctx, recorded); err != nil {
			return false, fmt.Errorf("persist executed bindings: %w", err)
		}
	}
	return run.changed || stateChanged, nil
}
