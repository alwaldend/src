package main

import (
	"fmt"
)

type OutputType string

const (
	OutputTypeJson     = "json"
	OutputTypeMarkdown = "markdown"
	OutputTypeTemplate = "template"
)

// String is used both by fmt.Print and by Cobra in help text
func (self *OutputType) String() string {
	return string(*self)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (self *OutputType) Set(v string) error {
	switch v {
	case OutputTypeJson, OutputTypeMarkdown, OutputTypeTemplate:
		*self = OutputType(v)
		return nil
	default:
		return fmt.Errorf(
			"invalid value: %s: not one of (%s, %s, %s)",
			v,
			OutputTypeJson, OutputTypeMarkdown, OutputTypeTemplate,
		)
	}
}

// Type is only used in help text
func (self *OutputType) Type() string {
	return "ParseOutputType"
}

type ParsePathsResult struct {
	Dir     string    `json:"dir"`
	Readmes []*Readme `json:"readmes"`
}

type Readme struct {
	Path             string         `json:"path"`
	PathRelative     string         `json:"path_relative"`
	Dir              string         `json:"dir"`
	DirRelative      string         `json:"dir_relative"`
	DirRelativeSplit []string       `json:"dir_relative_split"`
	DirName          string         `json:"dir_name"`
	Data             map[string]any `json:"data"`
}

type ParseConfig struct {
	OutputType OutputType
	UseGit     bool
	ReadmeName string
	RootPath   string
	Exclude    []string
}

type Config struct {
	ChangeDir string
}
