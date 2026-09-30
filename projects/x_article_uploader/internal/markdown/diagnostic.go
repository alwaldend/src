package markdown

import (
	"fmt"
	"sort"
	"strings"
)

// Severity distinguishes a conversion outcome. A continuing diagnostic reports
// content that survives with lost formatting; a failing diagnostic reports
// content that would be lost, so conversion refuses to emit a document.
type Severity string

const (
	// Continuing reports a construct whose content is preserved.
	Continuing Severity = "continuing"
	// Failing reports a construct whose content would be lost.
	Failing Severity = "failing"
)

// Position is a source location within the post body, one-based.
type Position struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Diagnostic names an unrepresentable construct and where it came from.
type Diagnostic struct {
	Severity   Severity `json:"severity"`
	Code       string   `json:"code"`
	Message    string   `json:"message"`
	Source     string   `json:"source"`
	Position   Position `json:"position"`
	ByteOffset int      `json:"byte_offset"`
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%s:%d:%d: %s: %s", d.Source, d.Position.Line, d.Position.Column, d.Code, d.Message)
}

// positionAt converts a byte offset in the body into a one-based line and
// column within the source file. bodyLine is the file line the body begins on.
// The offset is a byte offset, because diagnostics also carry ByteOffset for
// slicing, but the column counts characters: a multi-byte character occupies
// one column, so the reported position is the one an editor shows.
func positionAt(body []byte, bodyLine, offset int) Position {
	if offset < 0 {
		offset = 0
	}
	if offset > len(body) {
		offset = len(body)
	}
	line := bodyLine
	column := 1
	for _, r := range string(body[:offset]) {
		if r == '\n' {
			line++
			column = 1
			continue
		}
		column++
	}
	return Position{Line: line, Column: column}
}

// DiagnosticError is the conversion failure carrying every failing diagnostic.
type DiagnosticError struct {
	Source      string
	Diagnostics []Diagnostic
}

func (e *DiagnosticError) Error() string {
	parts := make([]string, 0, len(e.Diagnostics))
	for _, d := range e.Diagnostics {
		parts = append(parts, d.String())
	}
	sort.Strings(parts)
	return fmt.Sprintf("%s: conversion failed: %s", e.Source, strings.Join(parts, "; "))
}
