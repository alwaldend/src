package policy

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/planner"
	lua "github.com/yuin/gopher-lua"
)

// adapter owns the interpreter and translates its values to planner data.
type adapter struct {
	// state is confined to sequential configuration and callback execution.
	state *lua.LState
	// planner owns inventory, rule descriptions, graph traversal, and decisions.
	planner *planner.Planner
	// now fixes all Lua time helpers and activity ages.
	now time.Time
	// callbacks maps opaque planner identifiers to interpreter-owned functions.
	callbacks []*lua.LFunction
}

// Build registers Lua configuration as Go rules and evaluates a read-only plan.
// Interpreter values remain inside this adapter; the planner receives plain data.
func Build(ctx context.Context, account string, load func(context.Context, model.ItemKind, model.TimeWindow) ([]model.Item, error), source, filename string, now time.Time) (model.Plan, error) {
	state := lua.NewState(lua.Options{SkipOpenLibs: true})
	defer state.Close()
	state.SetContext(ctx)
	for _, library := range []struct {
		// name identifies the restricted standard library.
		name string
		// open initializes only that library.
		open lua.LGFunction
	}{{lua.BaseLibName, lua.OpenBase}, {lua.TabLibName, lua.OpenTable}, {lua.StringLibName, lua.OpenString}, {lua.MathLibName, lua.OpenMath}} {
		state.Push(state.NewFunction(library.open))
		state.Push(lua.LString(library.name))
		state.Call(1, 0)
	}
	for _, name := range []string{"dofile", "loadfile", "load", "loadstring", "require", "module", "print", "_printregs", "collectgarbage"} {
		state.SetGlobal(name, lua.LNil)
	}
	e := &adapter{state: state, planner: planner.New(ctx, account, load, now), now: now}
	e.bind()
	fn, err := state.LoadString(source)
	if err != nil {
		return model.Plan{}, fmt.Errorf("parse policy %q: %w", filename, err)
	}
	if err := state.CallByParam(lua.P{Fn: fn, NRet: 0, Protect: true}); err != nil {
		return model.Plan{}, fmt.Errorf("execute policy %q: %w", filename, err)
	}
	plan, err := e.planner.Evaluate(e.evaluate)
	if err != nil {
		return model.Plan{}, fmt.Errorf("evaluate policy %q: %w", filename, err)
	}
	return plan, nil
}

// bind exposes strict configuration functions with no interpreter state in the planner.
func (e *adapter) bind() {
	e.state.SetGlobal("set_time_window", e.state.NewFunction(func(l *lua.LState) int {
		checkArity(l, 0, 2)
		if err := e.planner.SetTimeWindow(model.TimeWindow{Start: checkTime(l, 1), End: checkTime(l, 2)}); err != nil {
			l.RaiseError("set_time_window: %s", err)
		}
		return 0
	}))
	e.state.SetGlobal("days_ago", e.state.NewFunction(e.daysAgo))
	for _, name := range planner.Operations() {
		for _, predicate := range []bool{false, true} {
			binding := name
			if predicate {
				binding += "_when"
			}
			e.state.SetGlobal(binding, e.state.NewFunction(func(l *lua.LState) int {
				checkArity(l, 1, 2)
				rule := e.readOptions(l, 2)
				rule.Operation = name
				if predicate {
					rule.Predicate = e.remember(l.CheckFunction(1))
				} else {
					rule.IDs = checkIDs(l, 1)
				}
				if err := e.planner.Add(rule); err != nil {
					l.RaiseError("%s: %s", binding, err)
				}
				return 0
			}))
		}
	}
	for _, kind := range []model.ItemKind{model.ItemPost, model.ItemRepost} {
		e.state.SetGlobal("get_"+string(kind), e.state.NewFunction(func(l *lua.LState) int {
			checkArity(l, 1, 1)
			id := checkID(l, 1)
			item, ok, err := e.planner.Get(kind, id)
			if err != nil {
				l.RaiseError("get %s: %s", kind, err)
			}
			if ok {
				l.Push(e.record(item))
			} else {
				l.Push(lua.LNil)
			}
			return 1
		}))
		e.state.SetGlobal("get_"+string(kind)+"s", e.state.NewFunction(func(l *lua.LState) int {
			checkArity(l, 0, 0)
			items, err := e.planner.List(kind)
			if err != nil {
				l.RaiseError("list %s: %s", kind, err)
			}
			l.Push(e.records(items))
			return 1
		}))
	}
	for name, kind := range map[string]model.ReferenceKind{"get_replies": model.ReferenceReply, "get_quotes": model.ReferenceQuote} {
		e.state.SetGlobal(name, e.state.NewFunction(func(l *lua.LState) int {
			checkArity(l, 1, 1)
			id := checkID(l, 1)
			items, err := e.planner.Children(kind, id)
			if err != nil {
				l.RaiseError("%s: %s", name, err)
			}
			l.Push(e.records(items))
			return 1
		}))
	}
}

// checkArity rejects arguments that would otherwise be silently ignored by Lua.
func checkArity(l *lua.LState, minimum, maximum int) {
	if count := l.GetTop(); count < minimum || count > maximum {
		l.RaiseError("expected %d..%d arguments, received %d", minimum, maximum, count)
	}
}

// daysAgo returns a whole-second UTC timestamp without loading inventory.
func (e *adapter) daysAgo(l *lua.LState) int {
	checkArity(l, 1, 1)
	days, ok := l.Get(1).(lua.LNumber)
	minimum := time.Date(2010, time.November, 6, 0, 0, 0, 0, time.UTC)
	if !ok || math.IsNaN(float64(days)) || math.IsInf(float64(days), 0) || days < 0 || float64(days) > e.now.Sub(minimum).Hours()/24 {
		l.ArgError(1, "expected finite nonnegative days within the supported timestamp range")
		return 0
	}
	stamp := e.now.Add(-time.Duration(float64(days) * float64(24*time.Hour))).UTC().Truncate(time.Second)
	l.Push(lua.LString(stamp.Format(time.RFC3339)))
	return 1
}

// remember retains executable Lua code behind an opaque planner callback identifier.
func (e *adapter) remember(fn *lua.LFunction) planner.CallbackID {
	e.callbacks = append(e.callbacks, fn)
	return planner.CallbackID(len(e.callbacks))
}

// readOptions parses a name and callback references without exposing Lua to the planner.
func (e *adapter) readOptions(l *lua.LState, index int) planner.Rule {
	result := planner.Rule{}
	value := l.Get(index)
	if value == lua.LNil {
		return result
	}
	if fn, ok := value.(*lua.LFunction); ok {
		result.Include = e.remember(fn)
		return result
	}
	table, ok := value.(*lua.LTable)
	if !ok {
		l.ArgError(index, "expected an options table or preservation inclusion function")
		return result
	}
	table.ForEach(func(key, value lua.LValue) {
		name, ok := key.(lua.LString)
		if !ok {
			l.ArgError(index, "option names must be strings")
			return
		}
		if name == "name" {
			label, ok := value.(lua.LString)
			if !ok || strings.TrimSpace(string(label)) == "" {
				l.ArgError(index, "rule name must be a nonempty string")
				return
			}
			result.Name = string(label)
			return
		}
		if name != "include_post" && name != "follow_edge" {
			l.ArgError(index, "unknown option "+string(name))
			return
		}
		fn, ok := value.(*lua.LFunction)
		if !ok {
			l.ArgError(index, "callback options must contain functions")
			return
		}
		if name == "include_post" {
			result.Include = e.remember(fn)
		} else {
			result.Follow = e.remember(fn)
		}
	})
	return result
}

// evaluate invokes one interpreter-owned function with fresh plain-data projections.
func (e *adapter) evaluate(id planner.CallbackID, input planner.CallbackInput) (bool, error) {
	if id < 1 || int(id) > len(e.callbacks) {
		return false, fmt.Errorf("unknown callback identifier %d", id)
	}
	var value lua.LValue = lua.LNil
	if input.Post != nil {
		value = e.record(*input.Post)
	}
	if input.Edge != nil {
		edge := e.state.NewTable()
		for key, value := range map[string]string{"from_id": input.Edge.FromID, "to_id": input.Edge.ToID, "kind": string(input.Edge.Kind), "direction": input.Edge.Direction} {
			edge.RawSetString(key, lua.LString(value))
		}
		value = edge
	}
	args := []lua.LValue{value}
	if input.SeedID != "" {
		ctx := e.state.NewTable()
		ctx.RawSetString("seed_id", lua.LString(input.SeedID))
		args = append(args, ctx)
	}
	if err := e.state.CallByParam(lua.P{Fn: e.callbacks[id-1], NRet: 1, Protect: true}, args...); err != nil {
		return false, fmt.Errorf("call selection function: %w", err)
	}
	result := e.state.Get(-1)
	e.state.Pop(1)
	flag, ok := result.(lua.LBool)
	if !ok {
		return false, fmt.Errorf("selection callback returned %s; expected boolean", result.Type())
	}
	return bool(flag), nil
}

// records creates fresh Lua tables for one query result.
func (e *adapter) records(items []model.Item) *lua.LTable {
	result := e.state.NewTable()
	for _, item := range items {
		result.Append(e.record(item))
	}
	return result
}

// checkTime accepts an omitted bound or a strict RFC3339 timestamp string.
func checkTime(l *lua.LState, index int) time.Time {
	if l.Get(index) == lua.LNil {
		return time.Time{}
	}
	value, ok := l.Get(index).(lua.LString)
	if !ok {
		l.ArgError(index, "expected an RFC3339 timestamp string or nil")
		return time.Time{}
	}
	stamp, err := time.Parse(time.RFC3339, string(value))
	if err != nil || stamp.IsZero() {
		l.ArgError(index, "expected a supported RFC3339 timestamp string or nil")
		return time.Time{}
	}
	return stamp.UTC()
}

// checkID refuses numeric coercion, preserving the precision of X identifiers.
func checkID(l *lua.LState, index int) string {
	value, ok := l.Get(index).(lua.LString)
	if !ok || !model.ValidID(string(value)) {
		l.ArgError(index, "expected a canonical decimal identifier string")
		return ""
	}
	return string(value)
}

// checkIDs accepts one ID or a dense array, preserving the array's order.
func checkIDs(l *lua.LState, index int) []string {
	table, ok := l.Get(index).(*lua.LTable)
	if !ok {
		return []string{checkID(l, index)}
	}
	count := 0
	table.ForEach(func(lua.LValue, lua.LValue) { count++ })
	if count != table.Len() {
		l.ArgError(index, "expected a dense array of identifier strings")
	}
	ids := make([]string, count)
	for i := range ids {
		value, ok := table.RawGetInt(i + 1).(lua.LString)
		if !ok || !model.ValidID(string(value)) {
			l.ArgError(index, "expected a dense array of canonical decimal identifier strings")
		}
		ids[i] = string(value)
	}
	return ids
}

// record constructs fresh Lua data so callback mutations cannot alter inventory.
func (e *adapter) record(item model.Item) *lua.LTable {
	l := e.state
	result := l.NewTable()
	for key, value := range map[string]string{"id": item.ID, "author_id": item.AuthorID, "kind": string(item.Kind), "text": item.Text} {
		result.RawSetString(key, lua.LString(value))
	}
	if item.SourceID != "" {
		result.RawSetString("source_id", lua.LString(item.SourceID))
	}
	result.RawSetString("created_at", lua.LNumber(item.CreatedAt.Unix()))
	result.RawSetString("age_days", lua.LNumber(e.now.Sub(item.CreatedAt).Hours()/24))
	metrics := l.NewTable()
	for key, value := range item.Metrics {
		metrics.RawSetString(key, lua.LNumber(value))
	}
	result.RawSetString("metrics", metrics)
	refs := l.NewTable()
	isReply, isQuote := false, false
	for _, ref := range item.References {
		isReply = isReply || ref.Kind == model.ReferenceReply
		isQuote = isQuote || ref.Kind == model.ReferenceQuote
		table := l.NewTable()
		table.RawSetString("id", lua.LString(ref.ID))
		table.RawSetString("kind", lua.LString(ref.Kind))
		refs.Append(table)
	}

	if item.Language != nil {
		result.RawSetString("language", lua.LString(*item.Language))
	}
	if item.Attachments != nil {
		attachments := l.NewTable()
		for key, values := range map[string][]string{"media_keys": item.Attachments.MediaKeys, "poll_ids": item.Attachments.PollIDs, "media_source_ids": item.Attachments.MediaSourceIDs} {
			if values == nil {
				continue
			}
			array := l.NewTable()
			for _, value := range values {
				array.Append(lua.LString(value))
			}
			attachments.RawSetString(key, array)
		}
		result.RawSetString("attachments", attachments)
	}
	if item.Entities != nil {
		entities := l.NewTable()
		for key, values := range map[string][]model.TextEntity{"hashtags": item.Entities.Hashtags, "cashtags": item.Entities.Cashtags, "mentions": item.Entities.Mentions, "urls": item.Entities.URLs} {
			if values == nil {
				continue
			}
			array := l.NewTable()
			for _, value := range values {
				entity := l.NewTable()
				entity.RawSetString("start", lua.LNumber(value.Start))
				entity.RawSetString("end", lua.LNumber(value.End))
				for key, value := range map[string]*string{"tag": value.Tag, "id": value.ID, "username": value.Username, "url": value.URL, "expanded_url": value.ExpandedURL, "display_url": value.DisplayURL} {
					if value != nil {
						entity.RawSetString(key, lua.LString(*value))
					}
				}
				array.Append(entity)
			}
			entities.RawSetString(key, array)
		}
		result.RawSetString("entities", entities)
	}
	result.RawSetString("references", refs)
	result.RawSetString("is_reply", lua.LBool(isReply))
	result.RawSetString("is_quote", lua.LBool(isQuote))
	return result
}
