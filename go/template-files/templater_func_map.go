package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var TemplaterFuncMap = template.FuncMap{
	"timestamp_to_date": timestampToDate,
	"to_json":           toJson,
	"to_json_indent":    toJsonIndent,
	"to_html_table":     toHtmlTable,
	"last":              last[string],
	"set_map_key":       setMapKey[string],
	"unset_map_key":     unsetMapKey[string],
	"split":             split,
}

func split(val string, sep string) []string {
	return strings.Split(val, sep)
}

func last[T any](val []T) *T {
	if val != nil && len(val) != 0 {
		return &val[len(val)-1]
	}
	return nil
}

func setMapKey[T comparable](val map[T]any, key T, value any) bool {
	if val == nil {
		return false
	}
	val[key] = value
	return true
}

func unsetMapKey[T comparable](val map[T]any, keys ...T) bool {
	if val == nil {
		return false
	}
	for _, key := range keys {
		delete(val, key)
	}
	return true
}

func timestampToDate(timestamp float64) (string, error) {
	tm := time.Unix(int64(timestamp), 0)
	return tm.String(), nil
}

func toJson(val any) (string, error) {
	return toJsonIndent(val, "", "")
}

func toHtmlTable(val any) (string, error) {
	switch valueAsserted := val.(type) {
	case map[string]any:
		res := []string{
			`<table style="width: 100;">`,
		}
		for _, key := range slices.Sorted(maps.Keys(valueAsserted)) {
			curVal, ok := valueAsserted[key]
			if !ok {
				return "", fmt.Errorf("could not find key %s", curVal)
			}
			curValTable, err := toHtmlTable(curVal)
			if err != nil {
				return "", fmt.Errorf("could not convert key %s to table: %w", key, err)
			}
			res = append(res, fmt.Sprintf("<tr><th>%s</th><td>%s</td></tr>", key, curValTable))
		}
		res = append(res, "</table>")
		return strings.Join(res, "\n"), nil
	case []any:
		res := []string{}
		for _, curVal := range valueAsserted {
			curValTable, err := toHtmlTable(curVal)
			if err != nil {
				return "", fmt.Errorf("could not convert list %s to table: %w", valueAsserted, err)
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
		return fmt.Sprintf("%s", valueAsserted), nil
	}
}

func toJsonIndent(val any, prefix string, indent string) (string, error) {
	res, err := json.MarshalIndent(val, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
