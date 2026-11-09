package script

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"unicode/utf8"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/api/workerpb"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/luaapi"
	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/plan"
	lua "github.com/yuin/gopher-lua"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// collector owns the complete native plan and first API failure for one evaluation.
type collector struct {
	// ctx bounds argument conversion and collection.
	ctx context.Context
	// result contains only fully validated appended operations and helper output.
	result *workerpb.Plan
	// bytes counts encoded protobuf contribution without rescanning earlier operations.
	bytes int
	// failure prevents publishing partial output after protected Lua calls fail.
	failure error
}

// fail retains the first API error, even when the interpreter catches it.
func (c *collector) fail(err error) error {
	if c.failure == nil {
		c.failure = err
	}
	return fmt.Errorf("plan collection failed: %w", c.failure)
}

// append validates and copies one complete contribution before changing the plan.
func (c *collector) append(value *workerpb.Plan) error {
	if c.failure != nil {
		return fmt.Errorf("plan collection already failed: %w", c.failure)
	}
	if err := c.ctx.Err(); err != nil {
		return c.fail(fmt.Errorf("append plan under deadline: %w", err))
	}
	if err := plan.Validate(value); err != nil {
		return c.fail(fmt.Errorf("validate plan contribution: %w", err))
	}
	size := proto.Size(value)
	if len(value.Operations) > plan.MaxOperations-len(c.result.Operations) || size > plan.MaxBytes-c.bytes {
		return c.fail(fmt.Errorf("plan operation or byte bound exceeded"))
	}
	if value.CleanupReport != nil && c.result.CleanupReport != nil {
		return c.fail(fmt.Errorf("cleanup report already supplied"))
	}
	copied := proto.Clone(value).(*workerpb.Plan)
	c.result.Operations = append(c.result.Operations, copied.Operations...)
	if copied.CleanupReport != nil {
		c.result.CleanupReport = copied.CleanupReport
	}
	c.bytes += size
	return nil
}

// install exposes a read-only API whose names and arguments come from protobuf.
func (c *collector) install(state *lua.LState) error {
	bindings := state.NewTable()
	fields := luaapi.OperationFields()
	for index := 0; index < fields.Len(); index++ {
		field := fields.Get(index)
		if field.Message() == nil {
			return fmt.Errorf("plan operation payload must be a message")
		}
		bindings.RawSetString(string(field.Name()), state.NewFunction(func(l *lua.LState) int {
			operation, err := c.operation(l, field)
			if err != nil {
				c.fail(fmt.Errorf("validate plan method arguments: %w", err))
				l.RaiseError("invalid plan.%s arguments", field.Name())
				return 0
			}
			if err := c.append(&workerpb.Plan{Operations: []*workerpb.Operation{operation}}); err != nil {
				l.RaiseError("plan.%s collection failed", field.Name())
			}
			return 0
		}))
	}
	api := state.NewUserData()
	metatable := state.NewTable()
	metatable.RawSetString("__index", state.NewFunction(func(l *lua.LState) int {
		name, ok := l.Get(2).(lua.LString)
		value := lua.LValue(lua.LNil)
		if ok {
			value = bindings.RawGetString(string(name))
		}
		if value == lua.LNil {
			c.fail(fmt.Errorf("unknown plan method"))
			l.RaiseError("unknown plan method")
			return 0
		}
		l.Push(value)
		return 1
	}))
	metatable.RawSetString("__newindex", state.NewFunction(func(l *lua.LState) int {
		c.fail(fmt.Errorf("plan methods cannot be replaced"))
		l.RaiseError("plan methods cannot be replaced")
		return 0
	}))
	metatable.RawSetString("__metatable", lua.LFalse)
	state.SetMetatable(api, metatable)
	state.SetGlobal("plan", api)
	return nil
}

// operation converts one method's exact argument list directly into its protobuf payload.
func (c *collector) operation(l *lua.LState, field protoreflect.FieldDescriptor) (*workerpb.Operation, error) {
	operation := &workerpb.Operation{}
	message := operation.ProtoReflect().Mutable(field).Message()
	converter := conversion{ctx: c.ctx, active: map[*lua.LTable]bool{}}
	argument := luaapi.OperationArgument(field)
	tableIndex := 1
	if argument != nil {
		maximum := 1
		if field.Message().Fields().Len() > 1 {
			maximum = 2
		}
		if l.GetTop() < 1 || l.GetTop() > maximum {
			return nil, fmt.Errorf("invalid positional argument count")
		}
		value, err := converter.value(l.Get(1), argument, message.NewField(argument), 0)
		if err != nil {
			return nil, fmt.Errorf("convert positional argument: %w", err)
		}
		message.Set(argument, value)
		if l.GetTop() == 1 || l.Get(2) == lua.LNil {
			return operation, nil
		}
		tableIndex = 2
	} else if l.GetTop() != 1 {
		return nil, fmt.Errorf("expected one payload table")
	}
	options, ok := l.Get(tableIndex).(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("expected an argument table")
	}
	if argument != nil && options.RawGetString(string(argument.Name())) != lua.LNil {
		return nil, fmt.Errorf("positional argument repeated in options")
	}
	if err := converter.fill(options, message, 0); err != nil {
		return nil, fmt.Errorf("convert argument table: %w", err)
	}
	return operation, nil
}

// conversion bounds direct Lua-to-protobuf traversal without JSON intermediates.
type conversion struct {
	// ctx interrupts result conversion under the job deadline.
	ctx context.Context
	// active detects cyclic tables before descending into a message or list.
	active map[*lua.LTable]bool
	// values bounds traversal independently of encoded protobuf size.
	values int
	// bytes counts retained string bytes before final exact protobuf size validation.
	bytes int
}

// enter admits one table for the current recursive path.
func (c *conversion) enter(table *lua.LTable, depth int) error {
	if err := c.ctx.Err(); err != nil {
		return fmt.Errorf("convert plan table: %w", err)
	}
	c.values++
	if depth > maxDepth || c.values > maxValues || c.active[table] {
		return fmt.Errorf("excessive or cyclic plan table")
	}
	c.active[table] = true
	return nil
}

// fill accepts only declared snake-case fields and one payload per protobuf oneof.
func (c *conversion) fill(table *lua.LTable, message protoreflect.Message, depth int) error {
	if err := c.enter(table, depth); err != nil {
		return fmt.Errorf("enter plan message: %w", err)
	}
	defer delete(c.active, table)
	selected := map[protoreflect.FullName]bool{}
	for key, value := table.Next(lua.LNil); key != lua.LNil; key, value = table.Next(key) {
		name, ok := key.(lua.LString)
		if !ok {
			return fmt.Errorf("plan fields must use string names")
		}
		field := message.Descriptor().Fields().ByName(protoreflect.Name(name))
		if field == nil {
			return fmt.Errorf("unknown plan field")
		}
		if group := field.ContainingOneof(); group != nil {
			if selected[group.FullName()] {
				return fmt.Errorf("operation has multiple payloads")
			}
			selected[group.FullName()] = true
		}
		if field.IsList() {
			entries, ok := value.(*lua.LTable)
			if !ok {
				return fmt.Errorf("repeated plan field must be a table")
			}
			if err := c.enter(entries, depth+1); err != nil {
				return fmt.Errorf("enter plan list: %w", err)
			}
			count := 0
			for key, _ := entries.Next(lua.LNil); key != lua.LNil; key, _ = entries.Next(key) {
				if err := c.ctx.Err(); err != nil {
					return fmt.Errorf("count plan list: %w", err)
				}
				count++
				index, ok := key.(lua.LNumber)
				if count > plan.MaxOperations || !ok || index < 1 || float64(index) != math.Trunc(float64(index)) || index > plan.MaxOperations {
					return fmt.Errorf("invalid plan list indices or size")
				}
			}
			list := message.Mutable(field).List()
			for index := 1; index <= count; index++ {
				item := entries.RawGetInt(index)
				if item == lua.LNil {
					return fmt.Errorf("sparse plan list")
				}
				converted, err := c.value(item, field, list.NewElement(), depth+2)
				if err != nil {
					return fmt.Errorf("convert plan list element: %w", err)
				}
				list.Append(converted)
			}
			delete(c.active, entries)
		} else {
			converted, err := c.value(value, field, message.NewField(field), depth+1)
			if err != nil {
				return fmt.Errorf("convert plan field: %w", err)
			}
			message.Set(field, converted)
		}
	}
	return nil
}

// value converts a field without Lua coercion or loss of integer precision.
func (c *conversion) value(value lua.LValue, field protoreflect.FieldDescriptor, destination protoreflect.Value, depth int) (protoreflect.Value, error) {
	if err := c.ctx.Err(); err != nil {
		return protoreflect.Value{}, fmt.Errorf("convert plan value: %w", err)
	}
	c.values++
	if c.values > maxValues {
		return protoreflect.Value{}, fmt.Errorf("plan value bound exceeded")
	}
	switch field.Kind() {
	case protoreflect.MessageKind:
		table, ok := value.(*lua.LTable)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("plan payload must be a table")
		}
		if err := c.fill(table, destination.Message(), depth); err != nil {
			return protoreflect.Value{}, fmt.Errorf("convert plan payload: %w", err)
		}
		return destination, nil
	case protoreflect.StringKind:
		text, ok := value.(lua.LString)
		if !ok || len(text) > plan.MaxBytes-c.bytes || !utf8.ValidString(string(text)) {
			return protoreflect.Value{}, fmt.Errorf("invalid plan string or size")
		}
		c.bytes += len(text)
		return protoreflect.ValueOfString(string(text)), nil
	case protoreflect.BoolKind:
		boolean, ok := value.(lua.LBool)
		if !ok {
			return protoreflect.Value{}, fmt.Errorf("plan flag must be boolean")
		}
		return protoreflect.ValueOfBool(bool(boolean)), nil
	case protoreflect.Int32Kind, protoreflect.Int64Kind:
		number, ok := value.(lua.LNumber)
		n := float64(number)
		if !ok || math.IsNaN(n) || math.IsInf(n, 0) || math.Trunc(n) != n || math.Abs(n) > maxExactInteger {
			return protoreflect.Value{}, fmt.Errorf("plan integer is not exact")
		}
		if field.Kind() == protoreflect.Int32Kind {
			if n < math.MinInt32 || n > math.MaxInt32 {
				return protoreflect.Value{}, fmt.Errorf("plan integer exceeds int32")
			}
			return protoreflect.ValueOfInt32(int32(n)), nil
		}
		return protoreflect.ValueOfInt64(int64(n)), nil
	default:
		return protoreflect.Value{}, fmt.Errorf("unsupported protobuf field kind")
	}
}

// messageTable follows protobuf descriptors, including empty lists and present scalars.
// Only archive records cross into Lua; operation plans stay native protobuf.
func messageTable(l *lua.LState, message protoreflect.Message) (*lua.LTable, error) {
	if err := l.Context().Err(); err != nil {
		return nil, fmt.Errorf("project protobuf message: %w", err)
	}
	table := l.NewTable()
	fields := message.Descriptor().Fields()
	for index := 0; index < fields.Len(); index++ {
		field := fields.Get(index)
		if field.HasPresence() && !message.Has(field) {
			continue
		}
		value := message.Get(field)
		if field.IsList() {
			result := l.NewTable()
			list := value.List()
			for i := 0; i < list.Len(); i++ {
				item, err := fieldValue(l, field, list.Get(i))
				if err != nil {
					return nil, fmt.Errorf("project protobuf list: %w", err)
				}
				result.Append(item)
			}
			table.RawSetString(string(field.Name()), result)
		} else {
			item, err := fieldValue(l, field, value)
			if err != nil {
				return nil, fmt.Errorf("project protobuf field: %w", err)
			}
			table.RawSetString(string(field.Name()), item)
		}
	}
	return table, nil
}

// fieldValue projects supported kinds, preserving exact integers beyond Lua's number range.
func fieldValue(l *lua.LState, field protoreflect.FieldDescriptor, value protoreflect.Value) (lua.LValue, error) {
	if err := l.Context().Err(); err != nil {
		return nil, fmt.Errorf("project protobuf value: %w", err)
	}
	switch field.Kind() {
	case protoreflect.MessageKind:
		result, err := messageTable(l, value.Message())
		if err != nil {
			return nil, fmt.Errorf("project nested protobuf message: %w", err)
		}
		return result, nil
	case protoreflect.StringKind:
		return lua.LString(value.String()), nil
	case protoreflect.BoolKind:
		return lua.LBool(value.Bool()), nil
	case protoreflect.EnumKind:
		item := field.Enum().Values().ByNumber(value.Enum())
		if item == nil {
			return nil, fmt.Errorf("unknown protobuf enum value")
		}
		return lua.LString(item.Name()), nil
	case protoreflect.Int32Kind, protoreflect.Int64Kind:
		number := value.Int()
		if number > maxExactInteger || number < -maxExactInteger {
			return lua.LString(strconv.FormatInt(number, 10)), nil
		}
		return lua.LNumber(number), nil
	default:
		return nil, fmt.Errorf("unsupported protobuf projection kind")
	}
}
