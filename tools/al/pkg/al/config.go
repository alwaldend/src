package al

import (
	"context"
	"fmt"
	"os"

	"git.alwaldend.com/alwaldend/src/tools/al/api/al_proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"sigs.k8s.io/yaml"
)

func LoadConfigs(ctx context.Context, paths ...string) (*al_proto.Config, error) {
	res := &al_proto.Config{}
	for _, path := range paths {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context ended: %w", ctx.Err())
		default:
		}
		config := &al_proto.Config{}
		configContentYaml, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("could not open config file %s: %w", path, err)
		}
		configContentJson, err := yaml.YAMLToJSON(configContentYaml)
		err = protojson.Unmarshal(configContentJson, config)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal config file %s: %w", path, err)
		}
		proto.Merge(res, config)
	}
	return res, nil
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
		file, err = os.OpenFile(out, os.O_CREATE | os.O_WRONLY, 0o444)
		if err != nil {
			return fmt.Errorf("could not open out file %s: %w", out, err)
		}
		defer file.Close()
	}
	file.Write(configMarshaled)
	file.WriteString("\n")
	return nil
}
