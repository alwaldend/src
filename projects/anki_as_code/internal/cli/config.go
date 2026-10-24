package cli

import (
	"fmt"
	"io"
	"os"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"
	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"github.com/spf13/cobra"
)

// newConfigCommand groups operations on text configuration independently of Anki collections.
func newConfigCommand(application ankiascodeapi.App) *cobra.Command {
	command := &cobra.Command{
		Use:     "config",
		Aliases: []string{"conf"},
		Short:   "Inspect and manipulate collection configuration files",
		Args:    cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return fmt.Errorf("config requires a subcommand")
		},
	}
	command.AddCommand(newConfigDumpCommand(application), newConfigExportCommand(application))
	return command
}

// newConfigExportCommand writes a serialized configuration set as a TOML directory tree.
func newConfigExportCommand(application ankiascodeapi.App) *cobra.Command {
	var output, format string
	command := &cobra.Command{
		Use:   "export [FILE]",
		Short: "Write a configuration set as TOML files",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, arguments []string) error {
			var data []byte
			var err error
			if len(arguments) == 0 || arguments[0] == "-" {
				data, err = io.ReadAll(command.InOrStdin())
			} else {
				data, err = os.ReadFile(arguments[0])
			}
			if err != nil {
				return fmt.Errorf("read configuration set: %w", err)
			}
			configs := &ankiascodepb.CollectionConfigSet{}
			if err := codec.Decode(data, format, configs); err != nil {
				return fmt.Errorf("decode configuration set: %w", err)
			}
			if err := application.ConfigExport(command.Context(), output, configs); err != nil {
				return fmt.Errorf("export configuration set: %w", err)
			}
			if _, err := fmt.Fprintln(command.ErrOrStderr(), "Exported configuration files to", output); err != nil {
				return fmt.Errorf("write configuration export status: %w", err)
			}
			return nil
		},
	}
	outputFlag.bindRequired(command, &output)
	command.Flags().StringVar(&format, "format", codec.JSON, "Input format: json or toml")
	return command
}

// newConfigDumpCommand defines a read-only configuration dump with selectable output format.
func newConfigDumpCommand(application ankiascodeapi.App) *cobra.Command {
	var text, format string
	command := &cobra.Command{
		Use:   "dump",
		Short: "Print parsed collection configurations",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			data, err := application.ConfigDump(command.Context(), text, format)
			if err != nil {
				return fmt.Errorf("dump configuration: %w", err)
			}
			if _, err := command.OutOrStdout().Write(data); err != nil {
				return fmt.Errorf("write configuration dump: %w", err)
			}
			return nil
		},
	}
	textFlag.bindRequired(command, &text)
	command.Flags().StringVar(&format, "format", "json", "Output format: json or toml")
	return command
}
