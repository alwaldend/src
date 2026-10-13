package codec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// DecodeTOML decodes TOML values into the supplied protobuf message.
func DecodeTOML(b []byte, message proto.Message) error {
	if message == nil || !message.ProtoReflect().IsValid() {
		return fmt.Errorf("nil protobuf input message")
	}
	var values map[string]any
	if _, err := toml.Decode(string(b), &values); err != nil {
		return fmt.Errorf("decode TOML: %w", err)
	}
	if err := checkJSONValue(values); err != nil {
		return fmt.Errorf("convert TOML values to JSON: %w", err)
	}
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("encode TOML values as JSON: %w", err)
	}
	if err := protojson.Unmarshal(data, message); err != nil {
		return fmt.Errorf("parse TOML with protobuf schema: %w", err)
	}
	return nil
}

// checkJSONValue rejects TOML datetimes that JSON would silently convert to strings.
// Protobuf remains responsible for validating the resulting JSON field types.
func checkJSONValue(value any) error {
	switch value := value.(type) {
	case time.Time:
		return fmt.Errorf("TOML dates and times require quoted text")
	case map[string]any:
		for name, entry := range value {
			if err := checkJSONValue(entry); err != nil {
				return fmt.Errorf("field %q: %w", name, err)
			}
		}
	case []any:
		for index, entry := range value {
			if err := checkJSONValue(entry); err != nil {
				return fmt.Errorf("array item %d: %w", index, err)
			}
		}
	case []map[string]any:
		for index, entry := range value {
			if err := checkJSONValue(entry); err != nil {
				return fmt.Errorf("table array item %d: %w", index, err)
			}
		}
	}
	return nil
}

// EncodeTOML converts protobuf JSON values directly to TOML.
// Numeric strings, field presence, and omitted defaults follow ProtoJSON.
func EncodeTOML(message proto.Message) ([]byte, error) {
	data, err := (protojson.MarshalOptions{UseProtoNames: true}).Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("encode protobuf JSON: %w", err)
	}
	var values map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&values); err != nil {
		return nil, fmt.Errorf("decode protobuf JSON values: %w", err)
	}
	var output bytes.Buffer
	encoder := toml.NewEncoder(&output)
	encoder.Indent = ""
	if err := encoder.Encode(values); err != nil {
		return nil, fmt.Errorf("encode TOML: %w", err)
	}
	return output.Bytes(), nil
}

// EncodeJSON serializes protobuf configuration with readable ProtoJSON and a final newline.
func EncodeJSON(message proto.Message) ([]byte, error) {
	data, err := (protojson.MarshalOptions{UseProtoNames: true, Indent: "  "}).Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("encode protobuf JSON: %w", err)
	}
	return append(data, '\n'), nil
}
