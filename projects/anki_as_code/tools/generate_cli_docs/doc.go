// Package main generates a short Markdown CLI index and one complete help page
// per root subcommand from Anki as code's Cobra tree. It renders help without
// executing application workflows or starting the native server; Bazel owns
// the fresh generated documentation directory.
package main
