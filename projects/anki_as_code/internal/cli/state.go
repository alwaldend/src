package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	ankiascodeapi "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodeapi"

	ankiascodepb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/ankiascodepb"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"github.com/spf13/cobra"
)

// newStateCommand groups state operations using application-managed collection access.
func newStateCommand(application ankiascodeapi.App) *cobra.Command {
	var input string
	command := &cobra.Command{
		Use:   "state",
		Short: "Inspect and manage collection-backed path bindings",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return fmt.Errorf("state requires a subcommand")
		},
	}
	command.PersistentFlags().StringVar(&input, collectionFlag.name, "", collectionFlag.usage)
	if err := command.MarkPersistentFlagRequired(collectionFlag.name); err != nil {
		panic(fmt.Errorf("register required state input: %w", err))
	}
	command.AddCommand(newStateDumpCommand(application, &input), newStateImportCommand(application, &input),
		newStateWriteCommand(&input, "replace", "Replace the whole stored state object", application.ColStateReplace),
		newStateWriteCommand(&input, "update", "Merge an external state object into stored state", application.ColStateUpdate))
	return command
}

// newStateDumpCommand prints stored bindings through application-managed access.
func newStateDumpCommand(application ankiascodeapi.App, input *string) *cobra.Command {
	var format string
	command := &cobra.Command{
		Use: "dump", Short: "Dump the stored state object", Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			value, err := application.ColStateDump(command.Context(), ankiascodeapi.CollectionPath(*input), commandStateKey(command))
			if err != nil {
				return fmt.Errorf("state dump: %w", err)
			}
			data, err := codec.Encode(value, codec.Format(format))
			if err != nil {
				return fmt.Errorf("encode dumped state: %w", err)
			}
			if _, err := command.OutOrStdout().Write(data); err != nil {
				return fmt.Errorf("write dumped state: %w", err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&format, "format", string(codec.TOML), "State format: toml or json")
	return command
}

// newStateImportCommand passes a complete PATH=ID batch to one managed state operation.
func newStateImportCommand(application ankiascodeapi.App, input *string) *cobra.Command {
	var text string
	command := &cobra.Command{
		Use: "import PATH=ID [PATH=ID...]", Short: "Bind external objects to configuration paths", Args: cobra.MinimumNArgs(1),
		RunE: func(command *cobra.Command, arguments []string) error {
			bindings := map[string]string{}
			for _, argument := range arguments {
				separator := strings.LastIndexByte(argument, '=')
				if separator <= 0 {
					return fmt.Errorf("binding %q must be PATH=ID", argument)
				}
				address := argument[:separator]
				if _, exists := bindings[address]; exists {
					return fmt.Errorf("binding path %q was supplied more than once", address)
				}
				bindings[address] = argument[separator+1:]
			}
			if err := application.ColStateImport(command.Context(), ankiascodeapi.CollectionPath(*input), text, commandStateKey(command), bindings); err != nil {
				return fmt.Errorf("state import: %w", err)
			}
			if _, err := fmt.Fprintln(command.ErrOrStderr(), "Imported collection bindings:", len(bindings)); err != nil {
				return fmt.Errorf("write import status: %w", err)
			}
			return nil
		},
	}
	textFlag.bindRequired(command, &text)
	return command
}

// newStateWriteCommand decodes an external object and invokes its selected state operation.
func newStateWriteCommand(input *string, name, summary string, write func(context.Context, ankiascodeapi.CollectionSource, string, *ankiascodepb.State) error) *cobra.Command {
	var format string
	command := &cobra.Command{
		Use: name + " [FILE]", Short: summary, Args: cobra.MaximumNArgs(1),
		RunE: func(command *cobra.Command, arguments []string) error {
			var data []byte
			var err error
			if len(arguments) == 0 || arguments[0] == "-" {
				data, err = io.ReadAll(command.InOrStdin())
			} else {
				data, err = os.ReadFile(arguments[0])
			}
			if err != nil {
				return fmt.Errorf("read state for %s: %w", name, err)
			}
			value := &ankiascodepb.State{}
			if err := codec.Decode(data, codec.Format(format), value); err != nil {
				return fmt.Errorf("decode state for %s: %w", name, err)
			}
			if err := write(command.Context(), ankiascodeapi.CollectionPath(*input), commandStateKey(command), value); err != nil {
				return fmt.Errorf("state %s: %w", name, err)
			}
			if _, err := fmt.Fprintf(command.ErrOrStderr(), "Collection state %s completed\n", name); err != nil {
				return fmt.Errorf("write %s status: %w", name, err)
			}
			return nil
		},
	}
	command.Flags().StringVar(&format, "format", string(codec.TOML), "State format: toml or json")
	return command
}
