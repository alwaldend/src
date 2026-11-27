package main

import (
	"encoding/json"
	"fmt"
	html_template "html/template"
	"maps"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// templaterFuncMap registers the stable names available in template expressions.
func templaterFuncMap() template.FuncMap {
	return template.FuncMap{
		"timestamp_to_date": timestampToDate,
		"to_json":           toJson,
		"to_json_indent":    toJsonIndent,
		"to_html_table":     toHtmlTable,
		"last":              sequenceEnd(true),
		"set_map_key":       setMapKey,
		"unset_map_key":     unsetMapKey,
		"split":             split,
		"sliceString":       sliceString,
		"indent":            indent,
		"html_escape":       html_template.HTMLEscapeString,
		"basename":          filepath.Base,
		"dirname":           filepath.Dir,
		"first":             sequenceEnd(false),
	}
}

// sliceString returns the byte range requested by a template.
func sliceString(val string, start, end int) (string, error) {
	if start < 0 || end < start || end > len(val) {
		return "", fmt.Errorf("invalid byte range [%d:%d] for string of length %d", start, end, len(val))
	}
	return val[start:end], nil
}

// indent inserts the supplied indentation immediately before each normalized newline.
func indent(val string, indent string) string {
	lines := strings.Split(strings.ReplaceAll(val, "\r\n", "\n"), "\n")
	return strings.Join(lines, indent+"\n")
}

// split divides a string at every occurrence of the separator.
func split(val string, sep string) []string {
	return strings.Split(val, sep)
}

// sequenceEnd selects the first or last element of a template sequence, including typed slices.
func sequenceEnd(last bool) func(any) (any, error) {
	return func(sequence any) (any, error) {
		if sequence == nil {
			return nil, nil
		}
		value := reflect.ValueOf(sequence)
		if value.Kind() != reflect.Array && value.Kind() != reflect.Slice {
			return nil, fmt.Errorf("expected a sequence, got %T", sequence)
		}
		if value.Len() == 0 {
			return nil, nil
		}
		index := 0
		if last {
			index = value.Len() - 1
		}
		return value.Index(index).Interface(), nil
	}
}

// setMapKey assigns a template map entry and reports whether a map was supplied.
func setMapKey(mapping, key, value any) (bool, error) {
	switch entries := mapping.(type) {
	case nil:
		return false, nil
	case map[string]any:
		name, ok := key.(string)
		if !ok {
			return false, fmt.Errorf("expected a string key, got %T", key)
		}
		if entries == nil {
			return false, nil
		}
		entries[name] = value
	case map[any]any:
		if key != nil && !reflect.ValueOf(key).Comparable() {
			return false, fmt.Errorf("invalid map key of type %T", key)
		}
		if entries == nil {
			return false, nil
		}
		entries[key] = value
	default:
		return false, fmt.Errorf("expected a data map, got %T", mapping)
	}
	return true, nil
}

// unsetMapKey removes the requested keys from a decoded data map.
func unsetMapKey(mapping any, keys ...any) (bool, error) {
	switch entries := mapping.(type) {
	case nil:
		return false, nil
	case map[string]any:
		for _, key := range keys {
			name, ok := key.(string)
			if !ok {
				return false, fmt.Errorf("expected a string key, got %T", key)
			}
			delete(entries, name)
		}
		return entries != nil, nil
	case map[any]any:
		for _, key := range keys {
			if key != nil && !reflect.ValueOf(key).Comparable() {
				return false, fmt.Errorf("invalid map key of type %T", key)
			}
			delete(entries, key)
		}
		return entries != nil, nil
	default:
		return false, fmt.Errorf("expected a data map, got %T", mapping)
	}
}

// timestampToDate formats a Unix timestamp in local time.
func timestampToDate(timestamp int64) (string, error) {
	tm := time.Unix(int64(timestamp), 0)
	return tm.String(), nil
}

// toJson encodes a template value as JSON with no indentation prefix.
func toJson(val any) (string, error) {
	result, err := toJsonIndent(val, "", "")
	if err != nil {
		return "", fmt.Errorf("encode template value: %w", err)
	}
	return result, nil
}

// toHtmlTable formats nested values as HTML lists and tables.
func toHtmlTable(val any) (string, error) {
	result, err := htmlTable(val, 0)
	if err != nil {
		return "", fmt.Errorf("render HTML table: %w", err)
	}
	return result, nil
}

// htmlTable bounds recursion even when template map operations construct a cycle.
func htmlTable(val any, depth int) (string, error) {
	if depth >= 128 {
		return "", fmt.Errorf("HTML table nesting exceeds 128 levels")
	}
	switch valueAsserted := val.(type) {
	case map[string]any:
		res := []string{
			`<table style="width: 100;">`,
		}
		for _, key := range slices.Sorted(maps.Keys(valueAsserted)) {
			curVal := valueAsserted[key]
			curValTable, err := htmlTable(curVal, depth+1)
			if err != nil {
				return "", fmt.Errorf("could not convert key %s to table: %w", key, err)
			}
			res = append(res, fmt.Sprintf("<tr><th>%s</th><td>%s</td></tr>", key, curValTable))
		}
		res = append(res, "</table>")
		return strings.Join(res, "\n"), nil
	case []any:
		res := []string{}
		for index, curVal := range valueAsserted {
			curValTable, err := htmlTable(curVal, depth+1)
			if err != nil {
				return "", fmt.Errorf("convert item %d to HTML: %w", index, err)
			}
			res = append(res, fmt.Sprintf("<li>%s</li>", curValTable))

		}
		slices.Sort(res)
		return fmt.Sprintf("<ul>%s</ul>", strings.Join(res, "\n")), nil
	case string:
		return fmt.Sprintf("<code>%s</code>", valueAsserted), nil
	case int:
		return strconv.Itoa(valueAsserted), nil
	case bool:
		return strconv.FormatBool(valueAsserted), nil
	case float64:
		return strconv.FormatFloat(valueAsserted, 'f', -1, 64), nil
	default:
		return fmt.Sprint(valueAsserted), nil
	}
}

// toJsonIndent encodes JSON with the requested line prefix and indentation.
func toJsonIndent(val any, prefix string, indent string) (string, error) {
	res, err := json.MarshalIndent(val, prefix, indent)
	if err != nil {
		return "", fmt.Errorf("marshal JSON: %w", err)
	}
	return string(res), nil
}
