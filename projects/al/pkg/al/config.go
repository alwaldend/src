package al

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"git.alwaldend.com/alwaldend/src/projects/al/api/al_proto"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/yaml"
)

// LoadConfigs merges JSON and YAML configuration in the supplied order.
func LoadConfigs(ctx context.Context, paths ...string) (*al_proto.Config, error) {
	return loadConfigs(ctx, nil, paths...)
}

// loadConfigs merges declarative configurations after explicit file substitution.
func loadConfigs(ctx context.Context, replacements *strings.Replacer, paths ...string) (*al_proto.Config, error) {
	res := &al_proto.Config{}
	for _, path := range paths {
		config, err := loadConfig(ctx, path, replacements)
		if err != nil {
			return nil, fmt.Errorf("could not load config %s: %w", path, err)
		}
		proto.Merge(res, config)
	}
	return res, nil
}

// loadConfig decodes one JSON or YAML configuration without code evaluation.
func loadConfig(ctx context.Context, path string, replacements *strings.Replacer) (*al_proto.Config, error) {
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
	if replacements != nil {
		configContent = []byte(replacements.Replace(string(configContent)))
	}
	switch extension {
	case ".yaml", ".json":
		configContentJson, err := yaml.YAMLToJSON(configContent)
		if err != nil {
			return nil, fmt.Errorf("could not parse config file %s", path)
		}
		err = protojson.Unmarshal(configContentJson, res)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal config file %s", path)
		}

	default:
		return nil, fmt.Errorf("invalid extension: %s", extension)
	}
	return res, nil
}

// DumpConfigs writes merged declarative configuration with private permissions.
func DumpConfigs(ctx context.Context, out string, paths ...string) error {
	return DumpConfigsWithSubstitutions(ctx, out, nil, paths...)
}

// DumpConfigsWithSubstitutions replaces explicit markers with JSON-escaped,
// whitespace-trimmed file contents before merging configuration. Callers must
// supply public build inputs only; secret values belong in runtime injection.
func DumpConfigsWithSubstitutions(ctx context.Context, out string, substitutions map[string]string, paths ...string) error {
	keys := make([]string, 0, len(substitutions))
	for key := range substitutions {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys)*2)
	for _, key := range keys {
		if key == "" {
			return fmt.Errorf("empty configuration substitution marker")
		}
		content, err := os.ReadFile(substitutions[key])
		if err != nil {
			return fmt.Errorf("read configuration substitution %q: %w", key, err)
		}
		encodedKey, err := json.Marshal(key)
		if err != nil {
			return fmt.Errorf("encode configuration substitution marker: %w", err)
		}
		encodedValue, err := json.Marshal(strings.TrimSpace(string(content)))
		if err != nil {
			return fmt.Errorf("encode configuration substitution value: %w", err)
		}
		pairs = append(pairs, string(encodedKey[1:len(encodedKey)-1]), string(encodedValue[1:len(encodedValue)-1]))
	}
	config, err := loadConfigs(ctx, strings.NewReplacer(pairs...), paths...)
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
		file, err = os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			return fmt.Errorf("could not open out file %s: %w", out, err)
		}
		defer file.Close()
		// Config dumps may include environment-derived secret values. Tighten
		// an existing file before writing and remove any previous trailing data.
		if err := file.Chmod(0o600); err != nil {
			return fmt.Errorf("could not restrict config output permissions: %w", err)
		}
		if err := file.Truncate(0); err != nil {
			return fmt.Errorf("could not truncate config output: %w", err)
		}
	}
	if _, err := file.Write(configMarshaled); err != nil {
		return fmt.Errorf("could not write config output: %w", err)
	}
	if _, err := file.WriteString("\n"); err != nil {
		return fmt.Errorf("could not finish config output: %w", err)
	}
	return nil
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
	case float64:
		return fp.Right(&al_proto.Json{Val: &al_proto.Json_ValFloat{ValFloat: valTyped}})
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
	if val == nil {
		return fp.Right[any](nil)
	}
	switch valTyped := val.Val.(type) {
	case nil:
		return fp.Right[any](nil)
	case *al_proto.Json_ValString:
		return fp.Right[any](valTyped.ValString)
	case *al_proto.Json_ValBool:
		return fp.Right[any](valTyped.ValBool)
	case *al_proto.Json_ValInt64:
		return fp.Right[any](valTyped.ValInt64)
	case *al_proto.Json_ValFloat:
		return fp.Right[any](valTyped.ValFloat)
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

// Parse a protobuf json into a proper protobuf message
func FromPbJsonToPb(val *al_proto.Json, target proto.Message) fp.EmptyEither {
	data, err := FromPbJson(val).Get()
	if err != nil {
		return fp.EmptyLeft(fmt.Errorf("could not convert pb json to data: %w", err))
	}
	if data == nil {
		return fp.EmptyRight()
	}
	dataJson, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fp.EmptyLeft(fmt.Errorf("could not convert data to json: %w", err))
	}
	if err := protojson.Unmarshal(dataJson, target); err != nil {
		return fp.EmptyLeft(fmt.Errorf("could not unmarshal protobuf configuration"))
	}
	return fp.EmptyRight()
}
