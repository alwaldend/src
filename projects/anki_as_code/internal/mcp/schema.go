package mcp

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// schemaBuilder retains message definitions while traversing one protobuf graph.
// SDK inference describes encoding/json Go values, while these messages use
// ProtoJSON enum names, flattened oneofs and string-encoded 64-bit integers.
// Descriptors preserve that wire contract without maintaining per-RPC schemas.
type schemaBuilder struct {
	// definitions includes placeholders before visiting fields, permitting recursion.
	definitions map[string]any
}

// messageSchema describes the codec's canonical ProtoJSON object without Go struct inference.
func messageSchema(message protoreflect.MessageDescriptor) (map[string]any, error) {
	builder := schemaBuilder{definitions: map[string]any{}}
	reference, err := builder.message(message)
	if err != nil {
		return nil, fmt.Errorf("describe protobuf message %s: %w", message.FullName(), err)
	}
	return map[string]any{"type": "object", "allOf": []any{reference}, "$defs": builder.definitions}, nil
}

// message registers a reusable object and returns its local JSON Schema reference.
func (builder *schemaBuilder) message(message protoreflect.MessageDescriptor) (map[string]any, error) {
	name := string(message.FullName())
	reference := map[string]any{"$ref": "#/$defs/" + name}
	if _, exists := builder.definitions[name]; exists {
		return reference, nil
	}
	if message.ParentFile().Syntax() != protoreflect.Proto3 || strings.HasPrefix(name, "google.protobuf.") || message.ExtensionRanges().Len() != 0 {
		return nil, fmt.Errorf("unsupported protobuf message %s", name)
	}
	properties := map[string]any{}
	definition := map[string]any{"type": "object", "properties": properties, "additionalProperties": false}
	builder.definitions[name] = definition
	for index := 0; index < message.Fields().Len(); index++ {
		field := message.Fields().Get(index)
		value, err := builder.field(field)
		if err != nil {
			return nil, fmt.Errorf("describe field %s: %w", field.FullName(), err)
		}
		// ProtoJSON accepts null for an absent field, including optional and oneof fields.
		property := map[string]any{"anyOf": []any{value, map[string]any{"type": "null"}}}
		if comments := strings.TrimSpace(field.ParentFile().SourceLocations().ByDescriptor(field).LeadingComments); comments != "" {
			property["description"] = comments
		}
		properties[string(field.Name())] = property
	}
	var constraints []any
	for index := 0; index < message.Oneofs().Len(); index++ {
		oneof := message.Oneofs().Get(index)
		if oneof.IsSynthetic() {
			continue
		}
		// Exactly one branch matches an unset group or its single non-null member.
		unset := map[string]any{}
		alternatives := []any{map[string]any{"properties": unset}}
		for index := 0; index < oneof.Fields().Len(); index++ {
			fieldName := string(oneof.Fields().Get(index).Name())
			unset[fieldName] = map[string]any{"type": "null"}
			alternatives = append(alternatives, map[string]any{
				"required":   []string{fieldName},
				"properties": map[string]any{fieldName: map[string]any{"not": map[string]any{"type": "null"}}},
			})
		}
		constraints = append(constraints, map[string]any{"oneOf": alternatives})
	}
	if len(constraints) != 0 {
		definition["allOf"] = constraints
	}
	return reference, nil
}

// field describes cardinality around the underlying scalar or message value.
func (builder *schemaBuilder) field(field protoreflect.FieldDescriptor) (map[string]any, error) {
	if field.IsMap() {
		value, err := builder.scalar(field.MapValue())
		if err != nil {
			return nil, fmt.Errorf("describe map value: %w", err)
		}
		result := map[string]any{"type": "object", "additionalProperties": value}
		switch field.MapKey().Kind() {
		case protoreflect.StringKind:
		case protoreflect.BoolKind:
			result["propertyNames"] = map[string]any{"enum": []string{"true", "false"}}
		default:
			pattern := "^-?(0|[1-9][0-9]*)$"
			switch field.MapKey().Kind() {
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				pattern = "^(0|[1-9][0-9]*)$"
			}
			result["propertyNames"] = map[string]any{"pattern": pattern}
		}
		return result, nil
	}
	value, err := builder.scalar(field)
	if err != nil {
		return nil, fmt.Errorf("describe protobuf value: %w", err)
	}
	if field.IsList() {
		return map[string]any{"type": "array", "items": value}, nil
	}
	return value, nil
}

// scalar describes canonical ProtoJSON values and numeric enum compatibility.
func (builder *schemaBuilder) scalar(field protoreflect.FieldDescriptor) (map[string]any, error) {
	switch field.Kind() {
	case protoreflect.MessageKind:
		value, err := builder.message(field.Message())
		if err != nil {
			return nil, fmt.Errorf("describe nested message: %w", err)
		}
		return value, nil
	case protoreflect.StringKind:
		return map[string]any{"type": "string"}, nil
	case protoreflect.BytesKind:
		return map[string]any{"type": "string", "contentEncoding": "base64"}, nil
	case protoreflect.BoolKind:
		return map[string]any{"type": "boolean"}, nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return map[string]any{"type": "integer", "minimum": -2147483648, "maximum": 2147483647}, nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return map[string]any{"type": "integer", "minimum": 0, "maximum": uint64(4294967295)}, nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return map[string]any{"type": "string", "pattern": "^-?(0|[1-9][0-9]*)$", "description": "Signed 64-bit decimal integer encoded as a JSON string."}, nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return map[string]any{"type": "string", "pattern": "^(0|[1-9][0-9]*)$", "description": "Unsigned 64-bit decimal integer encoded as a JSON string."}, nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return map[string]any{"anyOf": []any{map[string]any{"type": "number"}, map[string]any{"enum": []string{"NaN", "Infinity", "-Infinity"}}}}, nil
	case protoreflect.EnumKind:
		var names []string
		for index := 0; index < field.Enum().Values().Len(); index++ {
			names = append(names, string(field.Enum().Values().Get(index).Name()))
		}
		return map[string]any{"anyOf": []any{map[string]any{"type": "string", "enum": names}, map[string]any{"type": "integer", "minimum": -2147483648, "maximum": 2147483647}}}, nil
	default:
		return nil, fmt.Errorf("unsupported protobuf kind %s", field.Kind())
	}
}

// protobufProperty attaches a descriptor schema's definitions to a tool envelope.
func protobufProperty(envelope map[string]any, name string, message protoreflect.MessageDescriptor) error {
	value, err := messageSchema(message)
	if err != nil {
		return fmt.Errorf("describe tool property %q: %w", name, err)
	}
	envelope["$defs"] = value["$defs"]
	envelope["properties"].(map[string]any)[name] = value["allOf"].([]any)[0]
	return nil
}
