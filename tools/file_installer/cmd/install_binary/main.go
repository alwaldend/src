package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// main adapts installer arguments and reports operation failures.
func main() {
	var source, name, directory string
	command := &cobra.Command{
		Use:          "install_binary",
		Short:        "Install a standalone executable atomically",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(command *cobra.Command, arguments []string) error {
			if err := install(source, name, directory); err != nil {
				return fmt.Errorf("install executable: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&source, "source", "", "Standalone source executable")
	command.Flags().StringVar(&name, "name", "", "Installed command name")
	command.Flags().StringVar(&directory, "directory", "", "Destination directory (default: ~/.local/bin)")
	for _, flag := range []string{"source", "name"} {
		if err := command.MarkFlagRequired(flag); err != nil {
			fmt.Fprintf(os.Stderr, "configure installer: %s\n", err)
			os.Exit(1)
		}
	}
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

// install copies a regular executable before atomically replacing its destination.
func install(source, name, directory string) (result error) {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("invalid command name %q", name)
	}
	if directory == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("resolve home directory: %w", err)
		}
		directory = filepath.Join(home, ".local", "bin")
	}
	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open source %q: %w", source, err)
	}
	defer func() {
		if err := input.Close(); err != nil {
			result = errors.Join(result, fmt.Errorf("close source %q: %w", source, err))
		}
	}()
	info, err := input.Stat()
	if err != nil {
		return fmt.Errorf("inspect source %q: %w", source, err)
	}
	if !info.Mode().IsRegular() || info.Mode().Perm()&0o111 == 0 {
		return fmt.Errorf("source %q is not a regular executable", source)
	}
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return fmt.Errorf("create destination directory %q: %w", directory, err)
	}
	output, err := os.CreateTemp(directory, "."+name+"-*")
	if err != nil {
		return fmt.Errorf("create temporary executable: %w", err)
	}
	temporary := output.Name()
	defer func() {
		if err := os.Remove(temporary); err != nil && !errors.Is(err, os.ErrNotExist) {
			result = errors.Join(result, fmt.Errorf("remove temporary executable: %w", err))
		}
	}()
	_, copyErr := io.Copy(output, input)
	modeErr := output.Chmod(0o755)
	closeErr := output.Close()
	if err := errors.Join(copyErr, modeErr, closeErr); err != nil {
		return fmt.Errorf("write temporary executable: %w", err)
	}
	destination := filepath.Join(directory, name)
	if err := os.Rename(temporary, destination); err != nil {
		return fmt.Errorf("replace executable %q: %w", destination, err)
	}
	return nil
}
