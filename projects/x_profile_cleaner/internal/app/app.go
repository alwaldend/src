package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/oauth1"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/policy"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/xapi"
)

// App owns planning and execution dependencies. Configure it before use.
type App struct {
	// client performs account-scoped API operations.
	client *xapi.Client
	// policyTimeout bounds Lua evaluation, including its lazy reads and quota waits.
	policyTimeout time.Duration
	// logger reports operation progress through the command-owned handler.
	logger *slog.Logger
}

// Configure creates the API client and sets the positive policy evaluation limit.
// It performs no requests and leaves existing configuration intact on failure.
func (a *App) Configure(baseURL string, credentials oauth1.Credentials, policyTimeout time.Duration, logger *slog.Logger) error {
	if policyTimeout <= 0 {
		return fmt.Errorf("policy timeout must be positive")
	}
	if logger == nil {
		return fmt.Errorf("progress logger is required")
	}
	client, err := xapi.New(baseURL, credentials, logger)
	if err != nil {
		return fmt.Errorf("configure X client: %w", err)
	}
	a.client, a.policyTimeout, a.logger = client, policyTimeout, logger
	return nil
}

// PreparePlan requires exactly one source: a saved plan or a Lua configuration.
// Saved plans are validated without collecting inventory or evaluating policies.
func (a *App) PreparePlan(ctx context.Context, configPath, planPath string) (model.Plan, error) {
	if (configPath == "") == (planPath == "") {
		return model.Plan{}, fmt.Errorf("provide exactly one configuration path or saved plan path")
	}
	if planPath != "" {
		a.logger.InfoContext(ctx, "Loading saved plan", "path", planPath)
		plan, err := readPlan(planPath)
		if err != nil {
			return model.Plan{}, fmt.Errorf("load saved plan: %w", err)
		}
		a.logger.InfoContext(ctx, "Saved plan loaded", "account", plan.AccountID, "actions", len(plan.Actions))
		return plan, nil
	}
	a.logger.InfoContext(ctx, "Loading cleanup policy", "path", configPath)
	source, err := ReadPolicy(configPath)
	if err != nil {
		return model.Plan{}, fmt.Errorf("load cleanup policy: %w", err)
	}
	return a.PrepareSource(ctx, source, configPath)
}

// ReadPolicy loads a caller-selected policy file without evaluating it.
func ReadPolicy(path string) (string, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read configuration %q: %w", path, err)
	}
	return string(source), nil
}

// PrepareSource builds a read-only plan from editor or file text under the policy deadline.
func (a *App) PrepareSource(ctx context.Context, source, name string) (model.Plan, error) {
	inventory, err := a.client.Inventory(ctx)
	if err != nil {
		return model.Plan{}, fmt.Errorf("prepare cleanup inventory: %w", err)
	}
	policyContext, cancel := context.WithTimeout(ctx, a.policyTimeout)
	defer cancel()
	a.logger.InfoContext(ctx, "Evaluating cleanup policy", "timeout", a.policyTimeout)
	plan, err := policy.Build(policyContext, inventory.AccountID, inventory.Items, source, name, time.Now())
	if err != nil {
		return model.Plan{}, fmt.Errorf("prepare cleanup plan: %w", err)
	}
	a.logger.InfoContext(ctx, "Plan prepared", "account", plan.AccountID, "activities", len(plan.Decisions), "actions", len(plan.Actions))
	return plan, nil
}

// Apply executes only validated actions after checking the authenticated account.
// Independent endpoint workers retain per-endpoint and same-target plan order.
// The returned count includes all confirmed completions, not a prefix of the list.
// When supplied, confirmed receives each completed action's index exactly once,
// serially, including targets confirmed absent. Failed or uncertain requests omit it.
func (a *App) Apply(ctx context.Context, plan model.Plan, confirmed func(int)) (int, error) {
	if err := plan.Validate(); err != nil {
		return 0, fmt.Errorf("validate cleanup plan: %w", err)
	}
	account, err := a.client.Account(ctx)
	if err != nil {
		return 0, fmt.Errorf("identify execution account: %w", err)
	}
	if account != plan.AccountID {
		return 0, fmt.Errorf("authenticated account does not match plan account %s", plan.AccountID)
	}
	a.logger.InfoContext(ctx, "Applying cleanup plan", "account", account, "actions", len(plan.Actions))
	runContext, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	// step preserves list position and any earlier operation on the same target.
	type step struct {
		// index identifies this action in the supplied plan.
		index int
		// action is copied from the validated list without modification.
		action model.Action
		// previous waits for an earlier action on this target, across endpoints.
		previous <-chan struct{}
		// done releases dependent operations only after confirmed success.
		done chan struct{}
	}
	queues := map[model.ActionKind][]step{}
	previous := map[string]chan struct{}{}
	for index, action := range plan.Actions {
		done := make(chan struct{})
		queues[action.Kind] = append(queues[action.Kind], step{index: index, action: action, previous: previous[action.TargetID], done: done})
		previous[action.TargetID] = done
	}
	completed := 0
	var progress sync.Mutex
	var workers sync.WaitGroup
	for _, queue := range queues {
		workers.Add(1)
		go func(queue []step) {
			defer workers.Done()
			for _, next := range queue {
				if next.previous != nil {
					select {
					case <-next.previous:
					case <-runContext.Done():
						return
					}
				}
				if runContext.Err() != nil {
					return
				}
				a.logger.InfoContext(runContext, "Applying action", "action", next.index+1, "total", len(plan.Actions), "kind", next.action.Kind, "target", next.action.TargetID)
				if err := a.client.Execute(runContext, account, next.action); err != nil {
					cancel(fmt.Errorf("apply action %d of %d: %w", next.index+1, len(plan.Actions), err))
					return
				}
				progress.Lock()
				completed++
				if confirmed != nil {
					confirmed(next.index)
				}
				a.logger.InfoContext(runContext, "Action completed", "completed", completed, "total", len(plan.Actions), "action", next.index+1, "kind", next.action.Kind, "target", next.action.TargetID)
				progress.Unlock()
				close(next.done)
			}
		}(queue)
	}
	workers.Wait()
	if err := context.Cause(runContext); err != nil {
		a.logger.WarnContext(ctx, "Cleanup stopped; dispatched actions without confirmation may have completed", "completed", completed, "total", len(plan.Actions))
		return completed, fmt.Errorf("execute cleanup plan: %w", err)
	}
	a.logger.InfoContext(ctx, "Cleanup completed", "completed", completed, "total", len(plan.Actions))
	return completed, nil
}

// readPlan rejects unknown fields, trailing data, and malformed actions upfront.
func readPlan(path string) (model.Plan, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return model.Plan{}, fmt.Errorf("read plan %q: %w", path, err)
	}
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	var plan model.Plan
	if err := decoder.Decode(&plan); err != nil {
		return model.Plan{}, fmt.Errorf("decode plan %q: %w", path, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			return model.Plan{}, fmt.Errorf("decode trailing plan data: %w", err)
		}
		return model.Plan{}, fmt.Errorf("plan contains trailing JSON values")
	}
	if err := plan.Validate(); err != nil {
		return model.Plan{}, fmt.Errorf("validate plan %q: %w", path, err)
	}
	return plan, nil
}

// WritePlan atomically replaces a plan using a private file in the target directory.
func WritePlan(path string, plan model.Plan) error {
	if err := plan.Validate(); err != nil {
		return fmt.Errorf("validate plan before writing: %w", err)
	}
	content, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("encode plan: %w", err)
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".cleanup-plan-*")
	if err != nil {
		return fmt.Errorf("create plan output: %w", err)
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(append(content, '\n')); err != nil {
		closeErr := file.Close()
		if closeErr != nil {
			return fmt.Errorf("write plan: %w; close output: %w", err, closeErr)
		}
		return fmt.Errorf("write plan: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close plan output: %w", err)
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("publish plan %q: %w", path, err)
	}
	return nil
}
