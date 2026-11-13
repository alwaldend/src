package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"git.alwaldend.com/alwaldend/src/internal/x_sdk/pkg/oauth1"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/cleanup"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	operationplan "git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/plan"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/worker"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/xapi"
)

// Config supplies operator-owned worker references and mutation-only authorization.
type Config struct {
	// APIBaseURL identifies the X mutation endpoint, used only by Apply.
	APIBaseURL string
	// Credentials never reach the worker and are validated only for explicit apply.
	Credentials oauth1.Credentials
	// EvaluationTimeout bounds waiting for archive fetch, parsing and Lua evaluation.
	EvaluationTimeout time.Duration
	// Logger reports bounded progress at the application boundary.
	Logger *slog.Logger
	// WorkerSocket identifies the separately supervised source-owning worker.
	WorkerSocket string
	// ArchiveSource selects a provider configured by that worker's operator.
	ArchiveSource string
	// ArchiveKey is an opaque object key interpreted by the selected provider.
	ArchiveKey string
	// AccountID is the operator-selected expected archive and execution account.
	AccountID string
}

// App coordinates archive-backed planning and exact mutation execution.
// Configure it before use; the HTTP adapter serializes operations.
type App struct {
	// client is initialized only on explicit apply and retains quota state on resume.
	client *xapi.Client
	// config retains command-owned resources without granting worker X authorization.
	config Config
	// logger reports operation progress through the command-owned handler.
	logger *slog.Logger
}

// Configure records worker and mutation configuration without constructing an X client.
func (a *App) Configure(config Config) error {
	if config.EvaluationTimeout <= 0 || config.Logger == nil {
		return fmt.Errorf("positive evaluation timeout and progress logger are required")
	}
	if config.AccountID != "" && !model.ValidID(config.AccountID) {
		return fmt.Errorf("invalid expected account ID")
	}
	a.config, a.logger, a.client = config, config.Logger, nil
	return nil
}

// PreparePlan requires exactly one source: a saved plan or a Lua configuration.
// Saved plans are validated without collecting inventory or evaluating policies.
func (a *App) PreparePlan(ctx context.Context, configPath, planPath string) (operationplan.Document, error) {
	if (configPath == "") == (planPath == "") {
		return operationplan.Document{}, fmt.Errorf("provide exactly one configuration path or saved plan path")
	}
	if planPath != "" {
		a.logger.InfoContext(ctx, "Loading saved plan", "path", planPath)
		plan, err := readPlan(planPath)
		if err != nil {
			return operationplan.Document{}, fmt.Errorf("load saved plan: %w", err)
		}
		a.logger.InfoContext(ctx, "Saved plan loaded", "account", plan.AccountID, "actions", len(plan.Plan.Operations))
		return plan, nil
	}
	a.logger.InfoContext(ctx, "Loading cleanup policy", "path", configPath)
	source, err := ReadPolicy(configPath)
	if err != nil {
		return operationplan.Document{}, fmt.Errorf("load cleanup policy: %w", err)
	}
	document, err := a.PrepareSource(ctx, source, configPath)
	if err != nil {
		return operationplan.Document{}, fmt.Errorf("prepare cleanup source: %w", err)
	}
	return document, nil
}

// ReadPolicy loads a caller-selected policy file without evaluating it.
func ReadPolicy(path string) (string, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read configuration %q: %w", path, err)
	}
	return string(source), nil
}

// PrepareSource submits a source reference; archive data stays in the worker.
func (a *App) PrepareSource(ctx context.Context, source, name string) (operationplan.Document, error) {
	if a.config.WorkerSocket == "" {
		return operationplan.Document{}, fmt.Errorf("provide --worker-socket, --archive-source, --archive-key and --account for archive planning")
	}
	client, err := worker.Dial(a.config.WorkerSocket)
	if err != nil {
		return operationplan.Document{}, fmt.Errorf("connect archive worker: %w", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			a.logger.Warn("Worker connection cleanup failed", "error", err)
		}
	}()
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return operationplan.Document{}, fmt.Errorf("create archive operation ID: %w", err)
	}
	request := worker.Request{OperationID: hex.EncodeToString(id[:]), AccountID: a.config.AccountID, EvaluatedAt: time.Now().UTC(), Script: cleanup.Script(source), SourceName: a.config.ArchiveSource, ArchiveKey: a.config.ArchiveKey}
	operation, cancel := context.WithTimeout(ctx, a.config.EvaluationTimeout)
	defer cancel()
	a.logger.InfoContext(ctx, "Evaluating archive policy", "source", request.SourceName, "timeout", a.config.EvaluationTimeout)
	terminal := false
	defer func() {
		a.releaseJob(client, request, terminal)
	}()
	if _, err := client.Submit(operation, request); err != nil {
		return operationplan.Document{}, fmt.Errorf("submit archive policy %q: %w", name, err)
	}
	for {
		job, err := client.Get(operation, request)
		if err != nil {
			return operationplan.Document{}, fmt.Errorf("fetch archive result: %w", err)
		}
		if job.State != workerpb.EvaluationState_RUNNING {
			terminal = true
			if job.State != workerpb.EvaluationState_SUCCEEDED || job.Plan == nil {
				if job.Failure != nil {
					return operationplan.Document{}, fmt.Errorf("archive worker failed (%s): %s", job.Failure.Code, job.Failure.Message)
				}
				return operationplan.Document{}, fmt.Errorf("archive evaluation cancelled")
			}
			if err := cleanup.Validate(job.Plan); err != nil {
				return operationplan.Document{}, fmt.Errorf("validate cleanup result: %w", err)
			}
			plan := operationplan.Document{AccountID: request.AccountID, EvaluatedAt: request.EvaluatedAt, Plan: job.Plan}
			a.logger.InfoContext(ctx, "Plan prepared", "account", plan.AccountID, "activities", len(plan.Plan.GetCleanupReport().GetDecisions()), "actions", len(plan.Plan.Operations))
			return plan, nil
		}
		select {
		case <-operation.Done():
			return operationplan.Document{}, fmt.Errorf("wait for archive policy: %w", operation.Err())
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// releaseJob cancels abandoned work or releases a terminal result, without changing it.
func (a *App) releaseJob(client *worker.Client, request worker.Request, terminal bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if !terminal {
		if err := client.Cancel(ctx, request.AccountID, request.OperationID); err != nil {
			a.logger.Debug("Worker cancellation ended", "error", err)
			return
		}
		for !terminal {
			job, err := client.Get(ctx, request)
			if err != nil {
				a.logger.Debug("Worker cleanup ended", "error", err)
				return
			}
			terminal = job.State != workerpb.EvaluationState_RUNNING
			if !terminal {
				select {
				case <-ctx.Done():
					return
				case <-time.After(10 * time.Millisecond):
				}
			}
		}
	}
	if err := client.Delete(ctx, request.AccountID, request.OperationID); err != nil {
		a.logger.Warn("Worker result cleanup failed; retention will release it", "error", err)
	}
}

// Apply executes validated actions bound to the explicit expected archive account.
// The X service enforces write-credential ownership; no identity read is performed.
// Independent endpoint workers retain per-endpoint and same-target plan order.
// The returned count includes all confirmed completions, not a prefix of the list.
// When supplied, confirmed receives each completed action's index exactly once,
// serially, including targets confirmed absent. Failed or uncertain requests omit it.
func (a *App) Apply(ctx context.Context, plan operationplan.Document, confirmed func(int)) (int, error) {
	if err := plan.Validate(); err != nil {
		return 0, fmt.Errorf("validate cleanup plan: %w", err)
	}
	if err := xapi.ValidatePlan(plan.Plan); err != nil {
		return 0, fmt.Errorf("validate executable operations: %w", err)
	}
	account := a.config.AccountID
	if !model.ValidID(account) {
		return 0, fmt.Errorf("provide --account for explicit execution")
	}
	if account != plan.AccountID {
		return 0, fmt.Errorf("selected account does not match plan account %s", plan.AccountID)
	}
	if a.client == nil {
		client, err := xapi.New(a.config.APIBaseURL, a.config.Credentials, a.logger)
		if err != nil {
			return 0, fmt.Errorf("configure mutation client: %w", err)
		}
		a.client = client
	}
	a.logger.InfoContext(ctx, "Applying cleanup plan", "account", account, "actions", len(plan.Plan.Operations))
	runContext, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)
	// step preserves list position and any earlier operation on the same target.
	type step struct {
		// index identifies this action in the supplied plan.
		index int
		// action is copied from the validated list without modification.
		action *workerpb.Operation
		// previous waits for an earlier action on this target, across endpoints.
		previous <-chan struct{}
		// done releases dependent operations only after confirmed success.
		done chan struct{}
	}
	queues := map[string][]step{}
	previous := map[string]chan struct{}{}
	for index, action := range plan.Plan.Operations {
		done := make(chan struct{})
		queues[operationplan.OperationName(action)] = append(queues[operationplan.OperationName(action)], step{index: index, action: action, previous: previous[operationplan.TargetID(action)], done: done})
		previous[operationplan.TargetID(action)] = done
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
				a.logger.InfoContext(runContext, "Applying action", "action", next.index+1, "total", len(plan.Plan.Operations), "kind", operationplan.OperationName(next.action), "target", operationplan.TargetID(next.action))
				if err := a.client.Execute(runContext, account, next.action); err != nil {
					cancel(fmt.Errorf("apply action %d of %d: %w", next.index+1, len(plan.Plan.Operations), err))
					return
				}
				progress.Lock()
				completed++
				if confirmed != nil {
					confirmed(next.index)
				}
				a.logger.InfoContext(runContext, "Action completed", "completed", completed, "total", len(plan.Plan.Operations), "action", next.index+1, "kind", operationplan.OperationName(next.action), "target", operationplan.TargetID(next.action))
				progress.Unlock()
				close(next.done)
			}
		}(queue)
	}
	workers.Wait()
	if err := context.Cause(runContext); err != nil {
		a.logger.WarnContext(ctx, "Cleanup stopped; dispatched actions without confirmation may have completed", "completed", completed, "total", len(plan.Plan.Operations))
		return completed, fmt.Errorf("execute cleanup plan: %w", err)
	}
	a.logger.InfoContext(ctx, "Cleanup completed", "completed", completed, "total", len(plan.Plan.Operations))
	return completed, nil
}

// readPlan rejects unknown fields, trailing data, and malformed actions upfront.
func readPlan(path string) (operationplan.Document, error) {
	file, err := os.Open(path)
	if err != nil {
		return operationplan.Document{}, fmt.Errorf("open plan %q: %w", path, err)
	}
	content, readErr := io.ReadAll(io.LimitReader(file, operationplan.MaxDocumentBytes+1))
	if err := errors.Join(readErr, file.Close()); err != nil {
		return operationplan.Document{}, fmt.Errorf("read and close plan %q: %w", path, err)
	}
	if len(content) > operationplan.MaxDocumentBytes {
		return operationplan.Document{}, fmt.Errorf("plan document exceeds JSON byte bound")
	}
	var plan operationplan.Document
	if err := json.Unmarshal(content, &plan); err != nil {
		return operationplan.Document{}, fmt.Errorf("decode plan %q: %w", path, err)
	}
	return plan, nil
}

// WritePlan atomically replaces a plan using a private file in the target directory.
func WritePlan(path string, plan operationplan.Document) error {
	if err := plan.Validate(); err != nil {
		return fmt.Errorf("validate plan before writing: %w", err)
	}
	content, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return fmt.Errorf("encode plan: %w", err)
	}
	content = append(content, '\n')
	if len(content) > operationplan.MaxDocumentBytes {
		return fmt.Errorf("formatted plan document exceeds JSON byte bound")
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".cleanup-plan-*")
	if err != nil {
		return fmt.Errorf("create plan output: %w", err)
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(content); err != nil {
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
