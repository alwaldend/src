package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
)

type Templater struct{}

type TemplateDataFile struct {
	// Data file path
	Path string
	// Data file basename
	Basename string
	// Data file dirname
	Dirname string
	// Data file lines
	Lines []string
	// Parsed data
	Data any
}

type TemplateContext struct {
	Data []*TemplateDataFile
}

func timestampToDate(timestamp float64) (string, error) {
	tm := time.Unix(int64(timestamp), 0)
	return tm.String(), nil
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
	funcMap := template.FuncMap{
		"timestamp_to_date": timestampToDate,
	}
	templateObj, err := template.New("template").Funcs(funcMap).Parse(string(templateBytes))
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
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read %s: %w", path, err)
	}
	expr, err := regexp.Compile("\r?\n")
	if err != nil {
		return nil, fmt.Errorf("could not compile regex: %w", err)
	}
	fileString := string(fileBytes)
	fileLines := expr.Split(fileString, -1)
	fileLinesCount := len(fileLines)
	if fileLinesCount != 0 && fileLines[fileLinesCount-1] == "" {
		fileLines = fileLines[:fileLinesCount-1]
		fileLinesCount -= 1
	}
	extension := filepath.Ext(path)
	var data any
	switch extension {
	case ".toml":
		_, err := toml.Decode(fileString, &data)
		if err != nil {
			return nil, err
		}
	case ".json":
		err := json.Unmarshal(fileBytes, &data)
		if err != nil {
			return nil, err
		}
	case ".txt":
	default:
		return nil, fmt.Errorf("unsupported extension %s: %s", extension, path)
	}
	return &TemplateDataFile{
		Path:     path,
		Dirname:  filepath.Dir(path),
		Basename: filepath.Base(path),
		Data:     data,
		Lines:    fileLines,
	}, nil
}
