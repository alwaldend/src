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

const (
	// TOML encodes configuration and state as schema-backed TOML.
	TOML = "toml"
	// JSON encodes configuration and state using ProtoJSON.
	JSON = "json"
)

// Encode serializes a protobuf message as TOML or JSON for config and state callers.
// An empty format selects TOML; JSON uses readable ProtoJSON with a final newline.
func Encode(message proto.Message, format string) ([]byte, error) {
	var data []byte
	var err error
	switch format {
	case "", TOML:
		data, err = encodeTOML(message)
	case JSON:
		data, err = encodeJSON(message)
	default:
		return nil, fmt.Errorf("unsupported encoding format %q; use toml or json", format)
	}
	if err != nil {
		return nil, fmt.Errorf("encode protobuf message: %w", err)
	}
	return data, nil
}

// Decode parses TOML or JSON into the supplied protobuf message, checking its schema.
// An empty format selects TOML. Resource and state semantics remain caller-owned.
func Decode(data []byte, format string, message proto.Message) error {
	if message == nil || !message.ProtoReflect().IsValid() {
		return fmt.Errorf("nil protobuf input message")
	}
	var err error
	switch format {
	case "", TOML:
		err = decodeTOML(data, message)
	case JSON:
		err = protojson.Unmarshal(data, message)
	default:
		return fmt.Errorf("unsupported encoding format %q; use toml or json", format)
	}
	if err != nil {
		return fmt.Errorf("decode protobuf message: %w", err)
	}
	return nil
}

// decodeTOML decodes TOML values into the supplied protobuf message.
func decodeTOML(b []byte, message proto.Message) error {
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

// encodeTOML converts protobuf JSON values directly to TOML.
// Numeric strings, field presence, and omitted defaults follow ProtoJSON.
func encodeTOML(message proto.Message) ([]byte, error) {
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

// encodeJSON serializes a protobuf message with readable ProtoJSON and a final newline.
func encodeJSON(message proto.Message) ([]byte, error) {
	data, err := (protojson.MarshalOptions{UseProtoNames: true, Indent: "  "}).Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("encode protobuf JSON: %w", err)
	}
	return append(data, '\n'), nil
}
