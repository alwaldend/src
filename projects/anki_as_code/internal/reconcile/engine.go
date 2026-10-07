package reconcile

import (
	"context"
	"fmt"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/resource"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/state"
	"google.golang.org/protobuf/proto"
)

// Engine executes prepared plans using caller-owned native resources and state.
type Engine struct {
	// resources supplies generated Anki clients on the prepared connection.
	resources *resource.Store
	// states persists bindings in the namespace captured by the caller.
	states *state.Store
}

// New borrows stores for the same prepared collection and state namespace.
// The caller owns connection, process and collection lifetime.
func New(resources *resource.Store, states *state.Store) *Engine {
	return &Engine{resources: resources, states: states}
}

// Apply executes a supplied plan without loading configs or discovering resources.
// It rejects stale stored state before writes and leaves the supplied plan unchanged.
// Callers own snapshot freshness and exclusive offline access. Database writes may
// persist after a later failure; archive publication belongs to the collection owner.
func (engine *Engine) Apply(ctx context.Context, planned *collectionpb.CollectionPlan) (bool, error) {
	if planned == nil || planned.ExpectedState == nil || planned.ProposedBindings == nil {
		return false, fmt.Errorf("prepared collection plan is required")
	}
	stored, err := engine.states.Dump(ctx)
	if err != nil {
		return false, fmt.Errorf("read execution state: %w", err)
	}
	if !proto.Equal(stored, planned.ExpectedState) {
		return false, fmt.Errorf("stored bindings changed since planning; build a fresh plan")
	}
	run := nativeExecution{engine: engine, plan: proto.Clone(planned).(*collectionpb.CollectionPlan)}
	if err := run.execute(ctx); err != nil {
		return false, fmt.Errorf("execute collection plan: %w", err)
	}
	recorded := &collectionpb.State{Version: stored.Version, Resources: run.plan.ProposedBindings}
	stateChanged := !proto.Equal(stored, recorded)
	if stateChanged {
		if _, err := engine.states.Replace(ctx, recorded); err != nil {
			return false, fmt.Errorf("persist executed bindings: %w", err)
		}
	}
	return run.changed || stateChanged, nil
}
