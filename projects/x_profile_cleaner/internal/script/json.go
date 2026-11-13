package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/archive"
	lua "github.com/yuin/gopher-lua"
)

const (
	// maxDepth bounds nested Lua result values independently of byte size.
	maxDepth = 64
	// maxValues bounds Lua result traversal of repeated references and tiny values.
	maxValues = 1000000
	// maxExactInteger is the largest integer exactly representable by the Lua number type.
	maxExactInteger = 9007199254740991
)

// nullValue marks JSON null without confusing it with an absent Lua table entry.
type nullValue struct{}

// installJSON exposes an immutable null sentinel for archive data values.
func installJSON(l *lua.LState) {
	null := l.NewUserData()
	null.Value = nullValue{}
	l.G.Registry.RawSetString("script_json_null", null)
	library := l.NewTable()
	library.RawSetString("null", null)
	l.SetGlobal("json", library)
}

// Decode projects JSON to fresh Lua values while retaining exact string identifiers.
// Numbers outside Lua's exact integer range are rejected rather than rounded.
func Decode(l *lua.LState, raw json.RawMessage) (lua.LValue, error) {
	if err := archive.ValidateJSON(raw); err != nil {
		return nil, fmt.Errorf("validate JSON projection: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var input any
	if err := decoder.Decode(&input); err != nil {
		return nil, fmt.Errorf("decode JSON projection: %w", err)
	}
	value, err := project(l, input)
	if err != nil {
		return nil, fmt.Errorf("project JSON value: %w", err)
	}
	return value, nil
}

// project converts already validated JSON without exposing shared archive objects.
func project(l *lua.LState, value any) (lua.LValue, error) {
	switch v := value.(type) {
	case nil:
		return l.G.Registry.RawGetString("script_json_null"), nil
	case bool:
		return lua.LBool(v), nil
	case string:
		return lua.LString(v), nil
	case json.Number:
		number, err := v.Float64()
		if err != nil {
			return nil, fmt.Errorf("decode JSON number: %w", err)
		}
		if math.IsNaN(number) || math.IsInf(number, 0) || (math.Trunc(number) == number && math.Abs(number) > maxExactInteger) {
			return nil, fmt.Errorf("JSON number exceeds Lua precision")
		}
		return lua.LNumber(number), nil
	case []any:
		result := l.NewTable()
		for _, item := range v {
			entry, err := project(l, item)
			if err != nil {
				return nil, fmt.Errorf("project JSON array: %w", err)
			}
			result.Append(entry)
		}
		return result, nil
	case map[string]any:
		result := l.NewTable()
		for key, item := range v {
			entry, err := project(l, item)
			if err != nil {
				return nil, fmt.Errorf("project JSON object: %w", err)
			}
			result.RawSetString(key, entry)
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported JSON value")
	}
}
