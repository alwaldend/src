package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
)

//go:embed template-markdown.txt
var templateMarkdown string

type Outputter struct {
	state            *State
	templateMarkdown string
}

func NewOutputter(state *State) *Outputter {
	return &Outputter{state: state, templateMarkdown: templateMarkdown}
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
	return self.Output(buffer.String())
}

func (self *Outputter) OutputJson(data any) error {
	resultJson, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("could not serialise data: %w", err)
	}
	return self.Output(string(resultJson))
}

func (self *Outputter) Output(data string) error {
	_, err := fmt.Fprintf(self.state.Stdout, "%s\n", data)
	if err != nil {
		return fmt.Errorf("could not print: %w", err)
	}
	return nil
}
