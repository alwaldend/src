package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// execute dispatches README commands with caller-owned streams and cancellation.
func execute(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	var changeDir string
	root := &cobra.Command{
		Use:           "readme-tree",
		Short:         "Generate a tree from README metadata",
		SilenceErrors: true,
	}
	root.PersistentFlags().StringVarP(&changeDir, "change-dir", "C", "", "Resolve paths from this directory")
	root.AddCommand(newParseCommand(&changeDir))
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs(args)
	if err := root.ExecuteContext(ctx); err != nil {
		return fmt.Errorf("run readme-tree: %w", err)
	}
	return nil
}

// newParseCommand binds traversal and rendering flags for a single invocation.
func newParseCommand(changeDir *string) *cobra.Command {
	config := parseConfig{format: markdownFormat}
	cmd := &cobra.Command{
		Use:   "parse [path...]",
		Short: "Parse README files",
		RunE: func(cmd *cobra.Command, args []string) error {
			base, err := filepath.Abs(*changeDir)
			if err != nil {
				return fmt.Errorf("resolve working directory %q: %w", *changeDir, err)
			}
			if *changeDir != "" {
				info, err := os.Stat(base)
				if err != nil {
					return fmt.Errorf("inspect working directory %q: %w", base, err)
				}
				if !info.IsDir() {
					return fmt.Errorf("working directory %q is not a directory", base)
				}
			}
			config.root = resolvePath(base, config.root)
			for i, path := range config.exclude {
				config.exclude[i] = resolvePath(base, path)
			}
			result := parseResult{Dir: config.root}
			for _, path := range args {
				if err := result.collect(cmd.Context(), resolvePath(base, path), config); err != nil {
					return fmt.Errorf("collect README files under %q: %w", path, err)
				}
			}
			if err := render(cmd.OutOrStdout(), result, config.format); err != nil {
				return fmt.Errorf("render README tree: %w", err)
			}
			return nil
		},
	}
	flags := cmd.Flags()
	flags.VarP(&config.format, "output-type", "o", "Output format")
	flags.StringVarP(&config.name, "readme-name", "n", "README.md", "Filename to process")
	flags.StringVarP(&config.root, "root-path", "r", "", "Base for relative paths (default: working directory)")
	flags.StringSliceVarP(&config.exclude, "exclude", "e", nil, "Directories to exclude")
	flags.BoolVarP(&config.useGit, "git", "g", false, "Use Git to list tracked files")
	return cmd
}

// resolvePath anchors a relative path in the invocation's working directory.
func resolvePath(base, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(base, path)
}
