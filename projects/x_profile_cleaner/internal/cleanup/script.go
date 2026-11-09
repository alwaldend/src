package cleanup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/archive"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/script"
	lua "github.com/yuin/gopher-lua"
)

// Script explicitly invokes the cleanup library from a generic archive script.
// The closing delimiter cannot occur within the policy or overlap its suffix.
// This preserves policy text without escape conversion or an added newline.
func Script(source string) string {
	delimiter := "="
	for strings.Contains(source, "]"+delimiter+"]") || strings.HasSuffix(source, "]"+delimiter) {
		delimiter += "="
	}
	return "cleaner.plan([" + delimiter + "[\n" + source + "]" + delimiter + "])"
}

// library owns cleanup state for one generic script evaluation.
type library struct {
	// ctx bounds cleanup callbacks and archive reads.
	ctx context.Context
	// data is borrowed from the worker for this evaluation only.
	data *archive.Export
	// now fixes all age predicates.
	now time.Time
	// failure preserves helper errors even when Lua catches them.
	failure error
	// snapshot caches the validated cleanup projection for this call.
	snapshot *snapshot
}

// New creates an isolated cleanup library for script.Options.Library.
func New(ctx context.Context, data *archive.Export, now time.Time) script.Library {
	return &library{ctx: ctx, data: data, now: now}
}

// loadSnapshot loads cleanup records lazily and never reads X.
func (c *library) loadSnapshot() (*snapshot, error) {
	if c.snapshot != nil {
		return c.snapshot, nil
	}
	records, err := c.data.Tweets(c.ctx)
	if err != nil {
		return nil, fmt.Errorf("load cleanup snapshot: %w", err)
	}
	value, err := cleanupSnapshot(c.ctx, records)
	if err != nil {
		return nil, fmt.Errorf("build cleanup graph: %w", err)
	}
	c.snapshot = value
	return value, nil
}

// Install binds the cleanup DSL helper to a borrowed generic interpreter.
func (c *library) Install(state *lua.LState, appendPlan func(*workerpb.Plan) error) {
	bindings := state.NewTable()
	bindings.RawSetString("plan", state.NewFunction(func(l *lua.LState) int {
		source, valid := l.Get(1).(lua.LString)
		if l.GetTop() != 1 || !valid {
			c.failure = fmt.Errorf("cleanup helper requires one source string")
			l.RaiseError("cleaner.plan expects policy source")
			return 0
		}
		snapshot, err := c.loadSnapshot()
		if err != nil {
			c.failure = err
			l.RaiseError("cleanup archive validation failed")
			return 0
		}
		result, err := build(c.ctx, c.data.AccountID, snapshot.load, string(source), "cleanup-policy", c.now)
		if err != nil {
			c.failure = fmt.Errorf("evaluate cleanup policy: %w", err)
			l.RaiseError("cleanup policy failed")
			return 0
		}
		if err := Validate(result); err != nil {
			c.failure = fmt.Errorf("validate cleanup plan: %w", err)
			l.RaiseError("cleanup plan validation failed")
			return 0
		}
		if err := appendPlan(result); err != nil {
			c.failure = fmt.Errorf("append cleanup plan: %w", err)
			l.RaiseError("cleanup plan collection failed")
		}
		return 0
	}))
	state.SetGlobal("cleaner", bindings)
}

// Validate retains helper failures and verifies final explanations against the archive.
func (c *library) Validate(result *workerpb.Plan) error {
	if c.failure != nil {
		return fmt.Errorf("cleanup helper failed: %w", c.failure)
	}
	if result == nil || result.CleanupReport == nil {
		return nil
	}
	if err := Validate(result); err != nil {
		return fmt.Errorf("validate final cleanup report: %w", err)
	}
	inventory, err := c.loadSnapshot()
	if err != nil {
		return fmt.Errorf("verify final cleanup report: %w", err)
	}
	rows := make(map[string]int, len(inventory.records))
	for index, item := range inventory.records {
		rows[item.ID] = index
	}
	for _, decision := range result.CleanupReport.Decisions {
		index, ok := rows[decision.ItemId]
		if !ok {
			return fmt.Errorf("final cleanup report outside archive")
		}
		item := inventory.records[index]
		if string(item.Kind) != decision.Kind || item.SourceID != decision.SourceId || item.Text != decision.Text || !item.CreatedAt.Equal(decision.CreatedAt.AsTime()) {
			return fmt.Errorf("final cleanup report outside archive")
		}
	}
	return nil
}
