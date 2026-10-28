package graph

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// buildDepsCaptureFlags describes the complete caller-owned capture contract.
const buildDepsCaptureFlags = "--output=streamed_jsonproto --consistent_labels --noimplicit_deps --notool_deps --nonodep_deps --proto:output_rule_attrs=name,deps,srcs --proto:default_values=true --proto:flatten_selects=false --proto:rule_inputs_and_outputs=true"

// dependencyAttribute records normalized relevant metadata for conflict checks.
type dependencyAttribute struct {
	// Name is the literal deps or srcs attribute name.
	Name string `json:"name"`
	// Type is Bazel's declared attribute encoding.
	Type string `json:"type"`
	// ExplicitlySpecified distinguishes explicit and default values.
	ExplicitlySpecified *bool `json:"explicitlySpecified"`
	// Nodep marks label-valued metadata rather than dependency edges.
	Nodep bool `json:"nodep"`
	// SourceAspectName identifies aspect-owned attribute values to exclude.
	SourceAspectName string `json:"sourceAspectName"`
	// Labels contains normalized values, never selector condition labels.
	Labels []string `json:"-"`
}

// buildInputs validates capture metadata and intersects literal deps/srcs with inputs.
// The normalized metadata signature preserves conflicting-record detection even
// when the conflicting attributes would otherwise be excluded by the filter.
func buildInputs(rule *ruleRecord) ([]string, string, error) {
	if len(rule.Attribute) == 0 {
		return nil, "", fmt.Errorf("build-deps-only requires a name attribute marker on every RULE; recapture with %s", buildDepsCaptureFlags)
	}
	var attributes []json.RawMessage
	if err := json.Unmarshal(rule.Attribute, &attributes); err != nil {
		return nil, "", fmt.Errorf("decode rule attributes; recapture with %s: %w", buildDepsCaptureFlags, err)
	}
	marker := false
	metadata := map[string]string{}
	selected := map[string]bool{}
	for _, raw := range attributes {
		var identity namedRecord
		if err := json.Unmarshal(raw, &identity); err != nil {
			return nil, "", fmt.Errorf("decode attribute name: %w", err)
		}
		if identity.Name == "name" {
			marker = true
			continue
		}
		if identity.Name != "deps" && identity.Name != "srcs" {
			continue
		}
		var attribute dependencyAttribute
		if err := json.Unmarshal(raw, &attribute); err != nil {
			return nil, "", fmt.Errorf("decode %s metadata: %w", identity.Name, err)
		}
		if attribute.ExplicitlySpecified == nil {
			return nil, "", fmt.Errorf("attribute %s lacks explicitlySpecified metadata; recapture with %s", attribute.Name, buildDepsCaptureFlags)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(raw, &fields); err != nil {
			return nil, "", fmt.Errorf("decode %s metadata fields: %w", attribute.Name, err)
		}
		if _, err := attributeField[bool](fields, "nodep", false); err != nil {
			return nil, "", fmt.Errorf("attribute %s nodep metadata: %w", attribute.Name, err)
		}
		if _, err := attributeField[string](fields, "sourceAspectName", false); err != nil {
			return nil, "", fmt.Errorf("attribute %s aspect metadata: %w", attribute.Name, err)
		}
		values, err := attributeValues(attribute.Type, raw, false)
		if err != nil {
			return nil, "", fmt.Errorf("decode %s attribute: %w", attribute.Name, err)
		}
		for index, label := range values {
			values[index] = normalizeLabel(label)
			if _, err := parseTargetGroup(values[index]); err != nil {
				return nil, "", fmt.Errorf("attribute %s label %q: %w", attribute.Name, label, err)
			}
		}
		slices.Sort(values)
		attribute.Labels = slices.Compact(values)
		if len(attribute.Labels) == 0 {
			attribute.Labels = nil
		}
		header, err := json.Marshal(attribute)
		if err != nil {
			return nil, "", fmt.Errorf("encode normalized attribute %s: %w", attribute.Name, err)
		}
		labels, err := json.Marshal(attribute.Labels)
		if err != nil {
			return nil, "", fmt.Errorf("encode normalized labels for %s: %w", attribute.Name, err)
		}
		signature := string(header) + string(labels)
		if previous, exists := metadata[attribute.Name]; exists && previous != signature {
			return nil, "", fmt.Errorf("conflicting metadata for attribute %s", attribute.Name)
		}
		metadata[attribute.Name] = signature
		if *attribute.ExplicitlySpecified && !attribute.Nodep && attribute.SourceAspectName == "" {
			for _, label := range attribute.Labels {
				selected[label] = true
			}
		}
	}
	if !marker {
		return nil, "", fmt.Errorf("build-deps-only requires a name attribute marker on every RULE; recapture with %s", buildDepsCaptureFlags)
	}
	var inputs []string
	for _, label := range rule.RuleInput {
		label = normalizeLabel(label)
		if selected[label] {
			inputs = append(inputs, label)
		}
	}
	slices.Sort(inputs)
	return slices.Compact(inputs), metadata["deps"] + metadata["srcs"], nil
}

// attributeValues decodes only the declared label-bearing field of an attribute.
// Empty repeated fields are omitted by Bazel; ordinary null LABEL values also
// omit stringValue. Explicit JSON null is not part of this capture contract.
func attributeValues(kind string, raw json.RawMessage, selectorEntry bool) ([]string, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil, fmt.Errorf("decode attribute value object: %w", err)
	}
	switch kind {
	case "LABEL":
		if err := validateValueFields(fields, "stringValue"); err != nil {
			return nil, fmt.Errorf("label encoding: %w", err)
		}
		value, err := attributeField[string](fields, "stringValue", selectorEntry)
		if err != nil {
			return nil, fmt.Errorf("decode label: %w", err)
		}
		if _, present := fields["stringValue"]; !present {
			return nil, nil
		}
		return []string{value}, nil
	case "LABEL_LIST":
		if err := validateValueFields(fields, "stringListValue"); err != nil {
			return nil, fmt.Errorf("label list encoding: %w", err)
		}
		values, err := attributeField[[]string](fields, "stringListValue", false)
		if err != nil {
			return nil, fmt.Errorf("decode label list: %w", err)
		}
		return values, nil
	case "LABEL_DICT_UNARY", "LABEL_LIST_DICT", "LABEL_KEYED_STRING_DICT":
		field := "labelDictUnaryValue"
		if kind == "LABEL_LIST_DICT" {
			// Bazel's formatter uses this field despite the protobuf's legacy labelListDictValue.
			field = "stringListDictValue"
		} else if kind == "LABEL_KEYED_STRING_DICT" {
			field = "labelKeyedStringDictValue"
		}
		if err := validateValueFields(fields, field); err != nil {
			return nil, fmt.Errorf("%s encoding: %w", kind, err)
		}
		entries, err := attributeField[[]map[string]json.RawMessage](fields, field, false)
		if err != nil {
			return nil, fmt.Errorf("decode %s dictionary: %w", kind, err)
		}
		var labels []string
		for index, entry := range entries {
			key, err := attributeField[string](entry, "key", true)
			if err != nil {
				return nil, fmt.Errorf("dictionary entry %d key: %w", index, err)
			}
			if kind == "LABEL_LIST_DICT" {
				values, err := attributeField[[]string](entry, "value", false)
				if err != nil {
					return nil, fmt.Errorf("dictionary entry %d values: %w", index, err)
				}
				labels = append(labels, values...)
			} else {
				value, err := attributeField[string](entry, "value", true)
				if err != nil {
					return nil, fmt.Errorf("dictionary entry %d value: %w", index, err)
				}
				if kind == "LABEL_KEYED_STRING_DICT" {
					labels = append(labels, key)
				} else {
					labels = append(labels, value)
				}
			}
		}
		return labels, nil
	case "SELECTOR_LIST":
		if err := validateValueFields(fields, "selectorList"); err != nil {
			return nil, fmt.Errorf("selector encoding: %w", err)
		}
		if selectorEntry {
			return nil, errors.New("nested selector lists are unsupported")
		}
		labels, err := selectorValues(fields)
		if err != nil {
			return nil, fmt.Errorf("decode selector list: %w", err)
		}
		return labels, nil
	case "STRING", "STRING_LIST", "INTEGER", "INTEGER_LIST", "BOOLEAN", "TRISTATE", "LICENSE", "DISTRIBUTION_SET", "OUTPUT", "OUTPUT_LIST", "STRING_DICT", "STRING_LIST_DICT", "NODEP_LABEL", "NODEP_LABEL_LIST":
		return nil, nil
	default:
		return nil, fmt.Errorf("unsupported dependency attribute type %q", kind)
	}
}

// selectorValues unions literal alternatives without treating condition keys as edges.
func selectorValues(fields map[string]json.RawMessage) ([]string, error) {
	selector, err := attributeField[map[string]json.RawMessage](fields, "selectorList", true)
	if err != nil {
		return nil, fmt.Errorf("read selectorList: %w", err)
	}
	kind, err := attributeField[string](selector, "type", true)
	if err != nil {
		return nil, fmt.Errorf("read selector type: %w", err)
	}
	// Validate the underlying type even when the selector has no alternatives.
	if kind == "SELECTOR_LIST" {
		return nil, errors.New("nested selector lists are unsupported")
	}
	if _, err := attributeValues(kind, json.RawMessage(`{}`), false); err != nil {
		return nil, fmt.Errorf("unsupported selector value type %q: %w", kind, err)
	}
	elements, err := attributeField[[]map[string]json.RawMessage](selector, "elements", true)
	if err != nil {
		return nil, fmt.Errorf("read selector elements: %w", err)
	}
	if len(elements) == 0 {
		return nil, errors.New("selector list has no elements")
	}
	var labels []string
	for index, element := range elements {
		entries, err := attributeField[[]json.RawMessage](element, "entries", true)
		if err != nil {
			return nil, fmt.Errorf("selector element %d: %w", index, err)
		}
		if len(entries) == 0 {
			return nil, fmt.Errorf("selector element %d has no entries", index)
		}
		for _, raw := range entries {
			var entry map[string]json.RawMessage
			if err := json.Unmarshal(raw, &entry); err != nil {
				return nil, fmt.Errorf("decode selector entry: %w", err)
			}
			condition, err := attributeField[string](entry, "label", true)
			if err != nil {
				return nil, fmt.Errorf("read selector condition: %w", err)
			}
			if _, err := parseTargetGroup(condition); err != nil {
				return nil, fmt.Errorf("invalid selector condition %q: %w", condition, err)
			}
			useDefault, err := attributeField[bool](entry, "isDefaultValue", false)
			if err != nil {
				return nil, fmt.Errorf("read selector default marker: %w", err)
			}
			if useDefault {
				// None serializes the underlying type's default: absent label,
				// empty label collection, or a non-label value such as STRING "".
				defaults, err := attributeValues(kind, raw, false)
				if err != nil {
					return nil, fmt.Errorf("read None selector branch default: %w", err)
				}
				if len(defaults) != 0 {
					return nil, errors.New("None selector branch contains dependency labels")
				}
				continue
			}
			values, err := attributeValues(kind, raw, true)
			if err != nil {
				return nil, fmt.Errorf("read selector branch %q: %w", condition, err)
			}
			labels = append(labels, values...)
		}
	}
	return labels, nil
}

// validateValueFields rejects mismatched or unsupported label value encodings.
func validateValueFields(fields map[string]json.RawMessage, expected string) error {
	for name := range fields {
		if name != expected && name != "isDefaultValue" && (strings.HasSuffix(name, "Value") || name == "selectorList") {
			return fmt.Errorf("unexpected value field %q", name)
		}
	}
	return nil
}

// attributeField decodes a typed protobuf field, distinguishing absent from malformed.
func attributeField[T any](fields map[string]json.RawMessage, name string, required bool) (T, error) {
	var value T
	raw, present := fields[name]
	if !present {
		if required {
			return value, fmt.Errorf("missing %s", name)
		}
		return value, nil
	}
	if bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return value, fmt.Errorf("null %s is not emitted by the required Bazel capture", name)
	}
	if err := json.Unmarshal(raw, &value); err != nil {
		return value, fmt.Errorf("decode %s: %w", name, err)
	}
	return value, nil
}

// writeBuildDepsComment identifies the opt-in structural selection in Mermaid.
func writeBuildDepsComment(output *strings.Builder, enabled bool) {
	if enabled {
		output.WriteString("%% build_deps_only=true\n")
	}
}
