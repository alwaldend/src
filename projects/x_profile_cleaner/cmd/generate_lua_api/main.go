package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"git.alwaldend.com/alwaldend/src/projects/x_profile_cleaner/internal/luaapi"
	"github.com/spf13/cobra"
)

// main reports generator failures without printing source data or credentials.
func main() {
	if err := command().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// command exposes the repository-owned generation workflow with optional file output.
func command() *cobra.Command {
	var output string
	command := &cobra.Command{Use: "generate_lua_api", Short: "Generate Lua editor types from the worker plan schema", Args: cobra.NoArgs, SilenceUsage: true, SilenceErrors: true}
	command.RunE = func(cmd *cobra.Command, _ []string) error {
		data, err := luaapi.Generate()
		if err != nil {
			return fmt.Errorf("generate Lua plan types: %w", err)
		}
		if output == "" {
			if _, err := cmd.OutOrStdout().Write(data); err != nil {
				return fmt.Errorf("write Lua declarations: %w", err)
			}
			return nil
		}
		if !filepath.IsAbs(output) {
			return fmt.Errorf("output path must be absolute")
		}
		if err := write(output, data); err != nil {
			return fmt.Errorf("update Lua API file: %w", err)
		}
		return nil
	}
	command.Flags().StringVar(&output, "out", "", "absolute output path; omitted writes stdout")
	return command
}

// write replaces the public projection only after a complete successful write.
func write(path string, data []byte) (resultErr error) {
	file, err := os.CreateTemp(filepath.Dir(path), ".lua-api-*")
	if err != nil {
		return fmt.Errorf("create Lua declaration file: %w", err)
	}
	defer func() {
		if err := os.Remove(file.Name()); err != nil && !errors.Is(err, os.ErrNotExist) {
			resultErr = errors.Join(resultErr, fmt.Errorf("remove temporary declarations: %w", err))
		}
	}()
	if err := file.Chmod(0o644); err != nil {
		return errors.Join(fmt.Errorf("set Lua declaration mode: %w", err), file.Close())
	}
	if _, err := file.Write(data); err != nil {
		return errors.Join(fmt.Errorf("write Lua declaration file: %w", err), file.Close())
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close Lua declaration file: %w", err)
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("replace Lua declaration file: %w", err)
	}
	return nil
}
