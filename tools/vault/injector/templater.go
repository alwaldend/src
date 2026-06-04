package main

import (
	"bytes"
	"context"
	"fmt"
	"git.alwaldend.com/alwaldend/src/projects/al/pkg/fp"
	"strings"
	"text/template"
)

type templater struct {
}

type templateCtx struct {
	Res map[string]*ResourceResult
}

func (self *templater) Template(ctx context.Context, tpl string, data []*ResourceResult) fp.Either[string] {
	tmplCtx := &templateCtx{Res: map[string]*ResourceResult{}}
	for _, r := range data {
		tmplCtx.Res[r.Name] = r
	}
	tmpl, err := template.New("com.alwaldend.src.tools.vault.injector.templater").Funcs(
		template.FuncMap{
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
