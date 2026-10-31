// Package main runs scan-only TruffleHog Git scans, validates exact historical
// finding suppressions, and preserves original results as annotated JSON events.
// It owns the scanner process and its temporary files; callers supply the tools,
// repository location, and policy. It never verifies credentials or changes Git.
package main
