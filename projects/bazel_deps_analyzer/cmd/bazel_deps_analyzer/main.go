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

// newCommand exposes analysis modes without a root-level conversion route.
func newCommand() *cobra.Command {
	command := &cobra.Command{
		Use:           "bazel_deps_analyzer",
		Short:         "Analyze captured Bazel target and package dependencies",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          showHelp,
	}
	analyze := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze individual targets or aggregate packages",
		Args:  cobra.NoArgs,
		RunE:  showHelp,
	}
	analyze.AddCommand(newTargetCommand(), newPackageCommand())
	command.AddCommand(analyze)
	return command
}

// showHelp returns successful usage output without opening or decoding input.
func showHelp(command *cobra.Command, _ []string) error {
	if err := command.Help(); err != nil {
		return fmt.Errorf("show %s help: %w", command.CommandPath(), err)
	}
	return nil
}

// newTargetCommand preserves target-selection flags under analyze target.
func newTargetCommand() *cobra.Command {
	options := graph.Options{Depth: -1, Direction: "LR"}
	var inputPath, outputPath string
	command := &cobra.Command{
		Use:   "target",
		Short: "Convert a target dependency closure to Mermaid",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			convert := func(input io.Reader) ([]byte, error) {
				diagram, err := graph.Convert(input, options)
				if err != nil {
					return nil, fmt.Errorf("analyze target: %w", err)
				}
				return diagram, nil
			}
			if err := convertFiles(command.InOrStdin(), command.OutOrStdout(), inputPath, outputPath, convert); err != nil {
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
	command.Flags().BoolVar(&options.BuildDepsOnly, "build-deps-only", false, "Keep explicit deps/srcs label values intersected with captured rule inputs; requires complete attribute metadata")
	command.MarkFlagRequired("root")
	return command
}

// newPackageCommand owns aggregation flags independently of target traversal.
func newPackageCommand() *cobra.Command {
	var options graph.PackageOptions
	var inputPath, outputPath string
	command := &cobra.Command{
		Use:   "package",
		Short: "Aggregate every captured target into package nodes",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			convert := func(input io.Reader) ([]byte, error) {
				output, err := graph.ConvertPackages(input, options)
				if err != nil {
					return nil, fmt.Errorf("analyze packages: %w", err)
				}
				return output, nil
			}
			if err := convertFiles(command.InOrStdin(), command.OutOrStdout(), inputPath, outputPath, convert); err != nil {
				return fmt.Errorf("convert package graph: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&inputPath, "input", "-", "Input NDJSON capture, or - for stdin")
	command.Flags().StringVar(&outputPath, "output", "-", "Output graph file, or - for stdout")
	command.Flags().IntVar(&options.PackageDepth, "package-depth", 1, "Positive number of local package path segments to retain")
	command.Flags().StringVar(&options.Format, "format", "mermaid", "Output format: mermaid or json")
	command.Flags().StringVar(&options.Direction, "direction", "LR", "Mermaid flow direction: LR, RL, TB, BT")
	command.Flags().BoolVar(&options.IncludeFiles, "include-files", true, "Include source and generated file targets")
	command.Flags().BoolVar(&options.IncludeExternal, "include-external", true, "Include external repository targets in one aggregate")
	command.Flags().BoolVar(&options.MergeReciprocal, "merge-reciprocal", true, "Use double-headed Mermaid arrows for reciprocal relationships")
	command.Flags().BoolVar(&options.BuildDepsOnly, "build-deps-only", false, "Keep explicit deps/srcs label values intersected with captured rule inputs; requires complete attribute metadata")
	return command
}

// convertFiles closes caller-selected files and writes only after full conversion.
func convertFiles(input io.Reader, output io.Writer, inputPath, outputPath string, convert func(io.Reader) ([]byte, error)) (resultErr error) {
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
	diagram, err := convert(input)
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
