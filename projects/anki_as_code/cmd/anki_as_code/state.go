package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	collectionpb "git.alwaldend.com/alwaldend/src/projects/anki_as_code/api/collection"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/app"
	"git.alwaldend.com/alwaldend/src/projects/anki_as_code/internal/codec"
	"github.com/spf13/cobra"
)

// stateFlags holds command-owned flag values shared by the state subcommands.
type stateFlags struct {
	// collection selects the source database or collection package.
	collection string
	// key selects the state namespace; empty uses the store's default namespace.
	key string
}

// newStateCommand groups state operations using application-managed collection access.
func newStateCommand(application *app.App, flags *stateFlags) *cobra.Command {
	command := &cobra.Command{
		Use:   "state",
		Short: "Inspect and manage collection-backed path bindings",
		Args:  cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			return fmt.Errorf("state requires a subcommand")
		},
	}
	command.PersistentFlags().StringVar(&flags.collection, collectionFlag.name, "", collectionFlag.usage)
	if err := command.MarkPersistentFlagRequired(collectionFlag.name); err != nil {
		panic(fmt.Errorf("register required state input: %w", err))
	}
	command.AddCommand(newStateDumpCommand(application, flags), newStateImportCommand(application, flags),
		newStateWriteCommand(flags, "replace", "Replace the whole stored state object", application.ColStateReplace),
		newStateWriteCommand(flags, "update", "Merge an external state object into stored state", application.ColStateUpdate))
	return command
}

// newStateDumpCommand prints stored bindings through application-managed access.
func newStateDumpCommand(application *app.App, flags *stateFlags) *cobra.Command {
	var format string
	command := &cobra.Command{
		Use: "dump", Short: "Dump the stored state object", Args: cobra.NoArgs,
		RunE: func(command *cobra.Command, _ []string) error {
			value, err := application.ColStateDump(command.Context(), flags.collection, flags.key)
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
func newStateImportCommand(application *app.App, flags *stateFlags) *cobra.Command {
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
			if err := application.ColStateImport(command.Context(), flags.collection, text, flags.key, bindings); err != nil {
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
func newStateWriteCommand(flags *stateFlags, name, summary string, write func(context.Context, string, string, *collectionpb.State) error) *cobra.Command {
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
			value := &collectionpb.State{}
			if err := codec.Decode(data, codec.Format(format), value); err != nil {
				return fmt.Errorf("decode state for %s: %w", name, err)
			}
			if err := write(command.Context(), flags.collection, flags.key, value); err != nil {
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
