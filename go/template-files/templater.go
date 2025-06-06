package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/BurntSushi/toml"
)

type Templater struct{}

type TemplateDataFile struct {
	Path string
	Data any
}

type TemplateContext struct {
	Data []*TemplateDataFile
}

func NewTemplater() *Templater {
	return &Templater{}
}

func (self *Templater) TemplateFiles(dataPaths []string, templatePath string, outputPath string) error {
	ctx := &TemplateContext{}
	for _, path := range dataPaths {
		data, err := self.loadData(path)
		if err != nil {
			return fmt.Errorf("could not load data %s: %w", path, err)
		}
		ctx.Data = append(ctx.Data, data)
	}
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("could not read template %s: %w", templatePath, err)
	}
	templateObj, err := template.New("template").Parse(string(templateBytes))
	if err != nil {
		return fmt.Errorf("could not parse template %s: %w", templatePath, err)
	}
	var templateBuffer bytes.Buffer
	err = templateObj.Execute(&templateBuffer, ctx)
	if err != nil {
		return fmt.Errorf("could not execute template %s: %w", templatePath, err)
	}
	err = os.WriteFile(outputPath, templateBuffer.Bytes(), os.FileMode(0x655))
	if err != nil {
		return fmt.Errorf("could not write output %s: %w", outputPath, err)
	}
	return nil
}

func (self *Templater) loadData(path string) (*TemplateDataFile, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", path, err)
	}
	extension := filepath.Ext(path)
	var data any
	switch extension {
	case ".toml":
		_, err := toml.Decode(string(file), &data)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported extension %s: %s", extension, path)
	}
	return &TemplateDataFile{Path: path, Data: data}, nil
}
