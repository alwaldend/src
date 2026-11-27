package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
)

// outputMode supplies the creation permissions for new rendered files before umask.
const outputMode os.FileMode = 0o504

// maxInputBytes limits the combined size of the template and all data files.
const maxInputBytes = 64 << 20

// maxOutputBytes limits the retained render before publishing it.
const maxOutputBytes = 64 << 20

// maxDataLines bounds the line and NDJSON indexes shared across the data files.
const maxDataLines = 1 << 20

// renderBuffer holds the uncommitted render until execution succeeds.
type renderBuffer struct {
	// content holds at most maxOutputBytes of rendered text.
	content bytes.Buffer
}

// Write rejects a chunk before it could exceed the output limit.
func (buffer *renderBuffer) Write(chunk []byte) (int, error) {
	if len(chunk) > maxOutputBytes-buffer.content.Len() {
		return 0, fmt.Errorf("rendered output exceeds %d bytes", maxOutputBytes)
	}
	n, err := buffer.content.Write(chunk)
	if err != nil {
		return n, fmt.Errorf("buffer rendered output: %w", err)
	}
	return n, nil
}

// readInput reads at most the remaining input budget and accounts for growing files.
func readInput(path string, remaining int64) (content []byte, result error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close %q: %w", path, err))
		}
	}()
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("inspect %q: %w", path, err)
	}
	if info.Size() <= remaining {
		content, err = io.ReadAll(io.LimitReader(file, remaining+1))
		if err != nil {
			return nil, fmt.Errorf("read %q: %w", path, err)
		}
		if int64(len(content)) <= remaining {
			return content, nil
		}
	}
	return nil, fmt.Errorf("combined input exceeds %d bytes while reading %q", maxInputBytes, path)
}

// templateDataFile exposes a data file and its metadata to text/template.
type templateDataFile struct {
	// Path is the name supplied by the caller.
	Path string
	// Basename is the final component of Path.
	Basename string
	// BasenameWithoutExt removes the final filename extension.
	BasenameWithoutExt string
	// Extension selects the parser, including any caller override.
	Extension string
	// Dirname is the parent directory of Path.
	Dirname string
	// Lines holds the input split on LF or CRLF.
	Lines []string
	// Data holds the decoded value for structured inputs.
	Data any
}

// templateContext is the root value visible to template expressions.
type templateContext struct {
	// DataFiles follows the order of the caller's data arguments.
	DataFiles []*templateDataFile
}

// templateFiles loads the inputs and publishes output after successful rendering.
func templateFiles(ctx context.Context, dataPaths []string, templatePath, outputPath, extension string) error {
	templateBytes, err := readInput(templatePath, maxInputBytes)
	if err != nil {
		return fmt.Errorf("load template %q: %w", templatePath, err)
	}
	remainingBytes := int64(maxInputBytes - len(templateBytes))
	remainingLines := maxDataLines
	dataContext := templateContext{DataFiles: make([]*templateDataFile, 0, len(dataPaths))}
	for _, path := range dataPaths {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("load template inputs: %w", err)
		}
		content, err := readInput(path, remainingBytes)
		if err != nil {
			return fmt.Errorf("load template data: %w", err)
		}
		remainingBytes -= int64(len(content))
		data, err := loadData(path, extension, content, remainingLines)
		if err != nil {
			return fmt.Errorf("could not load data %s: %w", path, err)
		}
		remainingLines -= len(data.Lines)
		dataContext.DataFiles = append(dataContext.DataFiles, data)
	}
	templateObj, err := template.New("template").Funcs(templaterFuncMap()).Parse(string(templateBytes))
	if err != nil {
		return fmt.Errorf("could not parse template %s: %w", templatePath, err)
	}
	var templateBuffer renderBuffer
	err = templateObj.Execute(&templateBuffer, dataContext)
	if err != nil {
		return fmt.Errorf("could not execute template %s: %w", templatePath, err)
	}
	err = os.WriteFile(outputPath, templateBuffer.content.Bytes(), outputMode)
	if err != nil {
		return fmt.Errorf("could not write output %s: %w", outputPath, err)
	}
	return nil
}

// loadData decodes one data file and preserves its normalized lines for templates.
func loadData(path, extension string, fileBytes []byte, remainingLines int) (*templateDataFile, error) {
	fileString := string(fileBytes)
	lines := bytes.Count(fileBytes, []byte{'\n'})
	if len(fileBytes) != 0 && fileBytes[len(fileBytes)-1] != '\n' {
		lines++
	}
	if lines > remainingLines {
		return nil, fmt.Errorf("data files exceed %d lines", maxDataLines)
	}
	fileLines := strings.Split(strings.ReplaceAll(fileString, "\r\n", "\n"), "\n")
	fileLinesCount := len(fileLines)
	if fileLinesCount != 0 && fileLines[fileLinesCount-1] == "" {
		fileLines = fileLines[:fileLinesCount-1]
	}
	if extension == "" {
		extension = filepath.Ext(path)
	}
	var data any
	switch extension {
	case ".toml":
		_, err := toml.Decode(fileString, &data)
		if err != nil {
			return nil, fmt.Errorf("decode TOML: %w", err)
		}
	case ".json":
		err := json.Unmarshal(fileBytes, &data)
		if err != nil {
			return nil, fmt.Errorf("decode JSON: %w", err)
		}
	case ".ndjson":
		dataSlice := make([]any, 0, len(fileLines))
		for index, line := range fileLines {
			var lineJson any
			err := json.Unmarshal([]byte(line), &lineJson)
			if err != nil {
				return nil, fmt.Errorf("decode JSON line %d: %w", index+1, err)
			}
			dataSlice = append(dataSlice, lineJson)
		}
		data = dataSlice
	case ".yaml":
		err := yaml.Unmarshal(fileBytes, &data)
		if err != nil {
			return nil, fmt.Errorf("decode YAML: %w", err)
		}
	case ".txt":
	default:
		return nil, fmt.Errorf("unsupported extension %s: %s", extension, path)
	}
	return &templateDataFile{
		Path:               path,
		Dirname:            filepath.Dir(path),
		Basename:           filepath.Base(path),
		Extension:          extension,
		BasenameWithoutExt: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Data:               data,
		Lines:              fileLines,
	}, nil
}
