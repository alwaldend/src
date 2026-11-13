package script

import (
	"context"
	"fmt"

	lua "github.com/yuin/gopher-lua"
	"google.golang.org/protobuf/proto"
)

// recordFunction projects one validated protobuf record through its descriptor.
func recordFunction[T proto.Message](state *lua.LState, ctx context.Context, name string, load func(context.Context) (T, error), failure *error) *lua.LFunction {
	return datasetFunction(state, name, func(l *lua.LState) (lua.LValue, error) {
		value, err := load(ctx)
		if err != nil {
			return nil, fmt.Errorf("read archive record: %w", err)
		}
		result, err := messageTable(l, value.ProtoReflect())
		if err != nil {
			return nil, fmt.Errorf("project archive record: %w", err)
		}
		return result, nil
	}, failure)
}

// recordsFunction projects validated records into a fresh dense Lua array.
func recordsFunction[T proto.Message](state *lua.LState, ctx context.Context, name string, load func(context.Context) ([]T, error), failure *error) *lua.LFunction {
	return datasetFunction(state, name, func(l *lua.LState) (lua.LValue, error) {
		values, err := load(ctx)
		if err != nil {
			return nil, fmt.Errorf("read archive records: %w", err)
		}
		result := l.NewTable()
		for _, value := range values {
			item, err := messageTable(l, value.ProtoReflect())
			if err != nil {
				return nil, fmt.Errorf("project archive records: %w", err)
			}
			result.Append(item)
		}
		return result, nil
	}, failure)
}

// datasetFunction owns argument checks and sticky failures for named data accessors.
func datasetFunction(state *lua.LState, name string, read func(*lua.LState) (lua.LValue, error), failure *error) *lua.LFunction {
	return state.NewFunction(func(l *lua.LState) int {
		if l.GetTop() != 0 {
			l.RaiseError("archive.%s expects no arguments", name)
			return 0
		}
		value, err := read(l)
		if err != nil {
			if *failure == nil {
				*failure = fmt.Errorf("read typed archive %s: %w", name, err)
			}
			l.RaiseError("typed archive dataset read or projection failed")
			return 0
		}
		l.Push(value)
		return 1
	})
}
