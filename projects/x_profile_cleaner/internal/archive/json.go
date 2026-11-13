package archive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	// maxJSONInputBytes bounds a standalone JSON projection input.
	maxJSONInputBytes = 16 * 1024 * 1024
	// maxJSONDepth bounds untrusted archive nesting before typed decoding.
	maxJSONDepth = 64
	// maxJSONValues bounds untrusted JSON work independently of byte size.
	maxJSONValues = 1000000
)

// ValidateJSON rejects ambiguous duplicate keys, excessive nesting and trailing data.
func ValidateJSON(raw []byte) error {
	if len(raw) == 0 || len(raw) > maxJSONInputBytes || !utf8.Valid(raw) {
		return fmt.Errorf("invalid JSON size or encoding")
	}
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return fmt.Errorf("empty archive JSON")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	values := 0
	var visit func(int) error
	visit = func(depth int) error {
		values++
		if depth > maxJSONDepth || values > maxJSONValues {
			return fmt.Errorf("JSON complexity bound exceeded")
		}
		token, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("read JSON token: %w", err)
		}
		switch token {
		case json.Delim('{'):
			keys := map[string]bool{}
			for decoder.More() {
				key, err := decoder.Token()
				if err != nil {
					return fmt.Errorf("read JSON key: %w", err)
				}
				name, ok := key.(string)
				if !ok || keys[name] {
					return fmt.Errorf("invalid or duplicate JSON key")
				}
				keys[name] = true
				if err := visit(depth + 1); err != nil {
					return fmt.Errorf("read JSON member: %w", err)
				}
			}
			if _, err := decoder.Token(); err != nil {
				return fmt.Errorf("close JSON object: %w", err)
			}
		case json.Delim('['):
			for decoder.More() {
				if err := visit(depth + 1); err != nil {
					return fmt.Errorf("read JSON element: %w", err)
				}
			}
			if _, err := decoder.Token(); err != nil {
				return fmt.Errorf("close JSON array: %w", err)
			}
		default:
			if _, ok := token.(json.Delim); ok {
				return fmt.Errorf("unexpected JSON delimiter")
			}
		}
		return nil
	}
	if err := visit(0); err != nil {
		return fmt.Errorf("read JSON document: %w", err)
	}
	if _, err := decoder.Token(); err != io.EOF {
		return fmt.Errorf("trailing JSON data")
	}
	return nil
}
