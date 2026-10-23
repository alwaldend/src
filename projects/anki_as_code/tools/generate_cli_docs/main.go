package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/cli"
	"github.com/spf13/cobra"
)

// documentationApplication permits command construction while leaving workflows
// unavailable. Rendering help must never invoke the embedded application.
type documentationApplication struct {
	// App supplies the workflow interface without a native implementation.
	ankiascodeapi.App
}

// main runs the reference generator with an explicit build-owned output directory.
func main() {
	var output string
	command := &cobra.Command{
		Use: "generate_cli_docs", Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := generate(output); err != nil {
				return fmt.Errorf("generate CLI reference: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&output, "output", "", "Generated Markdown output directory")
	if err := command.MarkFlagRequired("output"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// generate writes a short index and one complete help page per root subcommand
// into a Bazel-owned directory without calling application workflow handlers.
func generate(output string) error {
	root := cli.New(documentationApplication{}, nil)
	root.InitDefaultHelpCmd()
	if err := os.MkdirAll(output, 0o755); err != nil {
		return fmt.Errorf("create reference directory %q: %w", output, err)
	}
	var index bytes.Buffer
	index.WriteString("Choose a command to see its usage, options, and complete nested help.\n\n")
	for _, command := range root.Commands() {
		if command.Hidden {
			continue
		}
		fmt.Fprintf(&index, "- [`%s`](%s/) — %s\n", command.Name(), command.Name(), command.Short)
		var page bytes.Buffer
		if err := renderCommand(&page, command); err != nil {
			return fmt.Errorf("render command tree %q: %w", command.Name(), err)
		}
		if err := writePage(output, command.Name()+".md", command.Name(), command.Short, page.Bytes()); err != nil {
			return fmt.Errorf("write command page %q: %w", command.Name(), err)
		}
	}
	if err := writePage(output, "_index.md", "CLI reference", root.Short, index.Bytes()); err != nil {
		return fmt.Errorf("write reference index: %w", err)
	}
	return nil
}

// writePage adds shared Hugo metadata and writes one generated Markdown page.
func writePage(directory, name, title, description string, body []byte) error {
	var page bytes.Buffer
	fmt.Fprintf(&page, "---\ntitle: %q\ndescription: %q\nal_source_path: projects/anki_as_code/tools/generate_cli_docs/main.go\n", title, description)
	if name == "_index.md" {
		page.WriteString("no_list: true\n")
	}
	page.WriteString("---\n\n")
	page.Write(body)
	path := filepath.Join(directory, name)
	if err := os.WriteFile(path, page.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write reference page %q: %w", path, err)
	}
	return nil
}

// renderCommand appends each visible command's help in sorted command order.
func renderCommand(page *bytes.Buffer, command *cobra.Command) error {
	if command.Hidden {
		return nil
	}
	command.InitDefaultHelpFlag()
	fmt.Fprintf(page, "## %s\n\n", command.CommandPath())
	if command.DisableFlagParsing {
		fmt.Fprintf(page, "Arguments and flags are forwarded to the native CLI. Run `%s --help` for native options.\n\n", command.CommandPath())
	}
	page.WriteString("```text\n")
	command.SetOut(page)
	if err := command.Help(); err != nil {
		return fmt.Errorf("render help for %q: %w", command.CommandPath(), err)
	}
	page.WriteString("```\n\n")
	for _, child := range command.Commands() {
		if err := renderCommand(page, child); err != nil {
			return fmt.Errorf("render subcommand of %q: %w", command.CommandPath(), err)
		}
	}
	return nil
}
