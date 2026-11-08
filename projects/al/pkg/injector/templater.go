package injector

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
)

// templater evaluates dependency-aware templates with redacted errors.
type templater struct{}

// templateCtx exposes declared dependency results and configured resource input.
type templateCtx struct {
	// Res exposes declared dependency results by name.
	Res map[string]*ResourceResult
	// Last is the final declared dependency result.
	Last *ResourceResult
	// Data supplies the current resource's configured input values.
	Data any
}

// Template renders local input and ordered dependencies with value-free errors.
func (self *templater) Template(ctx context.Context, tpl string, dependencies []*ResourceResult, data any) fp.Either[string] {
	tmplCtx := &templateCtx{Res: map[string]*ResourceResult{}, Data: data}
	for _, r := range dependencies {
		tmplCtx.Res[r.Name] = r
		tmplCtx.Last = r
	}
	tmpl, err := template.New("resource").Funcs(
		template.FuncMap{
			"to_json_indent": func(val any, prefix string, indent string) (string, error) {
				res, err := json.MarshalIndent(val, prefix, indent)
				if err != nil {
					return "", fmt.Errorf("could not marshal json: %w", err)
				}
				return string(res), nil
			},
			"to_json": func(val any) (string, error) {
				res, err := json.Marshal(val)
				if err != nil {
					return "", fmt.Errorf("could not marshal json: %w", err)
				}
				return string(res), nil
			},
			"b64decode": func(val string) (string, error) {
				res, err := base64.StdEncoding.DecodeString(val)
				if err != nil {
					return "", fmt.Errorf("could not decode base64 string: %w", err)
				}
				return string(res), nil
			},
			"b64encode": func(val string) (string, error) {
				return base64.StdEncoding.EncodeToString([]byte(val)), nil
			},
			"join": func(elems []any, sep string) string {
				elemsString := make([]string, len(elems))
				for _, elem := range elems {
					elemsString = append(elemsString, elem.(string))
				}
				return strings.Join(elemsString, sep)
			},
		},
	).Option("missingkey=error").Parse(tpl)
	if err != nil {
		return fp.Left[string](fmt.Errorf("could not parse template (details redacted)"))
	}
	var buff bytes.Buffer
	err = tmpl.Execute(&buff, tmplCtx)
	if err != nil {
		return fp.Left[string](fmt.Errorf("could not execute template (details redacted)"))
	}
	return fp.Right(buff.String())
}
