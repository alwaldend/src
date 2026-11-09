package script

import (
	"context"
	"fmt"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/archive"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/plan"
	lua "github.com/yuin/gopher-lua"
)

// NewState opens only computation libraries; its caller must close it.
func NewState(ctx context.Context) *lua.LState {
	state := lua.NewState(lua.Options{SkipOpenLibs: true})
	state.SetContext(ctx)
	for _, library := range []struct {
		// name selects one standard computation library.
		name string
		// open installs that library's native functions.
		open lua.LGFunction
	}{{lua.BaseLibName, lua.OpenBase}, {lua.TabLibName, lua.OpenTable}, {lua.StringLibName, lua.OpenString}, {lua.MathLibName, lua.OpenMath}} {
		state.Push(state.NewFunction(library.open))
		state.Push(lua.LString(library.name))
		state.Call(1, 0)
	}
	for _, name := range []string{"dofile", "loadfile", "load", "loadstring", "require", "module", "print", "_printregs", "collectgarbage"} {
		state.SetGlobal(name, lua.LNil)
	}
	return state
}

// Library installs trusted application helpers and validates their final result.
// A library instance belongs to one evaluation and must not retain state elsewhere.
type Library interface {
	// Install binds the library and its native plan contribution callback.
	Install(*lua.LState, func(*workerpb.Plan) error)
	// Validate checks helper failures and application-specific result provenance.
	Validate(*workerpb.Plan) error
}

// Options supplies trusted runtime configuration, independently of job payloads.
type Options struct {
	// Timeout shortens the parent deadline; zero inherits it and negatives are invalid.
	Timeout time.Duration
	// Library creates fresh helper state for each call; nil exposes generic APIs only.
	Library func(context.Context, *archive.Export, time.Time) Library
}

// Evaluator runs archive scripts with fresh interpreter and library state per call.
type Evaluator struct {
	// Options configures evaluation without changing its input or granting X access.
	Options Options
}

// Evaluate produces a complete protobuf plan without executing X operations.
func (e Evaluator) Evaluate(ctx context.Context, data *archive.Export, source string, now time.Time) (*workerpb.Plan, error) {
	if data == nil {
		return nil, fmt.Errorf("archive export is required")
	}
	if e.Options.Timeout < 0 {
		return nil, fmt.Errorf("evaluation timeout must not be negative")
	}
	if e.Options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.Options.Timeout)
		defer cancel()
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("before script evaluation: %w", err)
	}
	var library Library
	if e.Options.Library != nil {
		library = e.Options.Library(ctx, data, now)
		if library == nil {
			return nil, fmt.Errorf("library factory returned no library")
		}
	}
	result, err := evaluate(ctx, data, source, now, library)
	if stopped := ctx.Err(); stopped != nil {
		return nil, fmt.Errorf("script evaluation stopped: %w", stopped)
	}
	if library != nil {
		if failure := library.Validate(result); failure != nil {
			return nil, fmt.Errorf("validate script library: %w", failure)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("evaluate archive script: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("finish script evaluation: %w", err)
	}
	return result, nil
}

// evaluate owns the interpreter and shared archive/protobuf bindings for one call.
func evaluate(ctx context.Context, data *archive.Export, source string, now time.Time, library Library) (*workerpb.Plan, error) {
	state := NewState(ctx)
	defer state.Close()
	output := &collector{ctx: ctx, result: &workerpb.Plan{}}
	if err := output.install(state); err != nil {
		return nil, fmt.Errorf("install plan methods: %w", err)
	}
	installJSON(state)
	contextTable := state.NewTable()
	contextTable.RawSetString("account_id", lua.LString(data.AccountID))
	contextTable.RawSetString("evaluated_at", lua.LNumber(now.Unix()))
	state.SetGlobal("context", contextTable)
	var readFailure error
	archiveTable := state.NewTable()
	archiveTable.RawSetString("datasets", state.NewFunction(func(l *lua.LState) int {
		if l.GetTop() != 0 {
			l.RaiseError("archive.datasets expects no arguments")
			return 0
		}
		result := state.NewTable()
		for _, name := range data.Datasets() {
			result.Append(lua.LString(name))
		}
		l.Push(result)
		return 1
	}))
	archiveTable.RawSetString("read_raw", state.NewFunction(func(l *lua.LState) int {
		if l.GetTop() != 1 {
			l.RaiseError("archive.read_raw expects a dataset name")
			return 0
		}
		rows, err := data.Read(ctx, l.CheckString(1))
		if err != nil {
			l.RaiseError("archive dataset read failed")
			return 0
		}
		result := state.NewTable()
		for _, row := range rows {
			if err := ctx.Err(); err != nil {
				readFailure = fmt.Errorf("project archive dataset: %w", err)
				l.RaiseError("archive dataset read cancelled")
				return 0
			}
			value, err := Decode(state, row)
			if err != nil {
				readFailure = fmt.Errorf("project archive row: %w", err)
				l.RaiseError("archive row cannot be represented as JSON data")
				return 0
			}
			result.Append(value)
		}
		l.Push(result)
		return 1
	}))
	archiveTable.RawSetString("account", recordFunction(state, ctx, "account", data.Account, &readFailure))
	archiveTable.RawSetString("following", recordsFunction(state, ctx, "following", data.Following, &readFailure))
	archiveTable.RawSetString("likes", recordsFunction(state, ctx, "likes", data.Likes, &readFailure))
	archiveTable.RawSetString("tweets", recordsFunction(state, ctx, "tweets", data.Tweets, &readFailure))
	state.SetGlobal("archive", archiveTable)
	if library != nil {
		library.Install(state, output.append)
	}
	fn, err := state.Load(strings.NewReader(source), "worker-script")
	if err != nil {
		return nil, Failure("syntax", "worker-script", err)
	}
	if err := state.CallByParam(lua.P{Fn: fn, NRet: lua.MultRet, Protect: true}); err != nil {
		return nil, Failure("execution", "worker-script", err)
	}
	if err := data.Failure(); err != nil {
		return nil, Failure("archive", "worker-script", err)
	}
	if readFailure != nil {
		return nil, Failure("archive", "worker-script", readFailure)
	}
	if output.failure != nil {
		return nil, Failure("plan", "worker-script", output.failure)
	}
	if state.GetTop() != 0 {
		return nil, Failure("result", "worker-script", fmt.Errorf("scripts must collect operations without returning values"))
	}
	if err := plan.Validate(output.result); err != nil {
		return nil, Failure("plan", "worker-script", fmt.Errorf("validate collected plan: %w", err))
	}
	return output.result, nil
}
