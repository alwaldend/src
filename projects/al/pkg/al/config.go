package al

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"github.com/yuin/gopher-lua"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/yaml"
)

func LoadConfigs(ctx context.Context, paths ...string) (*al_proto.Config, error) {
	res := &al_proto.Config{}
	for _, path := range paths {
		config, err := loadConfig(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("could not load config %s: %w", path, err)
		}
		proto.Merge(res, config)
	}
	return res, nil
}

func loadConfig(ctx context.Context, path string) (*al_proto.Config, error) {
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("context ended: %w", ctx.Err())
	default:
	}
	extension := filepath.Ext(path)
	res := &al_proto.Config{}
	configContent, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config file %s: %w", path, err)
	}
	switch extension {
	case ".yaml", ".json":
		configContentJson, err := yaml.YAMLToJSON(configContent)
		err = protojson.Unmarshal(configContentJson, res)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal config file %s: %w", path, err)
		}
	case ".lua":
		state := lua.NewState()
		defer state.Close()
		state.SetGlobal("config", state.NewFunction(func(l *lua.LState) int {
			val := &al_proto.Config{}
			if err := parseTableArg(state, val); err != nil {
				state.ArgError(1, fmt.Sprintf("could not parse config: %s", err))
			} else {
				proto.Merge(res, val)
			}
			return 0
		}))
		state.SetGlobal("to_json", state.NewFunction(func(l *lua.LState) int {
			if res, err := toJson(l); err != nil {
				state.ArgError(1, fmt.Sprintf("could not convert to json: %s", err))
			} else {
				state.Push(lua.LString(res))
			}
			return 1
		}))
		state.SetGlobal("to_pb_json", state.NewFunction(func(l *lua.LState) int {
			if res, err := toPbJson(l); err != nil {
				state.ArgError(1, fmt.Sprintf("could not convert to protobuf json: %s", err))
			} else {
				state.Push(res)
			}
			return 1
		}))
		err := state.DoString(string(configContent))
		if err != nil {
			return nil, fmt.Errorf("could not run lua: %w", err)
		}
	default:
		return nil, fmt.Errorf("invalid extension: %s", extension)
	}
	return res, nil
}

func toPbJson(state *lua.LState) (lua.LValue, error) {
	data, err := parseLua(state.ToTable(1), nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse the table: %w", err)
	}
    data_pb, err := ToPbJson(data).Get()
    if err != nil {
        return nil, fmt.Errorf("could not covert data to protobuf json: %w", err)
    }
    data_pb_json, err := protojson.Marshal(data_pb)
    if err != nil {
        return nil, fmt.Errorf("could not marshal converted protobuf json: %w", err)
    }
    var data_parsed any
    if err := json.Unmarshal(data_pb_json, &data_parsed); err != nil {
        return nil, fmt.Errorf("could not parse converted protobuf as json: %w", err)
    }
    res, err := dumpLua(state, data_parsed)
    if err != nil {
        return nil, fmt.Errorf("could not dump json as lua: %w", err)
    }
    return res, nil
}

func toJson(state *lua.LState) ([]byte, error) {
	data, err := parseLua(state.ToTable(1), nil)
	if err != nil {
		return nil, fmt.Errorf("could not parse the table: %w", err)
	}
	res, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("could not marshal the data: %w", err)
	}
	return res, nil
}

func parseTableArg(state *lua.LState, target proto.Message) error {
	data, err := toJson(state)
	if err != nil {
		return fmt.Errorf("could not convert to json: %w", err)
	}
	err = protojson.Unmarshal(data, target)
	if err != nil {
		return fmt.Errorf("could not unmarshal as proto: %w", err)
	}
	return nil
}

func DumpConfigs(ctx context.Context, out string, paths ...string) error {
	config, err := LoadConfigs(ctx, paths...)
	if err != nil {
		return fmt.Errorf("could not load configs: %w", err)
	}
	configMarshaled, err := protojson.MarshalOptions{
		Indent: "    ",
	}.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal configs: %w", err)
	}
	var file *os.File
	if out == "" {
		file = os.Stdout
	} else {
		file, err = os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0o444)
		if err != nil {
			return fmt.Errorf("could not open out file %s: %w", out, err)
		}
		defer file.Close()
	}
	file.Write(configMarshaled)
	file.WriteString("\n")
	return nil
}

func dumpLua(state *lua.LState, val any) (lua.LValue, error) {
	switch converted := val.(type) {
	case bool:
		return lua.LBool(converted), nil
	case float64:
		return lua.LNumber(converted), nil
	case string:
		return lua.LString(converted), nil
	case json.Number:
		return lua.LString(converted), nil
	case []any:
		arr := state.CreateTable(len(converted), 0)
		for _, item := range converted {
            item_l, err := dumpLua(state, item)
            if err != nil {
                return nil, fmt.Errorf("could not convert item to lua: %w", err)
            }
			arr.Append(item_l)
		}
		return arr, nil
	case map[string]any:
		tbl := state.CreateTable(0, len(converted))
		for key, item := range converted {
            item_l, err := dumpLua(state, item)
            if err != nil {
                return nil, fmt.Errorf("could not covert item to lua: %w", err)
            }
			tbl.RawSetH(lua.LString(key), item_l)
		}
		return tbl, nil
	case nil:
		return lua.LNil, nil
    default:
        return nil, fmt.Errorf("unsupported type: %T: %v", val, val)
	}
}

func parseLua(val lua.LValue, enc map[*lua.LTable]bool) (any, error) {
	if enc == nil {
		enc = map[*lua.LTable]bool{}
	}
	switch converted := val.(type) {
	case lua.LBool:
		return bool(converted), nil
	case lua.LNumber:
		return float64(converted), nil
	case *lua.LNilType:
		return nil, nil
	case lua.LString:
		return string(converted), nil
	case *lua.LTable:
		if enc[converted] {
			return nil, fmt.Errorf("nested table: %s", val)
		}
		enc[converted] = true
		key, value := converted.Next(lua.LNil)
		switch key.Type() {
		case lua.LTNil: // empty table
			return []any{}, nil
		case lua.LTNumber:
			res := []any{}
			expectedKey := lua.LNumber(1)
			for key != lua.LNil {
				if key.Type() != lua.LTNumber {
					return nil, fmt.Errorf("invalid key: %s", key)
				}
				if expectedKey != key {
					return nil, fmt.Errorf("sparse array: expected %s, got %s", expectedKey, key)
				}
				valueParsed, err := parseLua(value, enc)
				if err != nil {
					return nil, fmt.Errorf("could not parse table value: %w", err)
				}
				res = append(res, valueParsed)
				expectedKey++
				key, value = converted.Next(key)
			}
			return res, nil
		case lua.LTString:
			res := make(map[string]any)
			for key != lua.LNil {
				if key.Type() != lua.LTString {
					return nil, fmt.Errorf("invalid key: %s", key)
				}
				valueParsed, err := parseLua(value, enc)
				if err != nil {
					return nil, fmt.Errorf("could not parse table value: %w", err)
				}
				res[key.String()] = valueParsed
				key, value = converted.Next(key)
			}
			return res, nil
		default:
			return nil, fmt.Errorf("invalid key: %s", key)
		}
	default:
		return nil, fmt.Errorf("invalid type: %s", val.Type())
	}
}

func VaultByName(config *al_proto.Config, name string) (*al_proto.VaultConn, error) {
	for i := range config.VaultConn {
		curVault := config.VaultConn[len(config.VaultConn)-1-i]
		if curVault.Name == name {
			return curVault, nil
		}
	}
	return nil, fmt.Errorf("missing vault with name %s", name)
}

func VaultAuthByName(config *al_proto.Config, name string) (*al_proto.VaultAuth, error) {
	for i := range config.VaultAuth {
		curAuth := config.VaultAuth[len(config.VaultAuth)-1-i]
		if curAuth.Name == name {
			return curAuth, nil
		}
	}
	return nil, fmt.Errorf("missing vault auth with name %s", name)
}

// Convert a value to a Protobuf json
func ToPbJson(val any) fp.Either[*al_proto.Json] {
	if val == nil {
		return fp.Right(&al_proto.Json{})
	}
	switch valTyped := val.(type) {
	case string:
		return fp.Right(&al_proto.Json{Val: &al_proto.Json_ValString{ValString: valTyped}})
	case bool:
		return fp.Right(&al_proto.Json{Val: &al_proto.Json_ValBool{ValBool: valTyped}})
	case int64:
		return fp.Right(&al_proto.Json{Val: &al_proto.Json_ValInt64{ValInt64: valTyped}})
	case []any:
		res := []*al_proto.Json{}
		for _, v := range valTyped {
			v2, err := ToPbJson(v).Get()
			if err != nil {
				return fp.Left[*al_proto.Json](err)
			}
			res = append(res, v2)
		}
		return fp.Right(&al_proto.Json{Val: &al_proto.Json_ValList{ValList: &al_proto.JsonList{Val: res}}})
	case map[string]any:
		res := map[string]*al_proto.Json{}
		for k, v := range valTyped {
			v2, err := ToPbJson(v).Get()
			if err != nil {
				return fp.Left[*al_proto.Json](err)
			}
			res[k] = v2
		}
		return fp.Right(&al_proto.Json{Val: &al_proto.Json_ValMap{ValMap: &al_proto.JsonMap{Val: res}}})
	default:
		return fp.Left[*al_proto.Json](fmt.Errorf("unsupported type %T: %s", val, val))
	}
}

// Parse Protobuf json into a normal type
func FromPbJson(val *al_proto.Json) fp.Result[any] {
	switch valTyped := val.Val.(type) {
	case *al_proto.Json_ValString, *al_proto.Json_ValBool, *al_proto.Json_ValInt64:
		return fp.Right[any](valTyped)
	case *al_proto.Json_ValList:
		res := []any{}
		for _, v := range valTyped.ValList.Val {
			v2, err := fp.Get(FromPbJson(v))
			if err != nil {
				return fp.Left[any](err)
			}
			res = append(res, v2)
		}
		return fp.Right[any](res)
	case *al_proto.Json_ValMap:
		res := map[string]any{}
		for k, v := range valTyped.ValMap.Val {
			v2, err := fp.Get(FromPbJson(v))
			if err != nil {
				return fp.Left[any](err)
			}
			res[k] = v2
		}
		return fp.Right[any](res)
	default:
		return fp.Left[any](fmt.Errorf("invalid item %T: %s", valTyped, val.Val))
	}
}
