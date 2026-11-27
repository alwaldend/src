package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
)

// templateMarkdown formats directory and title fields from README metadata.
//
//go:embed template-markdown.txt
var templateMarkdown string

// render writes the collected tree in the selected format.
func render(out io.Writer, result parseResult, format outputFormat) error {
	switch format {
	case jsonFormat:
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", "    ")
		if err := encoder.Encode(result); err != nil {
			return fmt.Errorf("encode JSON: %w", err)
		}
	case markdownFormat:
		tmpl, err := template.New("markdown").Parse(templateMarkdown)
		if err != nil {
			return fmt.Errorf("parse embedded Markdown template: %w", err)
		}
		if err := tmpl.Execute(out, result); err != nil {
			return fmt.Errorf("write Markdown tree: %w", err)
		}
		if _, err := fmt.Fprintln(out); err != nil {
			return fmt.Errorf("terminate Markdown output: %w", err)
		}
	default:
		return fmt.Errorf("unsupported output format %q", format)
	}
	return nil
}
