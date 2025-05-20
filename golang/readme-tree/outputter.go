package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
)

//go:embed template-markdown.txt
var templateMarkdown string

type Outputter struct {
	templateMarkdown string
	config           *Config
	stdout           io.Writer
}

func NewOutputter(config *Config, stdout io.Writer) *Outputter {
	return &Outputter{
		config:           config,
		templateMarkdown: templateMarkdown,
		stdout:           stdout,
	}
}

func (self *Outputter) Output(data any, output OutputType) error {
	switch output {
	case OutputTypeJson:
		return self.OutputJson(data)
	case OutputTypeMarkdown:
		return self.OutputMarkdown(data)
	default:
		return fmt.Errorf("output type %s is not supported", output)
	}
}

func (self *Outputter) OutputMarkdown(data any) error {
	template, err := template.New("markdown").Parse(self.templateMarkdown)
	if err != nil {
		return err
	}
	return self.OutputTemplate(template, data)
}

func (self *Outputter) OutputTemplate(t *template.Template, data any) error {
	var buffer bytes.Buffer
	err := t.Execute(&buffer, data)
	if err != nil {
		return err
	}
	return self.OutputString(buffer.String())
}

func (self *Outputter) OutputJson(data any) error {
	resultJson, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("could not serialise data: %w", err)
	}
	return self.OutputString(string(resultJson))
}

func (self *Outputter) OutputString(data string) error {
	_, err := fmt.Fprintf(self.stdout, "%s\n", data)
	if err != nil {
		return fmt.Errorf("could not print: %w", err)
	}
	return nil
}
