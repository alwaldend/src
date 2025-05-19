package main

import "io"

type ParsePathsResult struct {
	Dir     string    `json:"dir"`
	Readmes []*Readme `json:"readmes"`
}

type Readme struct {
	Path         string         `json:"path"`
	RelativePath string         `json:"relative_path"`
	Data         map[string]any `json:"data"`
}

type State struct {
	Getenv         func(string) string
	Stdin          io.Reader
	Stdout, Stderr io.Writer
}
