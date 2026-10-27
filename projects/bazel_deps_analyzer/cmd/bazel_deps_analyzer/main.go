package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"git.alwaldend.com/alwaldend/src/projects/bazel_deps_analyzer/internal/graph"
	"github.com/spf13/cobra"
)

// main reports command failures once and returns a failing process status.
func main() {
	if err := newCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// newCommand constructs the CLI with stdin/stdout defaults and required root.
func newCommand() *cobra.Command {
	options := graph.Options{Depth: -1, Direction: "LR"}
	var inputPath, outputPath string
	command := &cobra.Command{
		Use:           "bazel_deps_analyzer",
		Short:         "Convert Bazel streamed_jsonproto target graphs to Mermaid",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(command *cobra.Command, _ []string) error {
			if err := convertFiles(command.InOrStdin(), command.OutOrStdout(), inputPath, outputPath, options); err != nil {
				return fmt.Errorf("convert dependency diagram: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&inputPath, "input", "-", "Input NDJSON file, or - for stdin")
	command.Flags().StringVar(&outputPath, "output", "-", "Output Mermaid file, or - for stdout")
	command.Flags().StringVar(&options.Root, "root", "", "Complete root target label")
	command.Flags().IntVar(&options.Depth, "depth", -1, "Maximum dependency distance; -1 includes the complete closure")
	command.Flags().StringVar(&options.DiagramType, "diagram-type", "flowchart", "Mermaid view: flowchart dependencies or mindmap ownership inventory")
	command.Flags().StringVar(&options.Direction, "direction", "LR", "Flow direction: LR, RL, TB, BT")
	command.Flags().BoolVar(&options.IncludeFiles, "include-files", false, "Show source and generated file targets")
	command.Flags().BoolVar(&options.IncludeExternal, "include-external", false, "Show external repository targets")
	command.Flags().BoolVar(&options.ExternalLeaves, "external-leaves", false, "Show the first external targets as leaves without their dependencies")
	command.Flags().BoolVar(&options.GroupByPackage, "group-by-package", false, "Group short target names by full package, nesting external packages by repository")
	command.MarkFlagRequired("root")
	return command
}

// convertFiles closes caller-selected files and writes only after full conversion.
func convertFiles(input io.Reader, output io.Writer, inputPath, outputPath string, options graph.Options) (resultErr error) {
	if inputPath != "-" {
		file, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("open graph %q: %w", inputPath, err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				resultErr = errors.Join(resultErr, fmt.Errorf("close graph %q: %w", inputPath, err))
			}
		}()
		input = file
	}
	diagram, err := graph.Convert(input, options)
	if err != nil {
		return fmt.Errorf("convert graph %q: %w", inputPath, err)
	}
	if outputPath != "-" {
		file, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("create diagram %q: %w", outputPath, err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				resultErr = errors.Join(resultErr, fmt.Errorf("close diagram %q: %w", outputPath, err))
			}
		}()
		output = file
	}
	if _, err := output.Write(diagram); err != nil {
		return fmt.Errorf("write diagram %q: %w", outputPath, err)
	}
	return nil
}
