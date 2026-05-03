package al

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/tools/al/api/al_proto"
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
    configFunc := newConfigFunc(res)
	if extension == ".yaml" || extension == ".json" {
		configContentJson, err := yaml.YAMLToJSON(configContent)
		err = protojson.Unmarshal(configContentJson, res)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal config file %s: %w", path, err)
		}
	} else if extension == ".lua" {
        state := lua.NewState()
        defer state.Close()
        state.SetGlobal("config", state.NewFunction(configFunc))
        err := state.DoString(string(configContent))
        if err != nil {
            return nil, fmt.Errorf("could not run lua: %w", err)
        }
	} else {
		return nil, fmt.Errorf("invalid extension: %s", extension)
	}
	return res, nil
}

func newConfigFunc(config *al_proto.Config) func(state *lua.LState) int {
    return func(state *lua.LState) int {
        parser := &LuaParser{}
        data, err := parser.Parse(state.ToTable(1))
        if err != nil {
            state.ArgError(1, fmt.Sprintf("could not parse the table: %s", err))
        }
        dataJson, err := json.Marshal(data)
        if err != nil {
            state.ArgError(1, fmt.Sprintf("could not marshal the table: %s", err))
        }
        dataProto := &al_proto.Config{}
        err = protojson.Unmarshal(dataJson, dataProto)
        if err != nil {
            state.ArgError(1, fmt.Sprintf("could not convert table to proto: %s", err))
        }
        proto.Merge(config, dataProto)
        return 0
    }
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

type LuaParser struct {
    enc map[*lua.LTable]bool
}

func (self *LuaParser) Parse(val lua.LValue) (any, error) {
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
		if self.enc[converted] {
            return nil, fmt.Errorf("nested table: %s", val)
		}
		self.enc[converted] = true
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
                valueParsed, err := self.Parse(value)
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
                valueParsed, err := self.Parse(value)
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
