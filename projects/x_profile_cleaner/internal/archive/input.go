package archive

import (
	"bytes"
	"encoding/json"
	"fmt"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/model"
)

// objectReader reads literal source keys and shares the first parsing failure.
// It is parser state, not a second archive record model.
type objectReader struct {
	// fields contains the current JSON object's source values.
	fields map[string]json.RawMessage
	// failure belongs to the complete source record, including nested objects.
	failure *error
}

// readObject accepts an object or an absent optional value without coercion.
func readObject(raw json.RawMessage, failure *error, required bool) objectReader {
	result := objectReader{failure: failure}
	if len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		if required {
			result.fail(fmt.Errorf("missing required source object"))
		}
		return result
	}
	if err := json.Unmarshal(raw, &result.fields); err != nil {
		result.fail(fmt.Errorf("decode source object: %w", err))
	}
	return result
}

// sourceRecord opens exactly the named standard X row envelope.
func sourceRecord(row json.RawMessage, name string) objectReader {
	var failure error
	return readObject(row, &failure, true).object(name, true)
}

// fail preserves the first failure without exposing partially decoded records.
func (r objectReader) fail(err error) {
	if err != nil && *r.failure == nil {
		*r.failure = err
	}
}

// object reads a nested object using the literal source field name.
func (r objectReader) object(name string, required bool) objectReader {
	if *r.failure != nil {
		return objectReader{failure: r.failure}
	}
	return readObject(r.fields[name], r.failure, required)
}

// optionalString reads only JSON strings; missing/null values become absent.
func (r objectReader) optionalString(name string) *string {
	if *r.failure != nil {
		return nil
	}
	raw := r.fields[name]
	if len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		r.fail(fmt.Errorf("decode source field %s: %w", name, err))
		return nil
	}
	return &value
}

// text reads an optional string with the source format's empty default.
func (r objectReader) text(name string) string {
	value := r.optionalString(name)
	if value == nil {
		return ""
	}
	return *value
}

// requiredString requires presence while allowing an explicitly empty string.
func (r objectReader) requiredString(name string) string {
	value := r.optionalString(name)
	if value == nil {
		r.fail(fmt.Errorf("missing required source field %s", name))
		return ""
	}
	return *value
}

// identifier checks canonical nonzero decimal IDs without numeric coercion.
func (r objectReader) identifier(name string, required bool) string {
	value := r.text(name)
	if (required || value != "") && !model.ValidID(value) {
		r.fail(fmt.Errorf("invalid source identifier %s", name))
	}
	return value
}

// array accepts a JSON array; absent/null collections become empty.
func (r objectReader) array(name string) []json.RawMessage {
	if *r.failure != nil || len(r.fields[name]) == 0 {
		return nil
	}
	var result []json.RawMessage
	if err := json.Unmarshal(r.fields[name], &result); err != nil {
		r.fail(fmt.Errorf("decode source array %s: %w", name, err))
		return nil
	}
	return result
}

// strings decodes an array of strings without converting numeric IDs.
func (r objectReader) strings(name string) []string {
	if *r.failure != nil || len(r.fields[name]) == 0 {
		return nil
	}
	var result []string
	if err := json.Unmarshal(r.fields[name], &result); err != nil {
		r.fail(fmt.Errorf("decode source strings %s: %w", name, err))
		return nil
	}
	return result
}

// boolean reads the optional source boolean; missing/null means false.
func (r objectReader) boolean(name string) bool {
	if *r.failure != nil || len(r.fields[name]) == 0 {
		return false
	}
	var result bool
	if err := json.Unmarshal(r.fields[name], &result); err != nil {
		r.fail(fmt.Errorf("decode source flag %s: %w", name, err))
	}
	return result
}

// integer decodes an exact nonnegative counter; missing/null values stay absent.
func (r objectReader) integer(name string) *int64 {
	raw := r.fields[name]
	if *r.failure != nil || len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null")) {
		return nil
	}
	value, err := number(raw)
	if err != nil {
		r.fail(fmt.Errorf("decode source counter %s: %w", name, err))
		return nil
	}
	return &value
}
