package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"

	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
)

type Templater struct {
}

type templateCtx struct {
	Res   map[string]*ResourceResult
	Last  *ResourceResult
	Extra map[string]any
}

func (self *Templater) Template(ctx context.Context, tpl string, data []*ResourceResult, extra map[string]any) fp.Either[string] {
	tmplCtx := &templateCtx{Res: map[string]*ResourceResult{}, Extra: extra}
	for _, r := range data {
		tmplCtx.Res[r.Name] = r
		tmplCtx.Last = r
	}
	tmpl, err := template.New("com.alwaldend.src.tools.vault.injector.Templater.Template").Funcs(
		template.FuncMap{
			"b64decode": func(val string) (string, error) {
				res, err := base64.StdEncoding.DecodeString(val)
				if err != nil {
					return "", fmt.Errorf("could not decode base64 string: %w", err)
				}
				return string(res), nil
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
		return fp.Left[string](fmt.Errorf("could not parse template '%s': %w", tpl, err))
	}
	var buff bytes.Buffer
	err = tmpl.Execute(&buff, tmplCtx)
	if err != nil {
		return fp.Left[string](fmt.Errorf("could not execute template '%s': %w", tpl, err))
	}
	return fp.Right(buff.String())
}
