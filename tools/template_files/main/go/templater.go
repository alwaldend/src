package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
)

type Templater struct{}

// It has to be 7, otherwise the disk cache will start throwing
// 'permission denied' errors
var FILE_MODE = os.FileMode(0x744)

type TemplateDataFile struct {
	// Data file path
	Path string
	// Data file basename
	Basename string
	// Basename without the extension
	BasenameWithoutExt string
	// File extension (can be overriden using CLI)
	Extension string
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

func NewTemplater() *Templater {
	return &Templater{}
}

func (self *Templater) TemplateFiles(dataPaths []string, templatePath string, outputPath string, extension string) error {
	ctx := &TemplateContext{}
	for _, path := range dataPaths {
		data, err := self.loadData(path, extension)
		if err != nil {
			return fmt.Errorf("could not load data %s: %w", path, err)
		}
		ctx.Data = append(ctx.Data, data)
	}
	templateBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("could not read template %s: %w", templatePath, err)
	}
	templateObj, err := template.New("template").Funcs(TemplaterFuncMap()).Parse(string(templateBytes))
	if err != nil {
		return fmt.Errorf("could not parse template %s: %w", templatePath, err)
	}
	var templateBuffer bytes.Buffer
	err = templateObj.Execute(&templateBuffer, ctx)
	if err != nil {
		return fmt.Errorf("could not execute template %s: %w", templatePath, err)
	}
	err = os.WriteFile(outputPath, templateBuffer.Bytes(), FILE_MODE)
	if err != nil {
		return fmt.Errorf("could not write output %s: %w", outputPath, err)
	}
	return nil
}

func (self *Templater) loadData(path string, extension string) (*TemplateDataFile, error) {
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
	if extension == "" {
		extension = filepath.Ext(path)
	}
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
	case ".ndjson":
		dataSlice := []any{}
		for _, line := range fileLines {
			var lineJson any
			err := json.Unmarshal([]byte(line), &lineJson)
			if err != nil {
				return nil, err
			}
			dataSlice = append(dataSlice, lineJson)
		}
		data = dataSlice
	case ".yaml":
		err := yaml.Unmarshal(fileBytes, &data)
		if err != nil {
			return nil, err
		}
	case ".txt":
	default:
		return nil, fmt.Errorf("unsupported extension %s: %s", extension, path)
	}
	return &TemplateDataFile{
		Path:               path,
		Dirname:            filepath.Dir(path),
		Basename:           filepath.Base(path),
		Extension:          extension,
		BasenameWithoutExt: strings.TrimRight(filepath.Base(path), filepath.Ext(path)),
		Data:               data,
		Lines:              fileLines,
	}, nil
}
