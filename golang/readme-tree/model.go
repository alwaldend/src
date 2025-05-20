package main

import (
	"fmt"
	"io"
)

type ParseOutputType string

const (
	ParseOutputTypeJson     = "json"
	ParseOutputTypeMarkdown = "markdown"
	ParseOutputTypeTemplate = "template"
)

// String is used both by fmt.Print and by Cobra in help text
func (self *ParseOutputType) String() string {
	return string(*self)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (self *ParseOutputType) Set(v string) error {
	switch v {
	case ParseOutputTypeJson, ParseOutputTypeMarkdown, ParseOutputTypeTemplate:
		*self = ParseOutputType(v)
		return nil
	default:
		return fmt.Errorf(
			"invalid value: %s: not one of (%s, %s, %s)",
			v,
			ParseOutputTypeJson, ParseOutputTypeMarkdown, ParseOutputTypeTemplate,
		)
	}
}

// Type is only used in help text
func (self *ParseOutputType) Type() string {
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
	OutputType ParseOutputType
}

type Config struct {
	UseGit     bool
	ChangeDir  string
	ReadmeName string
	Parse      *ParseConfig
}

type State struct {
	Getenv         func(string) string
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	Config         *Config
}
